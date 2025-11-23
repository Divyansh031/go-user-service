CREATE KEYSPACE IF NOT EXISTS userservice
WITH replication = {'class': 'NetworkTopologyStrategy', 'replication_factor': 1};

USE userservice;

CREATE TABLE IF NOT EXISTS users (
    id text PRIMARY KEY,
    first_name text,
    last_name text,
    gender text,
    date_of_birth timestamp,
    phone_number text,
    email text,
    is_blocked boolean,
    created_at timestamp,
    updated_at timestamp
);

CREATE TABLE IF NOT EXISTS users_by_phone (
    phone_number text PRIMARY KEY,
    user_id text,
    created_at timestamp 
);

CREATE TABLE IF NOT EXISTS users_by_email (
    email text PRIMARY KEY,
    user_id text,
    created_at timestamp 
);

CREATE INDEX IF NOT EXISTS ON users (email);
CREATE INDEX IF NOT EXISTS ON users (phone_number);
CREATE INDEX IF NOT EXISTS ON users (is_blocked);
