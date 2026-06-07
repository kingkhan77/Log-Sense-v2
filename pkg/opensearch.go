package pkg

import (
	"fmt"

	opensearch "github.com/opensearch-project/opensearch-go"
)

func NewOpenSearch(
	cfg *Config,
) (*opensearch.Client, error) {

	client, err := opensearch.NewClient(
		opensearch.Config{
			Addresses: []string{
				fmt.Sprintf(
					"http://%s:%d",
					cfg.OpenSearch.Host,
					cfg.OpenSearch.Port,
				),
			},
			Username: cfg.OpenSearch.Username,
			Password: cfg.OpenSearch.Password,
		},
	)

	return client, err
}