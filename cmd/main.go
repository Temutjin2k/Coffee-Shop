package main

import (
	"database/sql"
	"fmt"
	"hot-coffee/internal/server"
	"log/slog"
	"os"

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

	server.ServerLaunch(db, logger)
}
