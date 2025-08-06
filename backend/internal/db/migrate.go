package db

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/cockroachdb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(reset bool, databaseURL string) {
	log.Println("🔧 Ejecutando migraciones...")

	if databaseURL == "" {
		log.Fatal("DATABASE_URL no está definido")
	}

	m, err := migrate.New(
		"file://internal/db/migrations",
		"cockroachdb://"+databaseURL,
	)
	if err != nil {
		log.Fatalf("Error creando migrador: %v", err)
	}

	if reset {
		log.Println("Ejecutando down de migraciones...")
		if err := m.Down(); err != nil {
			log.Fatalf("Error al hacer drop: %v", err)
		}
	}

	log.Println("Ejecutando migraciones...")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Error ejecutando migraciones: %v", err)
	}

	log.Println("✅ Migraciones ejecutadas con éxito")
}
