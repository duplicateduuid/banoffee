-- Add migration script here
CREATE TABLE "user" (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    username VARCHAR(255),
    avatar_url VARCHAR(255),
    header_url VARCHAR(255),
    bio VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ,
    banned_at TIMESTAMPTZ
);
