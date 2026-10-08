-- +goose Up
ALTER TABLE friend_requests DROP CONSTRAINT IF EXISTS friend_requests_from_user_id_to_user_id_key;

CREATE UNIQUE INDEX IF NOT EXISTS idx_friend_requests_pending_unique
ON friend_requests (from_user_id, to_user_id)
WHERE status = 'pending';

-- +goose Down
DROP INDEX IF EXISTS idx_friend_requests_pending_unique;

ALTER TABLE friend_requests
ADD CONSTRAINT friend_requests_from_user_id_to_user_id_key UNIQUE (from_user_id, to_user_id);