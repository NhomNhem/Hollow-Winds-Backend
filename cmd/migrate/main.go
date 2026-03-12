package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	// Database URL from environment or use default
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Try direct connection (port 5432) instead of pooler (port 6543)
		dbURL = "postgresql://postgres.sfkzrqxjylbedwwvdgzo:Truongcutedeptrai123!@db.sfkzrqxjylbedwwvdgzo.supabase.co:5432/postgres"
	}

	// Connect to database
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	log.Println("✅ Connected to database")

	// Read migration file
	migrationSQL, err := os.ReadFile("migrations/004_hollow_wilds_phase1.sql")
	if err != nil {
		log.Fatalf("Failed to read migration file: %v\n", err)
	}

	log.Println("📋 Running migration: 004_hollow_wilds_phase1.sql")

	// Execute migration
	_, err = conn.Exec(context.Background(), string(migrationSQL))
	if err != nil {
		log.Fatalf("Failed to run migration: %v\n", err)
	}

	log.Println("✅ Migration completed successfully!")

	// Verify tables exist
	log.Println("\n🔍 Verifying tables...")

	tables := []string{"players", "player_saves", "player_save_backups", "leaderboard_entries", "analytics_events"}
	for _, table := range tables {
		var exists bool
		err = conn.QueryRow(context.Background(),
			"SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = $1)",
			table,
		).Scan(&exists)

		if err != nil {
			log.Printf("❌ Failed to check table %s: %v\n", table, err)
		} else if exists {
			log.Printf("   ✅ Table '%s' exists\n", table)
		} else {
			log.Printf("   ❌ Table '%s' does NOT exist\n", table)
		}
	}

	fmt.Println("\n🎉 Hollow Wilds Phase 1 schema ready!")
}
