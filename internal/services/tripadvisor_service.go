package services

import (
	"context"

	sqlc_db "apify-poi-data/db/sqlc"
	"apify-poi-data/pkg/apify"
	tripsadvisor_v1 "apify-poi-data/proto/apify/tripsadvisor/v1"
)

type TripadvisorService struct {
	tripsadvisor_v1.UnimplementedTripadvisorServiceServer
	Database    *sqlc_db.Database
	ApifyClient *apify.Client
}

func (t *TripadvisorService) SearchTripadvisor(ctx context.Context, in *tripsadvisor_v1.SearchRequest) (*tripsadvisor_v1.SearchResponse, error) {
	return &tripsadvisor_v1.SearchResponse{
		Status: "OK",
	}, nil
}
