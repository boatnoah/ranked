package storage

import (
	"context"
	"database/sql"
)

type UserTrophyTxStore struct {
	tx *sql.Tx
}

func (ut *UserTrophyTxStore) Upsert(ctx context.Context, userID int64, delta int64) error {

	query := `
		INSERT INTO user_trophies (user_id, trophies) VALUES ($1, $2)
		ON CONFLICT (user_id)
		DO UPDATE SET
			trophies = user_trophies.trophies + EXCLUDED.trophies,
			updated_at = now()
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := ut.tx.QueryRowContext(
		ctx,
		query,
		userID,
		delta,
	).Err()

	if err != nil {
		return err
	}

	return nil

}

func (ut *UserTrophyTxStore) GetTrophyCountByID(ctx context.Context, userID int64) (*UserTrophies, error) {
	query := `
		SELECT trophies FROM user_trophies
		WHERE user_id = $1
		`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var userTrophies UserTrophies
	err := ut.tx.QueryRowContext(
		ctx,
		query,
		userID,
	).Scan(
		&userTrophies.UserID,
		&userTrophies.Trophies,
		&userTrophies.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &userTrophies, err

}

func (ut *UserTrophyTxStore) GetAllTrophies(ctx context.Context) ([]UserTrophies, error) {

	query := `
		SELECT user_id, trophies FROM user_trophies
	`

	rows, err := ut.tx.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	var userTrophies []UserTrophies

	for rows.Next() {
		var user UserTrophies
		err := rows.Scan(&user.UserID, &user.Trophies, &user.UpdatedAt)

		if err != nil {
			return nil, err
		}

		userTrophies = append(userTrophies, user)
	}

	return userTrophies, nil

}
