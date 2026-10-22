-- +goose Up
CREATE INDEX IF NOT EXISTS idx_friendships_user_b_user_a ON friendships (user_id_b, user_id_a);

-- +goose Down
DROP INDEX IF EXISTS idx_friendships_user_b_user_a;