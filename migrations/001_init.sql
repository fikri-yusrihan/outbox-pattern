-- migrations/001_init.sql

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS events;
DROP INDEX IF EXISTS idx_outbox_published;

-- This script initializes the database schema for the application.
-- It creates the necessary tables and indexes.
-- Ensure that the database is in a clean state before applying the migration.

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT NOT NULL,
);

CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type TEXT NOT NULL,
    payload TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    published_at TIMESTAMP,
    published BOOLEAN DEFAULT false
);

CREATE INDEX idx_outbox_published ON events (published, created_at);