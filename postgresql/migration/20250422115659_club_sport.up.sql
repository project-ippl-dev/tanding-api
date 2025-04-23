CREATE TABLE club_sport(
    id BIGSERIAL PRIMARY KEY,
    club_id UUID NOT NULL,
    sport_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NULL
);

CREATE INDEX idx_club_sport_club_id ON club_sport(club_id);
CREATE INDEX idx_club_sport_sport_id ON club_sport(sport_id);