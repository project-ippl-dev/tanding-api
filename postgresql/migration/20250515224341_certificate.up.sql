CREATE TABLE certificates(
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    event_id UUID NOT NULL,
    class_event_id UUID NOT NULL,
    reward_as VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);

CREATE INDEX idx_certificates ON certificates(user_id, event_id, class_event_id);