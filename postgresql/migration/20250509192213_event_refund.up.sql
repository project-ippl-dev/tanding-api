CREATE TYPE event_refund_status AS ENUM('waiting', 'approved', 'rejected');

CREATE TABLE event_refunds(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID NOT NULL,
    event_payment_receipt_id UUID NOT NULL,
    note TEXT NOT NULL,
    status event_refund_status NOT NULL DEFAULT 'waiting',
    admin_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NULL
);

CREATE INDEX idx_event_refunds_event_id ON event_refunds(event_id);
CREATE INDEX idx_event_refunds_event_payment_receipt_id ON event_refunds(event_payment_receipt_id);
CREATE INDEX idx_event_refunds_admin_id ON event_refunds(admin_id);