CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE role as ENUM('admin', 'user');

CREATE TABLE users(
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name VARCHAR(150) NOT NULL,
  born_at VARCHAR(150) NOT NULL DEFAULT '',
  born_on DATE,
  identity_number VARCHAR(16) NOT NULL DEFAULT '',
  phone VARCHAR(18) NOT NULL DEFAULT '',
  photo VARCHAR(255) NOT NULL DEFAULT '',
  role role NOT NULL DEFAULT 'user',
  gender VARCHAR(6) NOT NULL DEFAULT '',
  about TEXT NOT NULL DEFAULT '',
  status BOOL NOT NULL DEFAULT '1',
  can_participate BOOL NOT NULL DEFAULT '1',
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP
);