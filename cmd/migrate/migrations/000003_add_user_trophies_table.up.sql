CREATE TABLE IF NOT EXISTS user_trophies (
  user_id bigint PRIMARY KEY REFERENCES users(id),
  trophies int NOT NULL CHECK (trophies >= 0),
  updated_at timestamptz NOT NULL DEFAULT now()
);
