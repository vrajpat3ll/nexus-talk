-- Analytics service schema
-- Events, metrics, counters

CREATE TABLE events (
    id UUID PRIMARY KEY,
    user_id UUID,
    event_type VARCHAR(64),
    event_time TIMESTAMP DEFAULT now(),
    metadata JSONB
);

CREATE TABLE metrics (
    id SERIAL PRIMARY KEY,
    name VARCHAR(64),
    value BIGINT,
    collected_at TIMESTAMP DEFAULT now()
);

CREATE TABLE counters (
    name VARCHAR(64) PRIMARY KEY,
    value BIGINT,
    updated_at TIMESTAMP DEFAULT now()
);
