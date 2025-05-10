
CREATE TYPE event_role as ENUM ('owner', 'reviewer', 'contributor', 'admin');
CREATE TABLE event_privileges(
    id BIGSERIAL PRIMARY KEY,
    event_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role event_role NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NULL
);

CREATE INDEX idx_event_privileges_event_id ON event_privileges(event_id);
CREATE INDEX idx_event_privileges_user_id ON event_privileges(user_id);