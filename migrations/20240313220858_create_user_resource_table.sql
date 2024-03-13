-- Add migration script here
CREATE TABLE "user_resource" (
    user_id UUID,
    resource_id UUID,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ,
    PRIMARY KEY (user_d, resource_id),

    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES "user"(id),
    CONSTRAINT fk_resource FOREIGN KEY(resource_id) REFERENCES "resource"(id)
);