package interfaces

import (
	"log"

	"github.com/viteant/stockinsight/internal/db"
	"github.com/viteant/stockinsight/internal/finance/infrastructure/repository"
	"github.com/viteant/stockinsight/internal/finance/infrastructure/scraper"
	usecases "github.com/viteant/stockinsight/internal/finance/use-cases"
)

func SyncFinanceHandler() {
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
