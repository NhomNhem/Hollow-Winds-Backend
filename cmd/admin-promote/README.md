# Admin Promotion Tool

Quick tool to promote users to admin role.

## Usage

### On Fly.io (Recommended):

```bash
# SSH into your Fly.io app
fly ssh console -a gamefeel-backend

# Run the tool (promotes first user)
./cmd/admin-promote/admin-promote

# Or promote specific user
./cmd/admin-promote/admin-promote PLAYFAB_ID
```

### Locally:

```bash
# Set DATABASE_URL
export DATABASE_URL="your-database-url"

# Build and run
go run ./cmd/admin-promote

# Or with specific PlayFab ID
go run ./cmd/admin-promote YOUR_PLAYFAB_ID
```

### Via Supabase SQL Editor (Easiest):

```sql
-- List all users
SELECT playfab_id, username, is_admin FROM users;

-- Promote user to admin
UPDATE users SET is_admin = true WHERE playfab_id = 'YOUR_PLAYFAB_ID';

-- Verify
SELECT playfab_id, username FROM users WHERE is_admin = true;
```

## What it does:

1. Connects to database
2. If no argument: Lists all users and promotes the first one
3. If PlayFab ID provided: Promotes that specific user
4. Shows all current admin users

## Output:

```
✅ Connected to database
🔍 Listing all users in database:

1. PlayFab ID: ABC123 | Username: Player1
2. PlayFab ID: XYZ789 | Username: Player2

🎯 Promoting first user to admin: ABC123
✅ First user promoted to admin!

📊 Current admin users:
   👑 Player1 (ABC123)
```
