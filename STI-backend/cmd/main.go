package main

import (
	"errors"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	_ "github.com/viteant/stockinsight/docs"
	"github.com/viteant/stockinsight/internal/api"
	"github.com/viteant/stockinsight/internal/db"
	"github.com/viteant/stockinsight/internal/db/seeds/finances"
	"github.com/viteant/stockinsight/internal/db/seeds/stocks"
	financeinterfaces "github.com/viteant/stockinsight/internal/finance/interfaces"
	stockinterfaces "github.com/viteant/stockinsight/internal/stock/interfaces"
)

// @title StockInsight API
// @version 1.0
// @description API de acciones y recomendaciones
// @BasePath /api
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
			&cli.StringFlag{
				Name:  "export",
				Usage: "Exportar los datos de stocks a un archivo JSON",
			},
			&cli.StringFlag{
				Name:  "table",
				Usage: "Usa la tabla para la exportación",
			},
			&cli.StringFlag{
				Name:  "import",
				Usage: "Importar datos desde un archivo JSON",
			},
			&cli.BoolFlag{
				Name:  "update-finance",
				Usage: "Actualiza datos históricos de Yahoo Finance para todos los tickers",
			},
		},
		Action: func(c *cli.Context) error {
			if c.Bool("migrate") {
				migrate(c.Bool("reset"), databaseURL)
			} else if c.Bool("sync") {
				syncData()
			} else if c.Bool("serve") || c.NumFlags() == 0 {
				startServer()
			} else if c.Bool("update-finance") || c.NumFlags() == 0 {
				updateFinance()
			} else if path := c.String("export"); path != "" {
				if table := c.String("table"); table != "" {
					exportData(path, table)
				}
			} else if path := c.String("import"); path != "" {
				if table := c.String("table"); table != "" {
					importData(path, table)
				}
			} else {
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
	log.Println("🚀 Iniciando servidor en http://localhost:8080")

	dbConn := db.NewCockroachDB()
	defer dbConn.Close()

	app := fiber.New()

	api.RegisterRoutes(app, dbConn)
	app.Get("/swagger/*", swagger.HandlerDefault)

	log.Fatal(app.Listen(":8080"))
}

func syncData() {
	log.Println("🔄 Sincronizando datos de stocks con la API...")
	stockinterfaces.RunStockSync()
	log.Println("🔄 Sincronización de stocks completada.")
}

func exportData(path string, table string) {
	log.Println("Iniciando Exportación de datos...")
	dataBase := db.NewCockroachDB()
	defer dataBase.Close()
	var err error

	switch table {
	case "stocks":
		err = stocks.ExportStocksToJSON(dataBase, path)
	case "finances":
		err = finances.ExportFinanceDataToJSON(dataBase, path)
	default:
		err = errors.New("Nombre de la tabla no existe")
	}

	if err != nil {
		log.Fatalf("Error exportando datos: %v", err)
	}

	log.Printf("Datos exportados a %s", path)
}

func importData(path string, table string) {
	log.Println("Iniciando la importación de datos...")
	dataBase := db.NewCockroachDB()
	defer dataBase.Close()

	var err error

	switch table {
	case "stocks":
		err = stocks.ImportStocksFromJSON(dataBase, path)
	case "finances":
		err = finances.ImportFinanceDataFromJSON(dataBase, path)
	default:
		err = errors.New("Nombre de la tabla no existe")
	}

	if err != nil {
		log.Fatalf("Error exportando datos: %v", err)
	}

	log.Printf("Datos importados con éxito!")

}

func updateFinance() {
	financeinterfaces.SyncFinanceHandler()
}
