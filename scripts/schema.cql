-- Create keyspace
CREATE KEYSPACE IF NOT EXISTS userservice 
WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

USE userservice;

-- Main users table (query by ID)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    gender TEXT NOT NULL,
    date_of_birth TIMESTAMP NOT NULL,
    phone_number TEXT NOT NULL,
    email TEXT NOT NULL,
    is_blocked BOOLEAN DEFAULT false,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
) WITH CLUSTERING ORDER BY (id ASC)
AND default_time_to_live = 0
AND comment = 'Primary user data storage';

-- Lookup table for phone number
CREATE TABLE IF NOT EXISTS users_by_phone (
    phone_number TEXT PRIMARY KEY,
    user_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL
) WITH comment = 'Secondary index for phone number lookup';

-- Lookup table for email
CREATE TABLE IF NOT EXISTS users_by_email (
    email TEXT PRIMARY KEY,
    user_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL
) WITH comment = 'Secondary index for email lookup';

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS users_email_idx ON users (email);
CREATE INDEX IF NOT EXISTS users_phone_idx ON users (phone_number);
CREATE INDEX IF NOT EXISTS users_is_blocked_idx ON users (is_blocked);