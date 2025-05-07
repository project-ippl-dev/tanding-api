CREATE TABLE clubs(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(150) NOT NULL,
    logo VARCHAR(255) NOT NULL,
    phone VARCHAR(15) NOT NULL,
    short_name VARCHAR(5) NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_clubs_user_id ON clubs(user_id);