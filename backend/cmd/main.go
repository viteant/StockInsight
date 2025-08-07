package main

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"github.com/viteant/stockinsight/internal/db"
	"github.com/viteant/stockinsight/internal/db/seeds/finances"
	"github.com/viteant/stockinsight/internal/db/seeds/stocks"
	"github.com/viteant/stockinsight/internal/finance/infrastructure/repository"
	"github.com/viteant/stockinsight/internal/finance/infrastructure/scraper"
	usecases "github.com/viteant/stockinsight/internal/finance/use-cases"
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
	log.Println("🚀 Iniciando servidor...")
	// Aquí levantarás tu API más adelante
}

func syncData() {
	log.Println("🔄 Sincronizando datos de stocks con la API...")
	interfaces.RunStockSync()
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
	dataBase := db.NewCockroachDB()

	useCase := usecases.NewUpdateFinanceDataUseCase(
		repository.NewCockroachStockRepository(dataBase),
		repository.NewCockroachFinanceRepository(dataBase),
		scraper.NewYahooFinanceScraper(),
	)

	if err := useCase.Execute(); err != nil {
		log.Fatalf("Error ejecutando UpdateFinanceDataUseCase: %v", err)
	}
}
