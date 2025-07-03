package infra

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitPostgres() {
	connStr := os.Getenv("DATABASE_URL")
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("❌ Failed to connect to DB:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("❌ DB ping failed:", err)
	}

	fmt.Println("✅ Connected to PostgreSQL")
}
