package main

// @title MusicLibrary Swagger API
// @version 1.0

import (
	_ "MusicLibrary/docs"
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"os"
	"time"
)

type Song struct {
	ID          string    `json:"id"`
	Group       string    `json:"group"`
	Name        string    `json:"name"`
	ReleaseDate time.Time `json:"releaseDate"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
}

var log *logrus.Logger
var db *sql.DB

func initDB() (*sql.DB, error) {
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	log.Info("Successfully connected to database")
	return db, nil
}

func migrate() error {
	query := `
        CREATE TABLE IF NOT EXISTS musicdb (
            id SERIAL PRIMARY KEY,
            group_name TEXT NOT NULL,
            song_name TEXT NOT NULL,
            release_date DATE,
            text TEXT,
            link TEXT
        );
    `
	_, err := db.Exec(query)
	return err
}

func initLogger() *logrus.Logger {
	log := logrus.New()

	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	logLevel := "DEBUG"
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		log.Warn("Invalid LOG_LEVEL, defaulting to INFO")
		level = logrus.InfoLevel
	}

	log.SetLevel(level)

	return log
}

func main() {
	log = initLogger()

	if err := godotenv.Load("config.env"); err != nil {
		log.WithError(err).Fatal("Error loading .env file")
	}

	var err error

	db, err = initDB()

	if err != nil {
		log.WithError(err).Fatal("Fail: connect to DB")
	}
	defer db.Close()

	if err := migrate(); err != nil {
		log.WithError(err).Fatal("Failed to run migrations")
	}
	log.Info("Database migration completed")

	r := mux.NewRouter()

	r.HandleFunc("/songs", getSongs).Methods("GET")
	r.HandleFunc("/songs/{id}", getTextFromSong).Methods("GET")
	r.HandleFunc("/songs", addSong).Methods("POST")
	r.HandleFunc("/songs/{id}", editSong).Methods("PUT")
	r.HandleFunc("/songs/{id}", deleteSong).Methods("DELETE")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Info("Server started on port :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
