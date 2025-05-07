CREATE TABLE privilege_user(
    id BIGSERIAL PRIMARY KEY,
    privilege_id UUID NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMP NULL DEFAULT NOW(),
    updated_at TIMESTAMP NULL
);

CREATE INDEX idx_privilege_user_privilege_id ON privilege_user(privilege_id);
CREATE INDEX idx_privilege_user_user_id ON privilege_user(user_id);