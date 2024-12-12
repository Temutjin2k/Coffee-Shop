package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"hot-coffee/internal/server"

	_ "github.com/lib/pq"
)

const (
	dbhost   = "db"
	dbport   = 5432
	user     = "latte"
	password = "latte"
	dbname   = "frappuccino"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbhost, dbport, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		logger.Error("Error creating sql.DB", "Error", err)
		os.Exit(1)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		logger.Error("Error connecting to database", "Error", err)
		os.Exit(1)
	}
	query := `
	Insert into menu_items (Name, Description, Price, Field) values
    ('espresso', 'Heavy shot of coffee', 5.99, 'Idk what to write')
	`
	db.Exec(query)
	server.ServerLaunch(db, logger)
}
