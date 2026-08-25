-- +goose Up


CREATE TABLE IF NOT EXISTS profiles (
    user_id UUID PRIMARY KEY NOT NULL,
    display_name TEXT NOT NULL ,
    bio TEXT ,
    avatar_key TEXT ,
    is_private BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);


-- +goose Down
DROP TABLE IF EXISTS profiles;