-- RTC service schema
-- Calls, call participants, call events (for group and direct calls)

CREATE TABLE calls (
    id UUID PRIMARY KEY,
    thread_id UUID,
    started_at TIMESTAMP DEFAULT now(),
    ended_at TIMESTAMP
);

CREATE TABLE call_participants (
    call_id UUID REFERENCES calls(id),
    user_id UUID,
    joined_at TIMESTAMP DEFAULT now(),
    left_at TIMESTAMP,
    is_recorder BOOLEAN DEFAULT FALSE,
    PRIMARY KEY(call_id, user_id)
);

CREATE TABLE call_events (
    id UUID PRIMARY KEY,
    call_id UUID REFERENCES calls(id),
    user_id UUID,
    event_type VARCHAR(32),
    event_time TIMESTAMP DEFAULT now(),
    metadata JSONB
);
