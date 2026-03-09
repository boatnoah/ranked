CREATE TABLE IF NOT EXISTS matches (
  id bigserial PRIMARY KEY,
  user_id bigint NOT NULL REFERENCES users(id),
  result varchar(10) NOT NULL CHECK (result IN ('win', 'loss')),
  crowns int NOT NULL CHECK (crowns BETWEEN 1 AND 3),
  trophies_changed int NOT NULL,
  submitted_at timestamptz NOT NULL DEFAULT now()
);
