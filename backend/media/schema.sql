-- Media service schema
-- Files, thumbnails, media metadata, hashes

CREATE TABLE files (
    id UUID PRIMARY KEY,
    owner_id UUID,
    file_url VARCHAR(256),
    media_type VARCHAR(32),
    uploaded_at TIMESTAMP DEFAULT now(),
    size_bytes INTEGER,
    hash VARCHAR(128)
);

CREATE TABLE thumbnails (
    id UUID PRIMARY KEY,
    file_id UUID REFERENCES files(id),
    url VARCHAR(256),
    width INTEGER,
    height INTEGER
);

CREATE TABLE media_metadata (
    file_id UUID REFERENCES files(id),
    metadata JSONB,
    PRIMARY KEY(file_id)
);
