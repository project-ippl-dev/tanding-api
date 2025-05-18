CREATE TABLE class_events(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID NOT NULL,
    class_id UUID NOT NULL,
    price INTEGER NOT NULL,
    bracket_generate BOOLEAN NOT NULL DEFAULT 'false',
    bracket_lock BOOLEAN NOT NULL DEFAULT 'false',
    score_lock BOOLEAN NOT NULL DEFAULT 'false',
    match_index SMALLINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NULL
);

CREATE INDEX idx_class_events_event_id ON class_events(event_id);
CREATE INDEX idx_class_events_class_id ON class_events(class_id);
