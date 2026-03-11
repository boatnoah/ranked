package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	UserStorage interface {
		Create(context.Context, *User) error
		Delete(context.Context, int64) error
		GetByID(context.Context, int64) (*User, error)
		GetByEmail(context.Context, string) (*User, error)
	}
	MatchStore interface {
		Create(context.Context, int64, string, int, int64) error
		GetMatchesByID(context.Context, int64) ([]Matches, error)
	}
	TrophyStore interface {
		Upsert(context.Context, int64, int64) error
		GetTrophyCountByID(context.Context, int64) (*UserTrophies, error)
		GetTrophies(context.Context) ([]UserTrophies, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		UserStorage: &UserStore{db},
		TrophyStore: &UserTrophyStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
