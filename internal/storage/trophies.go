package storage

// user_id bigint PRIMARY KEY REFERENCES users(id),
//  trophies int NOT NULL CHECK (trophies >= 0),
//  updated_at timestamptz NOT NULL DEFAULT now()

type Trophy struct {
	UserID int64 `json:"user_id"`
	Trophies
}
