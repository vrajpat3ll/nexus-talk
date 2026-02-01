-- Notification service schema
-- Notifications, user device tokens

CREATE TABLE notifications (
    id UUID PRIMARY KEY,
    user_id UUID,
    message TEXT,
    sent_at TIMESTAMP DEFAULT now(),
    delivery_status VARCHAR(16)
);

CREATE TABLE user_device_tokens (
    user_id UUID,
    device_id UUID,
    push_token VARCHAR(256),
    last_registered TIMESTAMP DEFAULT now(),
    PRIMARY KEY(user_id, device_id)
);
