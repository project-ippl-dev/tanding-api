CREATE TYPE accomplishment_level as ENUM('region', 'province', 'national', 'international', 'others');

CREATE TABLE accomplishments(
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    title VARCHAR(150) NOT NULL,
    level accomplishment_level NOT NULL,
    ranking VARCHAR(150) NOT NULL,
    category VARCHAR(100) NOT NULL,
    sport VARCHAR(255) NOT NULL,
    description TEXT NULL,
    file_url VARCHAR(255) NOT NULL,
    month SMALLINT NOT NULL,
    year SMALLINT NOT NULL,
    created_at TIMESTAMP NULL DEFAULT NOW(),
    updated_at TIMESTAMP NULL
);

CREATE INDEX idx_accomplishment_user_id ON accomplishments(user_id);