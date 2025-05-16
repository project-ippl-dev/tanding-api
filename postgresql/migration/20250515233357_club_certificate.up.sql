CREATE TABLE club_certificates(
    id UUID PRIMARY KEY,
    club_id UUID NOT NULL,
    event_id UUID NOT NULL,
    reward_as VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);

CREATE INDEX idx_club_certificates ON club_certificates(club_id, event_id);