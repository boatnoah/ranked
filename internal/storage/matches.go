package storage

import (
	"context"
	"database/sql"
	"errors"
)

// id bigserial PRIMARY KEY,
// user_id bigint NOT NULL REFERENCES users(id),
// result varchar(10) NOT NULL CHECK (result IN ('win', 'loss')),
// crowns int NOT NULL CHECK (crowns BETWEEN 1 AND 3),
// trophies_changed int NOT NULL,
// submitted_at timestamptz NOT NULL DEFAULT now()
type Matches struct {
	ID              int64  `json:"id"`
	UserID          int64  `json:"user_id"`
	Result          string `json:"result"`
	Crowns          int    `json:"crowns"`
	TrophiesChanged int64  `json:"trophies_change"`
	SubmittedAt     string `json:"submitted_at"`
}

//
// 		Create(context.Context, int64, string, int, int64) error
// 		GetMatchesByID(context.Context, int64) ([]Matches, error)

type MatcheStore struct {
	db *sql.DB
}

func (ms *MatcheStore) Create(ctx context.Context, userID int64, result string, crowns int, delta int64) error {
	query := `
		INSERT INTO matches (user_id, result, crowns, trophies_changed) 
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := ms.db.ExecContext(ctx, query, userID, result, crowns, delta)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("Something went wrong unable to log match")
	}

	return nil

}

func (ms *MatcheStore) GetMatchByUserID(ctx context.Context, userID int64) ([]Matches, error) {
	// todo
}
