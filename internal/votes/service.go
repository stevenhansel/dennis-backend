package votes

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/stevenhansel/csm-ending-prediction-be/internal/errtrace"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/querier"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/querier/database"
)

type Querier interface {
	InsertVote(ctx context.Context, params *database.InsertVoteParams) error
	FindVotes(ctx context.Context, params *database.FindVotesParams) ([]*querier.Vote, error)
	UpdateVoteEpisodeSongID(ctx context.Context, voteID int, episodeSongID int) error
	FindEpisodeDetailByEpisodeSongID(ctx context.Context, episodeSongID int) (*querier.EpisodeDetail, error)
	FindEpisodeVotes(ctx context.Context, episodeID int) ([]*querier.EpisodeVote, error)
}

type VoteService struct {
	querier Querier
}

func NewService(querier Querier) *VoteService {
	return &VoteService{
		querier: querier,
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
		return s.querier.UpdateVoteEpisodeSongID(ctx, existingVoteID, params.EpisodeSongID)
	} else {
		return s.querier.InsertVote(ctx, &database.InsertVoteParams{
			IPAddress:     params.IPAddress,
			EpisodeSongID: params.EpisodeSongID,
		})
	}
}

func (c *VoteService) GetVotesByEpisodeID(ctx context.Context, episodeID int) ([]*querier.EpisodeVote, error) {
	return c.querier.FindEpisodeVotes(ctx, episodeID)
}
