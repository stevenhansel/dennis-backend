package votes

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/stevenhansel/csm-ending-prediction-be/internal/server/responseutil"
)

type VoteHttpController struct {
	responseutil *responseutil.Responseutil
	service      *VoteService
}

func NewVoteHttpController(responseutil *responseutil.Responseutil, service *VoteService) *VoteHttpController {
	return &VoteHttpController{
		responseutil: responseutil,
		service:      service,
	}
}

func (c *VoteHttpController) InsertVote(w http.ResponseWriter, r *http.Request) {
	res := c.responseutil.CreateResponse(w)

	decoder := json.NewDecoder(r.Body)
	var body InsertVoteParams
	err := decoder.Decode(&body)
	if err != nil {
		res.Error4xx(http.StatusBadRequest, "Request body is invalid")
		return
	}

	forwardedIP := r.Header.Get("X-Forwarded-For")
	splittedIPs := strings.Split(forwardedIP, ",")
	if len(splittedIPs) > 0 && splittedIPs[0] != "" {
		body.IPAddress = splittedIPs[0]
	} else {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			res.Error4xx(http.StatusBadRequest, "IP Address is invalid")
			return
		}

		clientIPAddress := net.ParseIP(ip)
		body.IPAddress = clientIPAddress.String()
	}

	if err := c.service.InsertVote(r.Context(), &body); err != nil {
		if realError := (ErrEpisodeNotFound{}); errors.As(err, &realError) {
			res.Error4xx(
				http.StatusNotFound,
				err.Error(),
			)
			return
		}

		res.Error5xx(err)
		return
	}

	res.JSON(http.StatusNoContent, nil)
}

func (c *VoteHttpController) GetVotesByEpisodeID(w http.ResponseWriter, r *http.Request) {
	res := c.responseutil.CreateResponse(w)

	strEpisodeID := chi.URLParam(r, "episodeId")
	episodeID, err := strconv.Atoi(strEpisodeID)
	if err != nil {
		res.Error4xx(http.StatusBadRequest, "Episode ID must be a valid integer")
		return
	}

	votes, err := c.service.GetVotesByEpisodeID(r.Context(), episodeID)
	if err != nil {
		res.Error5xx(err)
		return
	}

	res.JSON(http.StatusOK, votes)
}
