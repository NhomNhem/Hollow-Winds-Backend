-- Migration: Add Admin System
-- Description: Adds admin role system, audit logging, and user bans
-- Date: 2026-03-04

-- ============================================================================
-- 1. Add is_admin column to users table
-- ============================================================================
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS is_admin BOOLEAN DEFAULT false;

-- Create index for admin queries (partial index for better performance)
CREATE INDEX IF NOT EXISTS idx_users_is_admin 
ON users(is_admin) 
WHERE is_admin = true;

-- ============================================================================
-- 2. Create admin_actions audit table
-- ============================================================================
CREATE TABLE IF NOT EXISTS admin_actions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    admin_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    action_type VARCHAR(50) NOT NULL, -- 'adjust_gold', 'ban_user', 'unban_user', 'grant_talent', etc.
    target_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    details JSONB NOT NULL DEFAULT '{}', -- Action-specific data (reason, amount, old_value, new_value, etc.)
    ip_address VARCHAR(45), -- IPv4 or IPv6
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for audit log queries
CREATE INDEX IF NOT EXISTS idx_admin_actions_admin_id ON admin_actions(admin_id);
CREATE INDEX IF NOT EXISTS idx_admin_actions_target_user_id ON admin_actions(target_user_id);
CREATE INDEX IF NOT EXISTS idx_admin_actions_created_at ON admin_actions(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_admin_actions_action_type ON admin_actions(action_type);

-- ============================================================================
-- 3. Create user_bans table
-- ============================================================================
CREATE TABLE IF NOT EXISTS user_bans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    banned_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reason TEXT NOT NULL,
    banned_until TIMESTAMP, -- NULL = permanent ban
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    unbanned_at TIMESTAMP,
    unbanned_by UUID REFERENCES users(id) ON DELETE SET NULL,
    unban_reason TEXT
);

-- Indexes for ban lookups
CREATE INDEX IF NOT EXISTS idx_user_bans_user_id ON user_bans(user_id);
CREATE INDEX IF NOT EXISTS idx_user_bans_is_active ON user_bans(is_active, user_id) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_user_bans_created_at ON user_bans(created_at DESC);

-- ============================================================================
-- Comments for documentation
-- ============================================================================
COMMENT ON TABLE admin_actions IS 'Audit log for all admin operations';
COMMENT ON COLUMN admin_actions.action_type IS 'Type of admin action performed';
COMMENT ON COLUMN admin_actions.details IS 'JSON object with action-specific data (reason, amounts, etc.)';

COMMENT ON TABLE user_bans IS 'User ban records with history';
COMMENT ON COLUMN user_bans.banned_until IS 'NULL means permanent ban';
COMMENT ON COLUMN user_bans.is_active IS 'True if ban is currently active';
