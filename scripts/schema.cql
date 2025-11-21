CREATE KEYSPACE IF NOT EXISTS userservice
WITH replication = {'class': 'NetworkTopologyStrategy', 'replication_factor': 1};

USE userservice;

CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY,
    first_name text NOT NULL,
    last_name text NOT NULL,
    gender text NOT NULL,
    date_of_birth timestamp NOT NULL,
    phone_number text NOT NULL,
    email text NOT NULL,
    is_blocked boolean DEFAULT false,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS users_by_phone (
    phone_number text PRIMARY KEY,
    user_id uuid NOT NULL,
    created_at timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS users_by_email (
    email text PRIMARY KEY,
    user_id uuid NOT NULL,
    created_at timestamp NOT NULL
);

CREATE INDEX IF NOT EXISTS ON users (email);
CREATE INDEX IF NOT EXISTS ON users (phone_number);
CREATE INDEX IF NOT EXISTS ON users (is_blocked);
