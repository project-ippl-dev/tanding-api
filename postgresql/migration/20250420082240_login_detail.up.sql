CREATE TABLE login_details(
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);

CREATE INDEX idx_login_details ON login_details(user_id);