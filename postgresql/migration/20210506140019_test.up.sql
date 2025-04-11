CREATE TABLE test(
    id SERIAL PRIMARY KEY,
    uuid UUID NOT NULL,
    name VARCHAR(30) NOT NULL,
    created_at TIMESTAMP NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);