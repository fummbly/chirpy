-- +goose Up
CREATE TABLE chirps(
  id UUID PRIMARY KEY,
  created_at TIME NOT NULL,
  updated_at TIME NOT NULL,
  body TEXT NOT NULL,
  user_id UUID NOT NULL,
  FOREIGN KEY(user_id) REFERENCES users(id)
  ON DELETE CASCADE
);

-- +goose Down
DROP TABLE chirps;