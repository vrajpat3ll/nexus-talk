-- Messaging service schema
-- Messages, message states, threads, attachments

CREATE TABLE threads (
    id UUID PRIMARY KEY,
    type VARCHAR(16) CHECK (type IN ('direct','group','channel')),
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE messages (
    id UUID PRIMARY KEY,
    sender_id UUID,
    thread_id UUID REFERENCES threads(id),
    content TEXT,
    sent_at TIMESTAMP DEFAULT now(),
    content_type VARCHAR(32),
    is_deleted BOOLEAN DEFAULT FALSE,
    reply_to UUID,
    metadata JSONB
);

CREATE TABLE message_states (
    message_id UUID REFERENCES messages(id),
    user_id UUID,
    state VARCHAR(16), -- sent, delivered, read, recalled
    timestamp TIMESTAMP DEFAULT now(),
    PRIMARY KEY (message_id, user_id)
);

CREATE TABLE attachments (
    id UUID PRIMARY KEY,
    message_id UUID REFERENCES messages(id),
    file_url VARCHAR(256),
    media_type VARCHAR(32),
    size_bytes INTEGER,
    hash VARCHAR(128)
);
