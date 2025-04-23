CREATE TABLE club_participants(
                                  id BIGSERIAL PRIMARY KEY,
                                  club_id UUID NOT NULL,
                                  user_id UUID NOT NULL,
                                  sport_id UUID NOT NULL,
                                  club_approval BOOLEAN NULL,
                                  user_approval BOOLEAN NULL,
                                  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                  updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_club_participants_club_id ON club_participants(club_id);
CREATE INDEX idx_club_participants_user_id ON club_participants(user_id);
CREATE INDEX idx_club_participants_sport_id ON club_participants(sport_id);