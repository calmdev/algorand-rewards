package nodely

// indexerBaseURL is the base URL for the indexer API.
const indexerBaseURL = "https://mainnet-idx.4160.nodely.dev"

// ClientIndexer represents an HTTP client for the indexer API.
type ClientIndexer struct {
	Client
}

// NewClientIndexer creates a new ClientIndexer instance for the indexer API.
//
// Docs: https://nodely.io/docs/free/endpoints/#free-archival-indexer-api
func NewClientIndexer() *ClientIndexer {
	return &ClientIndexer{Client{BaseURL: indexerBaseURL}}
}
