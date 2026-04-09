// database/db.go
package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/lib/pq" // The underscore means we import it just for its side-effects (registering the Postgres driver)
)

// DB is a global variable. It holds our database connection pool.
// It starts with a capital 'D' so other folders (like our controllers) can access it.
var DB *sql.DB

func Connect() {
	// 1. Build the connection string using our environment variables
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	// 2. Open the database connection
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("❌ Failed to open database: %v", err)
	}

	// 3. Ping the database to ensure it's actually reachable
	err = DB.Ping()
	if err != nil {
		log.Fatalf("❌ Database is not responding: %v", err)
	}

	log.Println("✅ Successfully connected to PostgreSQL!")
}

func RunMigrations() {
	driver, err := postgres.WithInstance(DB, &postgres.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to initialize migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		log.Fatalf("❌ Failed to initialize migrations: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("ℹ️ Database migrations are already up to date")
		return
	}

	log.Println("✅ Database migrations applied successfully")
}
