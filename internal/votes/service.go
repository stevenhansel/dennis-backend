package votes

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/stevenhansel/csm-ending-prediction-be/internal/errtrace"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/querier"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/querier/database"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/socket"
)

type Querier interface {
	InsertVote(ctx context.Context, params *database.InsertVoteParams) error
	FindVotes(ctx context.Context, params *database.FindVotesParams) ([]*querier.Vote, error)
	UpdateVoteEpisodeSongID(ctx context.Context, voteID int, episodeSongID int) error
	FindEpisodeDetailByID(ctx context.Context, episodeID int) (*querier.EpisodeDetail, error)
	FindEpisodeDetailByEpisodeSongID(ctx context.Context, episodeSongID int) (*querier.EpisodeDetail, error)
	FindEpisodeVotes(ctx context.Context, episodeID int) ([]*querier.EpisodeVote, error)
}

type SocketStateQuerier interface {
	Publish(payload socket.Payload)
}

type VoteService struct {
	querier Querier
	socket  SocketStateQuerier
}

func NewService(querier Querier, socket SocketStateQuerier) *VoteService {
	return &VoteService{
		querier: querier,
		socket:  socket,
	}
}

type InsertVoteParams struct {
	IPAddress     string
	EpisodeSongID int `json:"episodeSongId"`
}

func (s *VoteService) InsertVote(ctx context.Context, params *InsertVoteParams) error {
	episode, err := s.querier.FindEpisodeDetailByEpisodeSongID(ctx, params.EpisodeSongID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errtrace.Wrap(NewErrEpisodeNotFound(fmt.Errorf("Episode not found")))
		}

		return errtrace.Wrap(err)
	}

	episodeSongIDs := make([]int, len(episode.Songs))
	for i, s := range episode.Songs {
		episodeSongIDs[i] = s.EpisodeSongID
	}

	votes, err := s.querier.FindVotes(ctx, &database.FindVotesParams{
		EpisodeSongIDs: episodeSongIDs,
	})
	if err != nil {
		return errtrace.Wrap(err)
	}

	var existingVoteID int
	for _, v := range votes {
		if v.IPAddress == params.IPAddress {
			existingVoteID = v.ID
			break
		}
	}

	if existingVoteID != 0 {
		err = s.querier.UpdateVoteEpisodeSongID(ctx, existingVoteID, params.EpisodeSongID)
	} else {
		err = s.querier.InsertVote(ctx, &database.InsertVoteParams{
			IPAddress:     params.IPAddress,
			EpisodeSongID: params.EpisodeSongID,
		})
	}

	if err != nil {
		return errtrace.Wrap(err)
	}

	return s.PublishVoteUpdate(ctx, episode.ID)
}

func (s *VoteService) HasVoted(ctx context.Context, episodeID int, ipAddress string) (bool, *int, error) {
	episode, err := s.querier.FindEpisodeDetailByID(ctx, episodeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil, errtrace.Wrap(NewErrEpisodeNotFound(fmt.Errorf("Episode not found")))
		}

		return false, nil, errtrace.Wrap(err)
	}

	episodeSongIDs := make([]int, len(episode.Songs))
	for i, s := range episode.Songs {
		episodeSongIDs[i] = s.EpisodeSongID
	}

	votes, err := s.querier.FindVotes(ctx, &database.FindVotesParams{
		EpisodeSongIDs: episodeSongIDs,
	})
	if err != nil {
		return false, nil, errtrace.Wrap(err)
	}

	for _, v := range votes {
		if v.IPAddress == ipAddress {
			return true, &v.EpisodeSongID, nil
		}
	}

	return false, nil, nil
}

func (c *VoteService) GetVotesByEpisodeID(ctx context.Context, episodeID int) ([]*querier.EpisodeVote, error) {
	return c.querier.FindEpisodeVotes(ctx, episodeID)
}

func (c *VoteService) PublishVoteUpdate(ctx context.Context, episodeID int) error {
	votes, err := c.GetVotesByEpisodeID(ctx, episodeID)
	if err != nil {
		return errtrace.Wrap(err)
	}

	c.socket.Publish(socket.Payload{
		Topic:     socket.NewVoteTopic,
		EpisodeID: episodeID,
		Message:   votes,
	})

	return nil
}
