CREATE TABLE event_participants(
  id BIGSERIAL PRIMARY KEY,
  event_registration_id uuid NOT NULL,
  user_id UUID NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NULL
);

CREATE INDEX idx_event_participants_event_registration_id ON event_participants(event_registration_id);
CREATE INDEX idx_event_participants_user_id ON event_participants(user_id);