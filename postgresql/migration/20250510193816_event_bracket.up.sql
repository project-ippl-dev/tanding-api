CREATE TYPE bracket_type AS ENUM('battle', 'bye');

--match order, match order in a class event
--match index is type of match like final, semi final, etc.
--event turn is 'partai pertandingan' in indonesia
--next match id, next match after this match
CREATE TABLE event_brackets(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID NOT NULL,
    class_event_id UUID NOT NULL,
    event_turn SMALLINT NOT NULL,
    match_index SMALLINT NOT NULL,
    match_order SMALLINT NOT NULL,
    next_match_id UUID NOT NULL,
    status bracket_type NOT NULL,
    is_active SMALLINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NULL
);

CREATE INDEX idx_event_brackets ON event_brackets(event_id, class_event_id, next_match_id, match_order, event_turn);