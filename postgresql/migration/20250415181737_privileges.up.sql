CREATE TYPE privilege_type as ENUM ('main', 'competition');

CREATE TABLE privileges(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(150) NOT NULL,
    type privilege_type NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NULL DEFAULT NOW(),
    updated_at TIMESTAMP NULL
);