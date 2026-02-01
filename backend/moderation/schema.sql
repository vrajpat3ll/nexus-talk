-- Moderation service schema
-- Reports, moderation actions, block list, automated sanctions

CREATE TABLE reports (
    id UUID PRIMARY KEY,
    reporter_id UUID,
    target_id UUID,
    reason_code VARCHAR(32),
    include_context BOOLEAN,
    created_at TIMESTAMP DEFAULT now(),
    status VARCHAR(16) -- open, closed, escalated
);

CREATE TABLE moderation_actions (
    id UUID PRIMARY KEY,
    report_id UUID REFERENCES reports(id),
    action VARCHAR(32), -- ban, warn, mute, etc
    moderator_id UUID,
    taken_at TIMESTAMP DEFAULT now(),
    comment TEXT
);

CREATE TABLE block_list (
    id UUID PRIMARY KEY,
    blocker_id UUID,
    blocked_id UUID,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE automod_strikes (
    id UUID PRIMARY KEY,
    user_id UUID,
    reason VARCHAR(128),
    issued_at TIMESTAMP DEFAULT now(),
    expiry TIMESTAMP,
    active BOOLEAN DEFAULT TRUE
);
