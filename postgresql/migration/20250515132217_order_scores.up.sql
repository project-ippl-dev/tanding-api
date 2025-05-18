CREATE TABLE order_scores(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_bracket_id UUID NOT NULL,
    round1 SMALLINT NOT NULL,
    round2 SMALLINT NOT NULL,
    round3 SMALLINT NOT NULL,
    extra SMALLINT NOT NULL,
    total SMALLINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);

CREATE INDEX idx_order_scores ON order_scores(order_bracket_id);