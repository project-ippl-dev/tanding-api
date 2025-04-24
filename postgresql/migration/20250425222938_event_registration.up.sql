CREATE TYPE event_registration_status AS ENUM('pending', 'canceled', 'waiting', 'approved', 'rejected');
CREATE TABLE event_registrations(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID NOT NULL,
    class_event_id UUID NOT NULL,
    club_id UUID NOT NULL,
    status event_registration_status NOT NULL DEFAULT 'pending',
    event_payment_receipt_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NULL
);

CREATE INDEX idx_event_registrations_event_id ON event_registrations(event_id);
CREATE INDEX idx_event_registrations_class_id ON event_registrations(class_event_id);
CREATE INDEX idx_event_registrations_club_id ON event_registrations(club_id);
CREATE INDEX idx_event_registrations_event_payment_receipt_id ON event_registrations(event_payment_receipt_id);