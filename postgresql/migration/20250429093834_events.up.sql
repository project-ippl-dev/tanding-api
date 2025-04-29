CREATE TYPE event_type as ENUM('competition', 'event');
CREATE TYPE remark_type as ENUM('unconfirmed', 'soon', 'open', 'closed', 'ongoing', 'done', 'rejected');

CREATE TABLE events(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    type event_type NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    prize_pool VARCHAR(255) NOT NULL,
    location TEXT NOT NULL,
    province VARCHAR(150) NOT NULL,
    city VARCHAR(255) NOT NULL,
    thumbnail VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    deadline TIMESTAMP NOT NULL,
    sport_id UUID NOT NULL,
    rules TEXT NOT NULL,
    quota INT NOT NULL,
    proposal_link VARCHAR(255) NOT NULL,
    status BOOLEAN NULL,
    remark remark_type NOT NULL DEFAULT 'unconfirmed',
    is_generate BOOLEAN NOT NULL DEFAULT false,
    open TIMESTAMP NOT NULL,
    order_number BIGSERIAL NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_events_user_id ON events(user_id);
CREATE INDEX idx_events_sport_id ON events(sport_id);