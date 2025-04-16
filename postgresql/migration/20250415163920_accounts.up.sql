CREATE TYPE account_type as ENUM ('manual', 'facebook', 'google');

CREATE TABLE accounts(
    id BIGSERIAL PRIMARY KEY,
    type account_type NOT NULL,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL,
    status BOOL NOT NULL DEFAULT '0',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NULL
);

CREATE INDEX idx_accounts_user_id ON accounts(user_id);