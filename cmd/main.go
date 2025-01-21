package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	"hot-coffee/internal/server"

	_ "github.com/lib/pq"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

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
