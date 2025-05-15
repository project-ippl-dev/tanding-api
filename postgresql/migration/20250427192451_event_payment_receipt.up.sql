CREATE TYPE event_receipt_status AS ENUM('waiting', 'approved', 'rejected', 'refund');

CREATE TABLE event_payment_receipts(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID NOT NULL,
    unique_number SMALLINT NOT NULL,
    payment_link VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL,
    admin_id UUID NOT NULL,
    total INTEGER NOT NULL,
    status event_receipt_status NOT NULL NOT NULL DEFAULT 'waiting',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_event_payment_receipts_user_id ON event_payment_receipts(user_id);
CREATE INDEX idx_event_payment_receipts_event_id ON event_payment_receipts(event_id);