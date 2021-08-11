package hackernews

import (
	"context"
)

type handler struct {
	Handler
	request RequestType
}

func NewRequestHandler(request RequestType) Fetcher {
	fbClient := NewFirebaseClientWithDefaultConfig(context.Background())
	return &handler{
		fbClient,
		request,
	}
}

func (h *handler) Fetch() ([]uint, error) {
	return h.FetchRequestedIDs(h.request)
}
