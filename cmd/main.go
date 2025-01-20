package main

import (
	"database/sql"
	"fmt"
	"hot-coffee/internal/server"
	"log/slog"
	"os"
	"time"

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
	// Если ошибка подключения, пробуем несколько раз
	if err != nil {
		logger.Error("Error connecting to database: %v", err)

		// Пробуем подключиться до 5 раз с интервалом в 5 секунд
		for i := 0; i < 5; i++ {
			logger.Info(fmt.Sprintf("Retrying to connect... Attempt #%d", i+1))
			time.Sleep(3 * time.Second)

			// Проверяем снова
			if err = db.Ping(); err == nil {
				logger.Info("Successfully connected to the database")
				break
			} else if i == 4 {
				logger.Error("Failed to connect to the database after 5 attempts. Exiting...")
				os.Exit(1)
			}
		}
	}

	server.ServerLaunch(db, logger)
}
