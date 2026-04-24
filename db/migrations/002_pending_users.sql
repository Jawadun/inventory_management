-- User Registration Approval System
-- Migration: 002_pending_users.sql
-- Run order: 2

CREATE TABLE IF NOT EXISTS pending_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    department VARCHAR(255),
    employee_id VARCHAR(50),
    phone VARCHAR(20),
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_pending_users_username ON pending_users(username);
CREATE INDEX IF NOT EXISTS idx_pending_users_status ON pending_users(status);