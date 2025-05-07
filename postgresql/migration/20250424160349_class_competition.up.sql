CREATE TYPE class_type AS ENUM('default', 'custom');
-- single =  Single Elimination, order = Giliran Urutan
CREATE TYPE match_type AS ENUM('single', 'order');

CREATE TABLE classes(
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  sport_id UUID NOT NULL,
  name VARCHAR(150) NOT NULL,
  class_competition_rule_id BIGINT NOT NULL,
  match_type match_type DEFAULT 'single' NOT NULL,
  type class_type NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NULL
);

CREATE INDEX idx_classes_sport_id ON classes(sport_id);
CREATE INDEX idx_classes_class_competition_rule_id ON classes(class_competition_rule_id);