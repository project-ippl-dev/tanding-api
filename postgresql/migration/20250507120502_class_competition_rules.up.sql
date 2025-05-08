CREATE TABLE class_competition_rules(
    id BIGSERIAL PRIMARY KEY NOT NULL,
    name VARCHAR(150) NOT NULL,
    male SMALLINT NOT NULL,
    female SMALLINT NOT NULL,
    total SMALLINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);