-- Migration: Platform Schema V1
-- Description: establishes multi-game support and tiered role system
-- Date: 2026-03-24

-- ============================================================================
-- 1. Create system_role type
-- ============================================================================
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'system_role') THEN
        CREATE TYPE system_role AS ENUM (
            'super_admin',
            'admin',
            'game_manager',
            'support',
            'user',
            'beta_tester'
        );
    END IF;
END $$;

-- ============================================================================
-- 2. Update users table with system_role
-- ============================================================================
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS system_role system_role DEFAULT 'user';

-- Migrate existing is_admin flag to system_role
UPDATE users 
SET system_role = 'admin' 
WHERE is_admin = true;

-- ============================================================================
-- 3. Create games table
-- ============================================================================
CREATE TABLE IF NOT EXISTS games (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug TEXT UNIQUE NOT NULL, -- e.g. 'hollow-wilds'
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Register Hollow Wilds as the first official game
INSERT INTO games (slug, name, description)
VALUES ('hollow-wilds', 'Hollow Wilds', 'Our flagship survival game')
ON CONFLICT (slug) DO NOTHING;

-- ============================================================================
-- 4. Create game_access table (User <-> Game mapping)
-- ============================================================================
CREATE TABLE IF NOT EXISTS game_access (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    role TEXT NOT NULL DEFAULT 'player', -- 'player', 'tester', 'developer'
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, game_id)
);

-- Grant all existing users access to Hollow Wilds as players
INSERT INTO game_access (user_id, game_id, role)
SELECT u.id, g.id, 'player'
FROM users u, games g
WHERE g.slug = 'hollow-wilds'
ON CONFLICT DO NOTHING;

-- ============================================================================
-- 5. Comments
-- ============================================================================
COMMENT ON COLUMN users.system_role IS 'Global platform-level role';
COMMENT ON TABLE games IS 'Registry of all games in the studio platform';
COMMENT ON TABLE game_access IS 'Tracks which users have access to which games and their role within that game';
