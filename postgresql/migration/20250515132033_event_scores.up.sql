CREATE TABLE event_scores(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_bracket_id UUID NOT NULL,
    home_round1 SMALLINT NOT NULL,
    home_round2 SMALLINT NOT NULL,
    home_round3 SMALLINT NOT NULL,
    home_extra SMALLINT NOT NULL,
    home_total SMALLINT NOT NULL,
    away_round1 SMALLINT NOT NULL,
    away_round2 SMALLINT NOT NULL,
    away_round3 SMALLINT NOT NULL,
    away_extra SMALLINT NOT NULL,
    away_total SMALLINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);

CREATE INDEX idx_event_scores ON event_scores(event_bracket_id);