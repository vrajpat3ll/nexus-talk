-- Identity service PostgreSQL schema
-- Users, Profiles, Devices, Sessions

-- User accounts and auth
CREATE TABLE users (
    id UUID PRIMARY KEY,
    username VARCHAR(32) UNIQUE NOT NULL,
    phone VARCHAR(16) UNIQUE,
    email VARCHAR(64),
    password_hash VARCHAR(128),
    created_at TIMESTAMP DEFAULT now(),
    is_active BOOLEAN DEFAULT TRUE
);

-- User profile data
CREATE TABLE user_profiles (
    user_id UUID PRIMARY KEY REFERENCES users(id),
    display_name VARCHAR(64),
    bio TEXT,
    avatar_url VARCHAR(256),
    privacy_settings JSONB
);

-- Devices linked to a user (multi-device support)
CREATE TABLE devices (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    device_info JSONB,
    push_token VARCHAR(256),
    created_at TIMESTAMP DEFAULT now()
);

-- Sessions for handling refresh tokens, multi-device login
CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    device_id UUID REFERENCES devices(id),
    refresh_token VARCHAR(256),
    issued_at TIMESTAMP DEFAULT now(),
    expires_at TIMESTAMP
);
