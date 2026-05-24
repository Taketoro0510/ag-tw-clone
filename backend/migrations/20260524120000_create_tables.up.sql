CREATE TABLE users (
    id UUID PRIMARY KEY,
    firebase_uid TEXT UNIQUE NOT NULL,
    email TEXT NOT NULL,
    display_name TEXT NOT NULL,
    avatar_url TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE posts (
    id UUID PRIMARY KEY,
    author_id UUID NOT NULL REFERENCES users(id),
    body TEXT NOT NULL,
    media_type TEXT CHECK (media_type IN ('image', 'video')),
    media_path TEXT,
    media_url TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_posts_timeline ON posts (created_at DESC, id DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_posts_author ON posts (author_id, created_at DESC) WHERE deleted_at IS NULL;

CREATE TABLE likes (
    post_id UUID NOT NULL REFERENCES posts(id),
    user_id UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (post_id, user_id)
);
CREATE INDEX idx_likes_user ON likes (user_id, created_at DESC);
