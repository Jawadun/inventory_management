-- IICT Inventory Management System - Database Schema
-- Migration: 001_initial_schema.sql
-- Run order: 1

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================
-- USERS & AUTHENTICATION
-- ============================================

-- Roles table
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert default roles
INSERT INTO roles (name, description) VALUES 
    ('admin', 'Full system access'),
    ('user', 'Can request and manage personal items'),
    ('viewer', 'Read-only public access');

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role_id INTEGER NOT NULL REFERENCES roles(id) DEFAULT 2,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    department VARCHAR(255),
    employee_id VARCHAR(50) UNIQUE,
    phone VARCHAR(20),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_role ON users(role_id);
CREATE INDEX idx_users_department ON users(department);

-- Session tokens for refresh
CREATE TABLE auth_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(500) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_auth_tokens_user ON auth_tokens(user_id);
CREATE INDEX idx_auth_tokens_token ON auth_tokens(token);

-- ============================================
-- INVENTORY CORE
-- ============================================

-- Categories
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    parent_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_categories_parent ON categories(parent_id);

-- Suppliers/Companies
CREATE TABLE suppliers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    contact_person VARCHAR(255),
    phone VARCHAR(20),
    email VARCHAR(255),
    address TEXT,
    notes TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_suppliers_name ON suppliers(name);

-- Items (inventory items)
CREATE TABLE items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    supplier_id UUID REFERENCES suppliers(id) ON DELETE SET NULL,
    sku VARCHAR(100) UNIQUE,
    barcode VARCHAR(100) UNIQUE,
    description TEXT,
    quantity INTEGER NOT NULL DEFAULT 0,
    min_quantity INTEGER DEFAULT 5,
    unit VARCHAR(50) DEFAULT 'pcs',
    location VARCHAR(255),
    storage_location VARCHAR(255),
    purchase_date DATE,
    purchase_price DECIMAL(12, 2),
    warranty_months INTEGER,
    status VARCHAR(50) DEFAULT 'available',
    image_url VARCHAR(500),
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_quantity CHECK (quantity >= 0),
    CONSTRAINT chk_min_quantity CHECK (min_quantity >= 0)
);

CREATE INDEX idx_items_category ON items(category_id);
CREATE INDEX idx_items_supplier ON items(supplier_id);
CREATE INDEX idx_items_sku ON items(sku);
CREATE INDEX idx_items_barcode ON items(barcode);
CREATE INDEX idx_items_status ON items(status);
CREATE INDEX idx_items_location ON items(location);

-- Item history (tracking all changes)
CREATE TABLE item_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    item_id UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    quantity_change INTEGER NOT NULL,
    previous_quantity INTEGER NOT NULL,
    new_quantity INTEGER NOT NULL,
    change_type VARCHAR(50) NOT NULL,
    reason TEXT,
    changed_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_item_history_item ON item_history(item_id);
CREATE INDEX idx_item_history_date ON item_history(created_at);

-- ============================================
-- REQUESTS & TRANSACTIONS
-- ============================================

-- Item request types
CREATE TABLE request_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description VARCHAR(255)
);

INSERT INTO request_types (name, description) VALUES 
    ('classroom', 'Issued to classroom'),
    ('lab', 'Issued to lab'),
    ('teachers_room', 'Issued to teachers room'),
    ('personal', 'Personal use');
CREATE TABLE item_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    item_id UUID NOT NULL REFERENCES items(id),
    request_type_id INTEGER REFERENCES request_types(id),
    quantity INTEGER NOT NULL DEFAULT 1,
    status VARCHAR(50) DEFAULT 'pending',
    reason TEXT,
    requested_by UUID REFERENCES users(id),
    requested_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    reviewed_by UUID REFERENCES users(id),
    reviewed_at TIMESTAMP,
    rejection_reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_request_quantity CHECK (quantity > 0)
);

CREATE INDEX idx_item_requests_user ON item_requests(user_id);
CREATE INDEX idx_item_requests_item ON item_requests(item_id);
CREATE INDEX idx_item_requests_status ON item_requests(status);
CREATE INDEX idx_item_requests_requested_by ON item_requests(requested_by);

-- Issue records (items issued to users)
CREATE TABLE issue_records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    request_id UUID REFERENCES item_requests(id) ON DELETE SET NULL,
    item_id UUID NOT NULL REFERENCES items(id),
    recipient_id UUID NOT NULL REFERENCES users(id),
    issued_by UUID REFERENCES users(id),
    quantity INTEGER NOT NULL DEFAULT 1,
    issue_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    due_date DATE,
    return_date DATE,
    actual_return_date DATE,
    return_condition VARCHAR(50),
    return_remarks TEXT,
    notes TEXT,
    status VARCHAR(50) DEFAULT 'issued',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_issue_quantity CHECK (quantity > 0)
);

CREATE INDEX idx_issue_records_item ON issue_records(item_id);
CREATE INDEX idx_issue_records_recipient ON issue_records(recipient_id);
CREATE INDEX idx_issue_records_status ON issue_records(status);
CREATE INDEX idx_issue_records_due_date ON issue_records(due_date);

-- Returns table (optional extension of issue_records)
CREATE TABLE returns (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    issue_id UUID NOT NULL REFERENCES issue_records(id) ON DELETE CASCADE,
    return_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    condition VARCHAR(100),
    remarks TEXT,
    checked_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- NOTICE BOARD
-- ============================================

CREATE TABLE notices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(500) NOT NULL,
    content TEXT NOT NULL,
    posted_by UUID REFERENCES users(id),
    is_pinned BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    priority INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notices_pinned ON notices(is_pinned) WHERE is_pinned = TRUE;
CREATE INDEX idx_notices_active ON notices(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_notices_date ON notices(created_at);

-- ============================================
-- STATS & PUBLIC VIEW
-- ============================================

-- Public stats (for unauthenticated users - cached/aggregated)
CREATE TABLE public_stats (
    id SERIAL PRIMARY KEY,
    total_items INTEGER DEFAULT 0,
    total_categories INTEGER DEFAULT 0,
    total_suppliers INTEGER DEFAULT 0,
    total_issued INTEGER DEFAULT 0,
    total_available INTEGER DEFAULT 0,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add trigger to update public_stats when items change (optional)
-- Will be implemented via application logic

-- ============================================
-- AUDIT LOG
-- ============================================

CREATE TABLE audit_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    action VARCHAR(100) NOT NULL,
    table_name VARCHAR(100),
    record_id UUID,
    changes JSONB,
    ip_address VARCHAR(50),
    user_agent VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_log_user ON audit_log(user_id);
CREATE INDEX idx_audit_log_table ON audit_log(table_name);
CREATE INDEX idx_audit_log_date ON audit_log(created_at);

-- ============================================
-- FUNCTIONS & TRIGGERS
-- ============================================

-- Function to auto-update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger for users
CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Trigger for items
CREATE TRIGGER update_items_updated_at 
    BEFORE UPDATE ON items 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Trigger for categories
CREATE TRIGGER update_categories_updated_at 
    BEFORE UPDATE ON categories 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Trigger for suppliers
CREATE TRIGGER update_suppliers_updated_at 
    BEFORE UPDATE ON suppliers 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Trigger for notices
CREATE TRIGGER update_notices_updated_at 
    BEFORE UPDATE ON notices 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Trigger for issue_records
CREATE TRIGGER update_issue_records_updated_at 
    BEFORE UPDATE ON issue_records 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- DEFAULT DATA
-- ============================================

-- Sample admin user (password: admin123 - MUST CHANGE IN PRODUCTION)
-- Using bcrypt hash for 'admin123'
INSERT INTO users (username, password_hash, role_id, full_name, email, department) VALUES 
    ('admin', '$2a$10$YourBcryptHashHereMustBeChanged', 1, 'System Administrator', 'admin@sust.edu.bd', 'IICT');

-- Sample categories
INSERT INTO categories (name, description) VALUES 
    ('Electronics', 'Electronic items and devices'),
    ('Furniture', 'Office and classroom furniture'),
    ('Stationery', 'Office stationery and supplies'),
    ('Lab Equipment', 'Laboratory instruments and tools'),
    ('Computer Accessories', 'Computer peripherals and accessories');

-- Initialize public stats
INSERT INTO public_stats (total_items, total_categories, total_suppliers, total_issued, total_available) VALUES 
    (0, 5, 0, 0, 0);