package storage

import "database/sql"

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
	SubmittedAt     string `json:"created_at"`
}

type MatcheStore struct {
	db *sql.DB
}
