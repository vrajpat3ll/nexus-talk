-- Payments service schema
-- Transactions, subscriptions, invoices

CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    user_id UUID,
    type VARCHAR(32), -- deposit, withdrawal, payment
    amount DECIMAL(12,2),
    currency VARCHAR(8),
    status VARCHAR(16), -- pending, complete, failed
    reference VARCHAR(64),
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE subscriptions (
    id UUID PRIMARY KEY,
    user_id UUID,
    channel_id UUID,
    status VARCHAR(16), -- active, cancelled
    started_at TIMESTAMP DEFAULT now(),
    expired_at TIMESTAMP
);

CREATE TABLE invoices (
    id UUID PRIMARY KEY,
    subscription_id UUID REFERENCES subscriptions(id),
    amount DECIMAL(12,2),
    due_at TIMESTAMP,
    paid_at TIMESTAMP,
    status VARCHAR(16) -- paid, unpaid, overdue
);
