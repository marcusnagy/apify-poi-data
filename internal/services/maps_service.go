package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/uber/h3-go/v4"

	sqlc_db "apify-poi-data/db/sqlc"
	"apify-poi-data/internal/models"
	"apify-poi-data/internal/services/converter"
	"apify-poi-data/pkg/apify"
	maps_v1 "apify-poi-data/proto/apify/maps/v1"
)

type MapsService struct {
	maps_v1.UnimplementedMapsServiceServer
	ApifyClient *apify.Client
	Database    *sqlc_db.Database
}

func (m *MapsService) castPOIToPlace(pois []models.POI) []models.Place {
	places := make([]models.Place, 0, len(pois))
	for _, poi := range pois {
		switch p := poi.(type) {
		case *models.Place:
			places = append(places, *p)
		default:
			fmt.Printf("Skipping non-Place POI: %s\n", poi.GetType())
		}
	}
	return places
}

func (m *MapsService) castPOIToPlaceScraper(pois []models.POI) []models.PlaceScraper {
	places := make([]models.PlaceScraper, 0, len(pois))
	for _, poi := range pois {
		switch p := poi.(type) {
		case *models.PlaceScraper:
			places = append(places, *p)
		default:
			fmt.Printf("Skipping non-Place POI: %s, Type: %s\n", poi.GetType(), p)
		}
	}
	return places
}

func (m *MapsService) processPOIGooglePlace(ctx context.Context, poi models.Place) error {
	poiParams := sqlc_db.InsertPOIParams{
		SearchString: pgtype.Text{
			String: poi.SearchString,
			Valid:  true,
		},
		Rank: pgtype.Int4{
			Int32: int32(poi.Rank),
			Valid: true,
		},
		SearchPageUrl: pgtype.Text{
			String: poi.SearchPageURL,
			Valid:  true,
		},
		IsAdvertisement: pgtype.Bool{
			Bool:  poi.IsAdvertisement,
			Valid: true,
		},
		Title: pgtype.Text{
			String: poi.Title,
			Valid:  true,
		},
		SubTitle: pgtype.Text{
			String: poi.SubTitle,
			Valid:  true,
		},
		CategoryName: pgtype.Text{
			String: poi.CategoryName,
			Valid:  true,
		},
		Address: pgtype.Text{
			String: poi.Address,
			Valid:  true,
		},
		Street: pgtype.Text{
			String: poi.Street,
			Valid:  true,
		},
		City: pgtype.Text{
			String: poi.City,
			Valid:  true,
		},
		PostalCode: pgtype.Text{
			String: poi.PostalCode,
			Valid:  true,
		},
		CountryCode: pgtype.Text{
			String: poi.CountryCode,
			Valid:  true,
		},
		ClaimThisBusiness: pgtype.Bool{
			Bool:  poi.ClaimThisBusiness,
			Valid: true,
		},
		LocationLat: pgtype.Float8{
			Float64: poi.Location.Lat,
			Valid:   true,
		},
		LocationLng: pgtype.Float8{
			Float64: poi.Location.Lng,
			Valid:   true,
		},
		TotalScore: pgtype.Float8{
			Float64: poi.TotalScore,
			Valid:   true,
		},
		PermanentlyClosed: pgtype.Bool{
			Bool:  poi.PermanentlyClosed,
			Valid: true,
		},
		TemporarilyClosed: pgtype.Bool{
			Bool:  poi.TemporarilyClosed,
			Valid: true,
		},
		PlaceID: pgtype.Text{
			String: poi.PlaceID,
			Valid:  true,
		},
		Categories: poi.Categories,
		Fid: pgtype.Text{
			String: poi.FID,
			Valid:  true,
		},
		Cid: pgtype.Text{
			String: poi.CID,
			Valid:  true,
		},
		ReviewsCount: pgtype.Int4{
			Int32: int32(poi.ReviewsCount),
			Valid: true,
		},
		ImagesCount: pgtype.Int4{
			Int32: int32(poi.ImagesCount),
			Valid: true,
		},
		ImageCategories:  poi.ImageCategories,
		HotelAds:         poi.HotelAds,
		OpeningHours:     poi.OpeningHours,
		PeopleAlsoSearch: poi.PeopleAlsoSearch,
		PlacesTags:       poi.PlacesTags,
		ReviewsTags:      poi.ReviewsTags,
		AdditionalInfo:   poi.AdditionalInfo,
		GasPrices:        poi.GasPrices,
		Url: pgtype.Text{
			String: poi.URL,
			Valid:  true,
		},
		ImageUrl: pgtype.Text{
			String: poi.ImageURL,
			Valid:  true,
		},
		Kgmid: pgtype.Text{
			String: poi.Kgmid,
			Valid:  true,
		},
		Neighborhood: pgtype.Text{
			String: models.ValueOrEmpty(poi.Neighborhood),
			Valid:  true,
		},
		Price: pgtype.Text{
			String: models.ValueOrEmpty(poi.Price),
			Valid:  true,
		},
		State: pgtype.Text{
			String: models.ValueOrEmpty(poi.State),
			Valid:  true,
		},
		Phone: pgtype.Text{
			String: models.ValueOrEmpty(poi.Phone),
			Valid:  true,
		},
		Website: pgtype.Text{
			String: models.ValueOrEmpty(poi.Website),
			Valid:  true,
		},
		PhoneUnformatted: pgtype.Text{
			String: models.ValueOrEmpty(poi.PhoneUnformatted),
			Valid:  true,
		},
		GoogleFoodUrl: pgtype.Text{
			String: models.ValueOrEmpty(poi.GoogleFoodURL),
			Valid:  false,
		},
	}

	if cell, err := h3.LatLngToCell(h3.LatLng{Lat: poi.Location.Lat, Lng: poi.Location.Lng}, DATABASE_RESOLUTION); err != nil {
		log.Printf("Failed to get H3Index: %v", err)
	} else {
		poiParams.H3Index = pgtype.Text{
			String: cell.String(),
			Valid:  true,
		}
	}

	// Convert to proto
	t, err := time.Parse(time.RFC3339, poi.ScrapedAt)
	if err != nil {
		log.Printf("Failed to parse time: %v", err)
		t = time.Now()
	}
	poiParams.ScrapedAt = pgtype.Timestamptz{
		Time:  t,
		Valid: true,
	}
	// Send to client
	err = m.insertPOItoDB(ctx, poiParams)
	if err != nil {
		return err
	}
	return nil
}

func (m *MapsService) processPOIGooglePlaceScraper(ctx context.Context, poi models.PlaceScraper) error {
	poiParams := sqlc_db.InsertPOIParams{
		SearchString: pgtype.Text{
			String: poi.SearchString,
			Valid:  true,
		},
		Rank: pgtype.Int4{
			Int32: int32(poi.Rank),
			Valid: true,
		},
		SearchPageUrl: pgtype.Text{
			String: poi.SearchPageURL,
			Valid:  true,
		},
		IsAdvertisement: pgtype.Bool{
			Bool:  poi.IsAdvertisement,
			Valid: true,
		},
		Title: pgtype.Text{
			String: poi.Title,
			Valid:  true,
		},
		SubTitle: pgtype.Text{
			String: models.ValueOrEmpty(poi.SubTitle),
			Valid:  true,
		},
		CategoryName: pgtype.Text{
			String: poi.CategoryName,
			Valid:  true,
		},
		Address: pgtype.Text{
			String: poi.Address,
			Valid:  true,
		},
		Street: pgtype.Text{
			String: poi.Street,
			Valid:  true,
		},
		City: pgtype.Text{
			String: poi.City,
			Valid:  true,
		},
		PostalCode: pgtype.Text{
			String: poi.PostalCode,
			Valid:  true,
		},
		CountryCode: pgtype.Text{
			String: poi.CountryCode,
			Valid:  true,
		},
		ClaimThisBusiness: pgtype.Bool{
			Bool:  poi.ClaimThisBusiness,
			Valid: true,
		},
		LocationLat: pgtype.Float8{
			Float64: poi.Location.Lat,
			Valid:   true,
		},
		LocationLng: pgtype.Float8{
			Float64: poi.Location.Lng,
			Valid:   true,
		},
		TotalScore: pgtype.Float8{
			Float64: poi.TotalScore,
			Valid:   true,
		},
		PermanentlyClosed: pgtype.Bool{
			Bool:  poi.PermanentlyClosed,
			Valid: true,
		},
		TemporarilyClosed: pgtype.Bool{
			Bool:  poi.TemporarilyClosed,
			Valid: true,
		},
		PlaceID: pgtype.Text{
			String: poi.PlaceID,
			Valid:  true,
		},
		Categories: poi.Categories,
		Fid: pgtype.Text{
			String: poi.FID,
			Valid:  true,
		},
		Cid: pgtype.Text{
			String: poi.CID,
			Valid:  true,
		},
		ReviewsCount: pgtype.Int4{
			Int32: int32(poi.ReviewsCount),
			Valid: true,
		},
		ImagesCount: pgtype.Int4{
			Int32: int32(poi.ImagesCount),
			Valid: true,
		},
		ImageCategories:  poi.ImageCategories,
		HotelAds:         poi.HotelAds,
		OpeningHours:     poi.OpeningHours,
		PeopleAlsoSearch: poi.PeopleAlsoSearch,
		PlacesTags:       poi.PlacesTags,
		ReviewsTags:      poi.ReviewsTags,
		AdditionalInfo:   poi.AdditionalInfo,
		GasPrices:        poi.GasPrices,
		Url: pgtype.Text{
			String: poi.URL,
			Valid:  true,
		},
		ImageUrl: pgtype.Text{
			String: poi.ImageURL,
			Valid:  true,
		},
		Kgmid: pgtype.Text{
			String: poi.Kgmid,
			Valid:  true,
		},
		Neighborhood: pgtype.Text{
			String: models.ValueOrEmpty(poi.Neighborhood),
			Valid:  true,
		},
		Price: pgtype.Text{
			String: models.ValueOrEmpty(poi.Price),
			Valid:  true,
		},
		State: pgtype.Text{
			String: models.ValueOrEmpty(poi.State),
			Valid:  true,
		},
		Phone: pgtype.Text{
			String: models.ValueOrEmpty(poi.Phone),
			Valid:  true,
		},
		Website: pgtype.Text{
			String: models.ValueOrEmpty(poi.Website),
			Valid:  true,
		},
		PhoneUnformatted: pgtype.Text{
			String: models.ValueOrEmpty(poi.PhoneUnformatted),
			Valid:  true,
		},
		GoogleFoodUrl: pgtype.Text{
			String: models.ValueOrEmpty(poi.GoogleFoodURL),
			Valid:  false,
		},
		SearchPageLoadedUrl: pgtype.Text{
			String: models.ValueOrEmpty(poi.SearchPageLoadedURL),
			Valid:  true,
		},
		Description: pgtype.Text{
			String: models.ValueOrEmpty(poi.Description),
			Valid:  true,
		},
		LocatedIn: pgtype.Text{
			String: models.ValueOrEmpty(poi.LocatedIn),
			Valid:  true,
		},
		PlusCode: pgtype.Text{
			String: models.ValueOrEmpty(poi.PlusCode),
			Valid:  true,
		},
		Menu: pgtype.Text{
			String: models.ValueOrEmpty(poi.Menu),
			Valid:  true,
		},
		ReserveTableUrl: pgtype.Text{
			String: models.ValueOrEmpty(poi.ReserveTableURL),
			Valid:  true,
		},
		HotelStars: pgtype.Text{
			String: models.ValueOrEmpty(poi.HotelStars),
			Valid:  true,
		},
		HotelDescription: pgtype.Text{
			String: models.ValueOrEmpty(poi.HotelDescription),
			Valid:  true,
		},
		CheckInDate: pgtype.Text{
			String: models.ValueOrEmpty(poi.CheckInDate),
			Valid:  true,
		},
		CheckOutDate: pgtype.Text{
			String: models.ValueOrEmpty(poi.CheckOutDate),
			Valid:  true,
		},
		SimilarHotelsNearby: poi.SimilarHotelsNearby,
		HotelReviewSummary:  poi.HotelReviewSummary,
		PopularTimesLiveText: pgtype.Text{
			String: models.ValueOrEmpty(poi.PopularTimesLiveText),
			Valid:  true,
		},
		PopularTimesLivePercent: pgtype.Int4{
			Int32: int32(models.ValueOrZero(poi.PopularTimesLivePercent)),
			Valid: true,
		},
		QuestionsAndAnswers:  poi.QuestionsAndAnswers,
		UpdatesFromCustomers: poi.UpdatesFromCustomers,
		WebResults:           poi.WebResults,
		ParentPlaceUrl: pgtype.Text{
			String: models.ValueOrEmpty(poi.ParentPlaceURL),
			Valid:  true,
		},
		TableReservationLinks: poi.TableReservationLinks,
		BookingLinks:          poi.BookingLinks,
		OrderBy:               poi.OrderBy,
		Images: pgtype.Text{
			String: models.ValueOrEmpty(poi.Images),
			Valid:  true,
		},
		ImageUrls:      poi.ImageUrls,
		Reviews:        poi.Reviews,
		UserPlaceNote:  poi.UserPlaceNote,
		RestaurantData: poi.RestaurantData,
		OwnerUpdates:   poi.OwnerUpdates,
	}

	if cell, err := h3.LatLngToCell(h3.LatLng{Lat: poi.Location.Lat, Lng: poi.Location.Lng}, DATABASE_RESOLUTION); err != nil {
		log.Printf("Failed to get H3Index: %v", err)
	} else {
		poiParams.H3Index = pgtype.Text{
			String: cell.String(),
			Valid:  true,
		}
	}

	if poi.PopularTimesHistogram != nil {
		popularTimesHistogram, err := json.Marshal(poi.PopularTimesHistogram)
		if err != nil {
			log.Printf("Failed to marshal PopularTimesHistogram: %v", err)
		}
		poiParams.PopularTimesHistogram = popularTimesHistogram
	}

	// Convert to proto
	t, err := time.Parse(time.RFC3339, poi.ScrapedAt)
	if err != nil {
		log.Printf("Failed to parse time: %v", err)
		t = time.Now()
	}
	poiParams.ScrapedAt = pgtype.Timestamptz{
		Time:  t,
		Valid: true,
	}
	// Send to client
	err = m.insertPOItoDB(ctx, poiParams)
	if err != nil {
		return err
	}
	return nil
}

func (m *MapsService) insertPOItoDB(ctx context.Context, params sqlc_db.InsertPOIParams) error {
	_, err := m.Database.Queries.InsertPOI(ctx, params)
	if err != nil {
		// Check if the error is a unique constraint violation
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			log.Printf("POI already exists: %v", err)
			return nil // Skip to the next POI
		}
		// Check if the error is sql.ErrNoRows
		if err == pgx.ErrNoRows {
			log.Printf("No rows in result set: %v", err)
			return nil // Skip to the next POI
		}
		log.Printf("Failed to insert POI: %v", err)
		return err
	}
	return nil
}

func (m *MapsService) handleGoogleMapsScraper(ctx context.Context, pois []models.PlaceScraper) {
	for _, poi := range pois {
		if err := m.processPOIGooglePlaceScraper(ctx, poi); err != nil {
			log.Printf("Failed to process POI: %v", err)

		}
	}
}

func (m *MapsService) handleGoogleMapsExtractor(ctx context.Context, pois []models.Place) {
	for _, poi := range pois {
		if err := m.processPOIGooglePlace(ctx, poi); err != nil {
			log.Printf("Failed to process POI: %v", err)
		}
	}
}

func (m *MapsService) InsertApifyDatasetItems(ctx context.Context, in *maps_v1.DatasetItemsRequest) (*maps_v1.DatasetItemsResponse, error) {
	data, err := m.ApifyClient.GetDataset(in.GetDatasetId())
	if err != nil {
		return nil, err
	}
	pois, err := models.ParsePOIsFromJSON(data)
	if err != nil {
		return nil, err
	}
	switch in.GetDatasetType() {
	case maps_v1.DatasetItemsRequest_GOOGLE_MAPS_EXTRACTOR:
		fmt.Println("Data received inside SearchGoogleMaps")
		d := m.castPOIToPlace(pois)
		m.handleGoogleMapsExtractor(ctx, d)
	case maps_v1.DatasetItemsRequest_GOOGLE_MAPS_SCRAPER:
		fmt.Println("Data received inside SearchGoogleMapsScraper")
		d := m.castPOIToPlaceScraper(pois)
		m.handleGoogleMapsScraper(ctx, d)

	}
	return &maps_v1.DatasetItemsResponse{
		Status: "success",
	}, nil
}

func (m *MapsService) SearchGoogleMapsExtractor(ctx context.Context, in *maps_v1.SearchRequest) (*maps_v1.SearchResponse, error) {

	req, err := converter.SearchRequestToInputPayloadMaps(in)
	if err != nil {
		return nil, err
	}

	resp := m.ApifyClient.ExtractPOIs(req, int(in.GetNumberOfResults()), true)

	select {
	case data := <-resp.Data:
		fmt.Println("Data received inside SearchGoogleMaps")
		d := m.castPOIToPlace(data)
		m.handleGoogleMapsExtractor(ctx, d)
		return &maps_v1.SearchResponse{
			Status: "success",
		}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-resp.Err:
		log.Printf("Error: %v", err)
		return nil, err
	}
}

func (m *MapsService) SearchGoogleMapsScraper(ctx context.Context, request *maps_v1.ScraperRequest) (*maps_v1.SearchResponse, error) {
	req, err := converter.SearchRequestToInputPayloadMapsScraper(request)
	if err != nil {
		return nil, err
	}

	resp := m.ApifyClient.ScrapePOIs(req, true)

	select {
	case data := <-resp.Data:
		fmt.Println("Data received inside SearchGoogleMapsScraper")
		d := m.castPOIToPlaceScraper(data)
		m.handleGoogleMapsScraper(ctx, d)
		return &maps_v1.SearchResponse{
			Status: "success",
		}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-resp.Err:
		log.Printf("Error: %v", err)
	}

	return nil, nil
}
