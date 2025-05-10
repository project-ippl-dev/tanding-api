CREATE TABLE order_brackets(
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  event_id UUID NOT NULL,
  class_event_id UUID NOT NULL,
  club_id UUID NOT NULL,
  event_registration_id UUID NOT NULL,
  order_by SMALLINT NOT NULL DEFAULT 0,
  rank SMALLINT NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NULL
);

CREATE INDEX idx_order_brackets_event_id ON order_brackets(event_id);
CREATE INDEX idx_order_brackets_class_event_id ON order_brackets(class_event_id);
CREATE INDEX idx_order_brackets_event_registration_id ON order_brackets(event_registration_id);
CREATE INDEX idx_order_brackets_club_id ON order_brackets(club_id);