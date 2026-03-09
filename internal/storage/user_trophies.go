package storage

import (
	"context"
	"database/sql"
)

// user_id bigint PRIMARY KEY REFERENCES users(id),
//  trophies int NOT NULL CHECK (trophies >= 0),
//  updated_at timestamptz NOT NULL DEFAULT now()

type UserTrophies struct {
	UserID    int64  `json:"user_id"`
	Trophies  int64  `json:"trophies"`
	UpdatedAt string `json:"updated_at"`
}

type UserTrophyStore struct {
	db *sql.DB
}

func (ut *UserTrophyStore) Upsert(ctx context.Context, userID int64, delta int64) error {

	return withTx(ut.db, ctx, func(tx *sql.Tx) error {
		query := `
		INSERT INTO user_trophies (user_id, trophies) VALUES ($1, $2)
		ON CONFLICT (user_id)
		DO UPDATE SET
			trophies = user_trophies.trophies + EXCLUDE.trophies
			updated_at = now()
		`

		ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
		defer cancel()

		err := tx.QueryRowContext(
			ctx,
			query,
			userID,
			delta,
		).Err()

		if err != nil {
			switch {
			case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
				return ErrDuplicateEmail
			case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
				return ErrDuplicateUsername
			default:
				return err
			}
		}

		return nil
	})

}

func (ut *UserTrophyStore) GetTrophyCountByID(ctx context.Context, userID int64) (int64, error) {
}

func (ut *UserTrophyStore) GetTrophies(ctx context.Context) ([]UserTrophies, error) {
}
