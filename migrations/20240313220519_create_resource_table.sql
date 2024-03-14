-- Add migration script here
CREATE TABLE "resource" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    url VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    image_url VARCHAR(255),
    author VARCHAR(255),
    description VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
