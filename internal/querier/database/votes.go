package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/errtrace"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/querier"
)

func toVotes(rows ...*VoteRow) []*querier.Vote {
	results := make([]*querier.Vote, len(rows))
	for i, r := range rows {
		results[i] = toVote(r)
	}

	return results
}

func toVote(row *VoteRow) *querier.Vote {
	return &querier.Vote{
		ID:            row.ID,
		IPAddress:     row.IPAddress,
		EpisodeSongID: row.EpisodeSongID,
		CreatedAt:     row.CreatedAt,
	}
}

type InsertVoteParams struct {
	IPAddress     string `db:"ip_address"`
	EpisodeSongID int    `db:"episode_song_id"`
}

func (d *DatabaseQuerier) InsertVote(ctx context.Context, params *InsertVoteParams) error {
	statement := `
  insert into "vote"
  ("ip_address", "episode_song_id")
  values (:ip_address, :episode_song_id)
  `
	if _, err := d.db.NamedExecContext(ctx, statement, params); err != nil {
		return errtrace.Wrap(err)
	}

	return nil
}

func (d *DatabaseQuerier) DeleteVoteByID(ctx context.Context, voteID int) error {
	statement := `
  delete from "vote"
  where "id" = $1
  `

	if _, err := d.db.ExecContext(ctx, statement, voteID); err != nil {
		return errtrace.Wrap(err)
	}

	return nil
}

func (d *DatabaseQuerier) UpdateVoteEpisodeSongID(ctx context.Context, voteID int, episodeSongID int) error {
	statement := `
    update "vote"
    set "episode_song_id" = $2
    where "id" = $1
  `

	if _, err := d.db.ExecContext(ctx, statement, voteID, episodeSongID); err != nil {
		return errtrace.Wrap(err)
	}

	return nil
}

type VoteRow struct {
	ID            int       `db:"vote_id"`
	IPAddress     string    `db:"vote_ip_address"`
	EpisodeSongID int       `db:"vote_episode_song_id"`
	CreatedAt     time.Time `db:"vote_created_at"`
}

type FindVotesParams struct {
	IPAddress      *string `db:"ip_address"`
	EpisodeSongID  *int    `db:"episode_song_id"`
	EpisodeSongIDs []int   `db:"episode_song_ids"`
}

func (d *DatabaseQuerier) FindVotes(ctx context.Context, params *FindVotesParams) ([]*querier.Vote, error) {
	var where []string
	if params.IPAddress != nil {
		where = append(where, `"v"."ip_address" = :ip_address`)
	}
	if params.EpisodeSongID != nil {
		where = append(where, `"v"."episode_song_id" = :episode_song_id`)
	}
	if len(params.EpisodeSongIDs) > 0 {
		where = append(where, `"v"."episode_song_id" in (:episode_song_ids)`)
	}

	var whereQuery string
	if len(where) > 0 {
		for i, w := range where {
			query := w
			if i != len(where)-1 {
				query += ` and `
			}

			whereQuery += query
		}
	}

	statement := fmt.Sprintf(`
    select
      "v"."id" as "vote_id",
      "v"."ip_address" as "vote_ip_address",
      "v"."episode_song_id" as "vote_episode_song_id",
      "v"."created_at" as "vote_created_at"
    from "vote" "v"
    where %s
    `,
		whereQuery,
	)

	query, args, err := sqlx.Named(statement, params)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	query = d.db.Rebind(query)

	rows, err := d.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	defer rows.Close()

	var results []*VoteRow
	for rows.Next() {
		var row VoteRow

		if err := rows.StructScan(&row); err != nil {
			return nil, errtrace.Wrap(err)
		}

		results = append(results, &row)
	}

	return toVotes(results...), nil
}
