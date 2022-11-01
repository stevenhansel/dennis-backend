package episodes

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/server/responseutil"
)

type EpisodeHttpController struct {
	responseutil *responseutil.Responseutil
	service      *EpisodeService
}

func NewEpisodeHttpController(responseutil *responseutil.Responseutil, service *EpisodeService) *EpisodeHttpController {
	return &EpisodeHttpController{
		responseutil: responseutil,
		service:      service,
	}
}

func (c *EpisodeHttpController) GetCurrentEpisode(w http.ResponseWriter, r *http.Request) {
	res := c.responseutil.CreateResponse(w)

	episode, err := c.service.FindCurrentEpisode(r.Context())
	if err != nil {
		res.Error5xx(err)
		return
	}

	res.JSON(http.StatusOK, episode)
}

func (c *EpisodeHttpController) GetEpisodeByID(w http.ResponseWriter, r *http.Request) {
	res := c.responseutil.CreateResponse(w)

	strEpisodeID := chi.URLParam(r, "episodeId")
	episodeID, err := strconv.Atoi(strEpisodeID)
	if err != nil {
		res.Error4xx(http.StatusBadRequest, "Episode ID must be a valid integer")
		return
	}

	episode, err := c.service.FindEpisodeDetailByID(r.Context(), episodeID)
	if err != nil {
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

	res.JSON(http.StatusOK, episode)
}

func (c *EpisodeHttpController) GetAllEpisodes(w http.ResponseWriter, r *http.Request) {
	res := c.responseutil.CreateResponse(w)

	episodes, err := c.service.FindAllEpisodes(r.Context())
	if err != nil {
		res.Error5xx(err)
		return
	}

	res.JSON(http.StatusOK, episodes)
}
