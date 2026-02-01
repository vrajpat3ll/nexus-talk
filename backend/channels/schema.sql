-- Channels & Groups service schema
-- Channels, group memberships, roles, scheduled posts

CREATE TABLE channels (
    id UUID PRIMARY KEY,
    owner_id UUID,
    name VARCHAR(64),
    description TEXT,
    is_public BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE memberships (
    channel_id UUID REFERENCES channels(id),
    user_id UUID,
    role VARCHAR(16) CHECK (role IN ('owner','admin','member')),
    joined_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY(channel_id, user_id)
);

CREATE TABLE scheduled_posts (
    id UUID PRIMARY KEY,
    channel_id UUID REFERENCES channels(id),
    author_id UUID,
    content TEXT,
    post_time TIMESTAMP
);
