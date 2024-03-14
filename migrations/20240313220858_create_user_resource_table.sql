-- Add migration script here
CREATE TABLE "user_resource" (
    user_id UUID,
    resource_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, resource_id),

    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES "user"(id),
    CONSTRAINT fk_resource FOREIGN KEY(resource_id) REFERENCES "resource"(id)
);