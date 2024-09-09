package infra

import (
	"fmt"
	"sync"

	"github.com/meilisearch/meilisearch-go"
)

var (
	meilisearchOnce   sync.Once
	meilisearchClient meilisearch.ServiceManager
)

func LoadMeilisearch() meilisearch.ServiceManager {
	meilisearchOnce.Do(func() {
		cfg := LoadConfig()

		client := meilisearch.New(
			fmt.Sprintf("http://%s:%d", cfg.Meilisearch.Host, cfg.Meilisearch.Port),
			meilisearch.WithAPIKey(cfg.Meilisearch.Key),
		)

		meilisearchClient = client
	})

	return meilisearchClient
}
