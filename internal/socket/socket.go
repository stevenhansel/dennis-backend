package socket

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/server/responseutil"
	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type Topic string

const (
	NewVoteTopic  Topic = "new_vote"
	NewSubscriber Topic = "new_subscriber"
)

type Payload struct {
	Topic     Topic       `json:"topic"`
	EpisodeID int         `json:"episodeId"`
	Message   interface{} `json:"message"`
}

type SocketState struct {
	publishLimiter          *rate.Limiter
	subscriberMessageBuffer int
	responseutil            *responseutil.Responseutil
	corsOrigins             []string

	subscribers      map[int]map[*Subscriber]struct{}
	subscribersMutex sync.Mutex
}

func NewSocketState(
	publishLimiter *rate.Limiter,
	subscriberMessageBuffer int,
	responseutil *responseutil.Responseutil,
) *SocketState {
	return &SocketState{
		publishLimiter:          publishLimiter,
		subscriberMessageBuffer: subscriberMessageBuffer,
		responseutil:            responseutil,
		subscribers:             make(map[int]map[*Subscriber]struct{}),
	}
}

type Subscriber struct {
	messages  chan Payload
	closeSlow func()
}

type NumOfSubscribersResponseBody struct {
	NumOfSubscribers int `json:"numOfSubscribers"`
}

func (s *SocketState) GetNumOfSubscribers(w http.ResponseWriter, r *http.Request) {
	res := s.responseutil.CreateResponse(w)

	sEpisodeID := chi.URLParam(r, "episodeId")
	episodeID, err := strconv.Atoi(sEpisodeID)
	if err != nil {
		res.Error4xx(http.StatusBadRequest, "Failed to parse request")
		return
	}

	var numOfSubscribers int
	if s, ok := s.subscribers[episodeID]; ok {
		numOfSubscribers = len(s)
	}

	res.JSON(http.StatusOK, &NumOfSubscribersResponseBody{
		NumOfSubscribers: numOfSubscribers,
	})
}

func (s *SocketState) SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	sEpisodeID := chi.URLParam(r, "episodeId")
	episodeID, err := strconv.Atoi(sEpisodeID)
	if err != nil {
		return
	}

	origins := []string{"*"}
	if len(s.corsOrigins) > 0 {
		origins = s.corsOrigins
	}
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: origins,
	})
	if err != nil {
		return
	}
	defer c.Close(websocket.StatusInternalError, "Websocket connection is closed")

	err = s.Subscribe(r.Context(), c, episodeID)
	if errors.Is(err, context.Canceled) {
		return
	}
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
		websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return
	}
	if err != nil {
		return
	}
}

func (s *SocketState) Subscribe(ctx context.Context, c *websocket.Conn, episodeID int) error {
	ctx = c.CloseRead(ctx)

	subscriber := &Subscriber{
		messages: make(chan Payload),
		closeSlow: func() {
			c.Close(websocket.StatusPolicyViolation, "connection too slow to keep up with messages")
		},
	}

	s.AddSubscriber(episodeID, subscriber)
	defer s.DeleteSubscriber(episodeID, subscriber)

	for {
		select {
		case msg := <-subscriber.messages:
			if err := writeTimeout(ctx, time.Second*5, c, msg); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}

}

func (s *SocketState) Publish(payload Payload) {
	s.subscribersMutex.Lock()
	defer s.subscribersMutex.Unlock()

	s.publishLimiter.Wait(context.Background())

	if _, ok := s.subscribers[payload.EpisodeID]; ok {
		for subscriber := range s.subscribers[payload.EpisodeID] {
			select {
			case subscriber.messages <- payload:
			default:
				go subscriber.closeSlow()
			}
		}
	}
}

func (s *SocketState) PublishSubscribersCount(episodeID int, subscriber **Subscriber) {
	payload := Payload{
		Topic:     NewSubscriber,
		EpisodeID: episodeID,
		Message: NumOfSubscribersResponseBody{
			NumOfSubscribers: len(s.subscribers[episodeID]),
		},
	}

	for sub := range s.subscribers[episodeID] {
		if sub != *subscriber {
			select {
			case sub.messages <- payload:
			default:
				go sub.closeSlow()
			}
		}
	}
}

func (s *SocketState) AddSubscriber(episodeID int, subscriber *Subscriber) {
	s.subscribersMutex.Lock()
	defer s.subscribersMutex.Unlock()

	if _, ok := s.subscribers[episodeID]; !ok {
		s.subscribers[episodeID] = map[*Subscriber]struct{}{}
	}
	s.subscribers[episodeID][subscriber] = struct{}{}

	s.PublishSubscribersCount(episodeID, &subscriber)
}

func (s *SocketState) DeleteSubscriber(episodeID int, subscriber *Subscriber) {
	s.subscribersMutex.Lock()
	defer s.subscribersMutex.Unlock()

	if _, ok := s.subscribers[episodeID]; ok {
		delete(s.subscribers[episodeID], subscriber)

		if len(s.subscribers[episodeID]) == 0 {
			delete(s.subscribers, episodeID)
		}
	}

	s.PublishSubscribersCount(episodeID, &subscriber)
}

func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, message interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return wsjson.Write(ctx, c, message)
}
