package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"github.com/viteant/stockinsight/internal/db"
	"github.com/viteant/stockinsight/internal/stock/interfaces"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  Archivo .env no encontrado, usando variables del sistema.")
	}

	databaseURL := os.Getenv("DATABASE_URI")

	app := &cli.App{
		Name:  "stock-app",
		Usage: "Gestión de aplicaciones de StockInsight",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "migrate",
				Usage: "Ejecutar las migraciones",
			},
			&cli.BoolFlag{
				Name:  "reset",
				Usage: "Resetear la base antes de migrar (solo con --migrate)",
			},
			&cli.BoolFlag{
				Name:  "serve",
				Usage: "Iniciar el servidor",
			},
			&cli.BoolFlag{
				Name:  "sync",
				Usage: "Sincronizar datos de stocks con la API",
			},
		},
		Action: func(c *cli.Context) error {
			switch {
			case c.Bool("migrate"):
				migrate(c.Bool("reset"), databaseURL)
			case c.Bool("sync"):
				syncData()
			case c.Bool("serve"), c.NumFlags() == 0:
				startServer()
			default:
				log.Println("Ninguna acción válida. Usa --help para ver opciones.")
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func migrate(reset bool, databaseURL string) {
	db.RunMigrations(reset, databaseURL)
}

func startServer() {
	log.Println("🚀 Iniciando servidor...")
	// Aquí levantarás tu API más adelante
}

func syncData() {
	log.Println("🔄 Sincronizando datos de stocks con la API...")
	interfaces.RunStockSync()
	log.Println("🔄 Sincronización de stocks completada.")
}
