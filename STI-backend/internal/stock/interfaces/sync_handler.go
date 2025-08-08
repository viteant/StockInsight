package interfaces

import (
	"github.com/viteant/stockinsight/internal/db"
	"github.com/viteant/stockinsight/internal/stock/infrastructure/api"
	"github.com/viteant/stockinsight/internal/stock/infrastructure/repository"
	"github.com/viteant/stockinsight/internal/stock/use_cases"
)

func RunStockSync() {
	dbConn := db.NewCockroachDB()
	defer dbConn.Close()

	fetcher := api.NewExternalAPIClient()
	repo := repository.NewCockroachStockRepository(dbConn)

	sync := use_cases.NewSyncService(fetcher, repo)
	if err := sync.Sync(); err != nil {
		panic(err)
	}
}
