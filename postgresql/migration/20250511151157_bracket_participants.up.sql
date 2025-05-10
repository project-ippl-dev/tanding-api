CREATE TYPE participant_type AS ENUM('home', 'away');

CREATE TABLE bracket_participants(
    id BIGSERIAL PRIMARY KEY,
    event_bracket_id UUID NOT NULL,
    event_registration_id UUID,
    type participant_type NOT NULL,
    is_bye BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);

CREATE INDEX idx_bracket_participants ON bracket_participants(event_registration_id, event_bracket_id);