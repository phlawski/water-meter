package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/phlawski/water-meter/internal/db"
	"github.com/phlawski/water-meter/internal/handler"
	_ "modernc.org/sqlite"
)

func main() {
	dsn := os.Getenv("DATABASE_PATH")
	if dsn == "" {
		dsn = "water-meter.db"
	}

	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer conn.Close()

	if err := migrate(conn); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	queries := db.New(conn)
	h := handler.New(queries)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", h.Index)
	r.Post("/readings", h.CreateReading)
	r.Post("/readings/{id}/delete", h.DeleteReading)
	r.Post("/config/price", h.SetPrice)
	r.Post("/config/internet", h.SetInternetContribution)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	addr := ":8080"
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

func migrate(conn *sql.DB) error {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS readings (
			id           INTEGER PRIMARY KEY AUTOINCREMENT,
			read_at      DATE    NOT NULL,
			value_m3     REAL    NOT NULL,
			price_per_m3 REAL    NOT NULL,
			notes        TEXT    NOT NULL DEFAULT ''
		);
		CREATE TABLE IF NOT EXISTS config (
			key   TEXT PRIMARY KEY,
			value TEXT NOT NULL
		);
	`)
	return err
}
