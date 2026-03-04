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
	migrationSQL, err := os.ReadFile("migrations/003_admin_system.sql")
	if err != nil {
		log.Fatalf("Failed to read migration file: %v\n", err)
	}

	log.Println("📋 Running migration: 003_admin_system.sql")

	// Execute migration
	_, err = conn.Exec(context.Background(), string(migrationSQL))
	if err != nil {
		log.Fatalf("Failed to run migration: %v\n", err)
	}

	log.Println("✅ Migration completed successfully!")

	// Verify tables exist
	log.Println("\n🔍 Verifying tables...")
	
	tables := []string{"users", "admin_actions", "user_bans"}
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

	// Check is_admin column
	var columnExists bool
	err = conn.QueryRow(context.Background(),
		"SELECT EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'is_admin')",
	).Scan(&columnExists)
	
	if err != nil {
		log.Printf("❌ Failed to check is_admin column: %v\n", err)
	} else if columnExists {
		log.Println("   ✅ Column 'users.is_admin' exists")
	} else {
		log.Println("   ❌ Column 'users.is_admin' does NOT exist")
	}

	fmt.Println("\n🎉 Admin system schema ready!")
}
