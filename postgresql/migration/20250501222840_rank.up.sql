CREATE TABLE ranks(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    club_id UUID NOT NULL,
    event_id UUID NOT NULL,
    class_event_id UUID NOT NULL,
    event_registration_id UUID NOT NULL,
    sport_id UUID NOT NULL,
    rank SMALLINT NOT NULL,
    point INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);

CREATE INDEX idx_club_rankings ON ranks(club_id, event_id, class_event_id, event_registration_id, sport_id);