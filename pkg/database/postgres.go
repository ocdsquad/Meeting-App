package database

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgresConnection(host string, user string, pass string, dbname string, port string) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, dbname,
	)

	var err error
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(20)                  // Maksimal 20 koneksi aktif
	db.SetMaxIdleConns(5)                   // Maksimal 5 koneksi idle
	db.SetConnMaxLifetime(30 * time.Minute) // Maksimal umur koneksi 30 menit

	log.Println("Database connected")
	return db, nil
}
