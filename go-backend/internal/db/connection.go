package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func NewCockroachDB() *sql.DB {
	dsn := os.Getenv("DATABASE_URI")
	if dsn == "" {
		log.Fatal("DATABASE_URI no está definido en el entorno")
	}

	db, err := sql.Open("postgres", "postgres://"+dsn)
	if err != nil {
		log.Fatalf("No se pudo abrir la conexión: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("No se pudo hacer ping a la base de datos: %v", err)
	}

	return db
}
