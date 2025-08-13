package inferno

import "github.com/vayzur/spark/pkg/httputil"

type InfernoClient struct {
	httpClient *httputil.Client
	address    string
	token      string
	nodeID     string
}

func NewInfernoClient(httpClient *httputil.Client, address, token, nodeID string) *InfernoClient {
	return &InfernoClient{
		httpClient: httpClient,
		address:    address,
		token:      token,
		nodeID:     nodeID,
	}
}
