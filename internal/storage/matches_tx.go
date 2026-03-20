package storage

import (
	"context"
	"database/sql"
	"errors"
)

type MatchTxStore struct {
	tx *sql.Tx
}

func (ms *MatchTxStore) Create(ctx context.Context, userID int64, result string, crowns int64, delta int64) error {

	query := `
		INSERT INTO matches (user_id, result, crowns, trophies_changed) VALUES ($1, $2, $3, $4)
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := ms.tx.ExecContext(ctx, query, userID, result, crowns, delta)

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

func (ms *MatchTxStore) GetMatchesByUserID(ctx context.Context, userID int64) ([]Matches, error) {

	query := `
		SELECT id, result, crowns, trophies_changed, submitted_at FROM matches
		WHERE = $1
	`

	rows, err := ms.tx.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	var matches []Matches

	for rows.Next() {
		var match Matches
		match.UserID = userID
		err := rows.Scan(&match.ID, &match.Result, &match.Crowns, &match.TrophiesChanged, &match.SubmittedAt)

		if err != nil {
			return nil, err
		}

		matches = append(matches, match)
	}
	return matches, nil
}
