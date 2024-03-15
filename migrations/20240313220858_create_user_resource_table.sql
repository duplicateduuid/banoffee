CREATE TYPE resource_status AS ENUM ('completed', 'ongoing', 'bookmarked');

CREATE TYPE rating AS ENUM ('one', 'two', 'three', 'four', 'five');

-- Add migration script here
CREATE TABLE "user_resource" (
    user_id UUID,
    resource_id UUID,
    status resource_status NOT NULL,
    review_rating rating,
    review_comment VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, resource_id),

    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES "user"(id),
    CONSTRAINT fk_resource FOREIGN KEY(resource_id) REFERENCES "resource"(id)
);