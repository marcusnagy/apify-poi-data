package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/uber/h3-go/v4"
	"google.golang.org/protobuf/types/known/structpb"

	sqlc_db "apify-poi-data/db/sqlc"
	poi_v1 "apify-poi-data/proto/apify/poi/v1"
)

const (
	DATABASE_RESOLUTION = 9
)

type PoiService struct {
	Database *sqlc_db.Database
	poi_v1.UnimplementedPoiServiceServer
}

func rawMessageToValue(raw json.RawMessage) (*structpb.Struct, error) {
	var m any
	if json.Valid(raw) {
		if err := json.Unmarshal(raw, &m); err != nil {
			return nil, err
		}
	} else {
		m = string(raw)
	}
	return structpb.NewStruct(map[string]any{"data": m})
}

func (p PoiService) toPOIs(res []sqlc_db.PoiDataSchemaGoogleMap) ([]*poi_v1.Poi, error) {
	var pois []*poi_v1.Poi
	for _, poi := range res {
		poi, err := p.toPOI(poi)
		if err != nil {
			return nil, err
		}
		pois = append(pois, poi)
	}
	return pois, nil
}

func (p PoiService) toPOI(poi sqlc_db.PoiDataSchemaGoogleMap) (*poi_v1.Poi, error) {
	var err error
	var hAds,
		peopleAlsoSearch,
		placesTags,
		reviewTags,
		gasPrices,
		similarHotelsNearby,
		hotelReviewSummary,
		popularTimesHistogram,
		questionsAndAnswers,
		webResults,
		tableReservationLinks,
		bookingLinks,
		reviews,
		restaurantData,
		ownerUpdates,
		orderBy,
		userPlaceNote,
		updatesFromCustomers,
		additionalInfo *structpb.Struct

	var openingHours []*poi_v1.OpeningHour

	if poi.HotelAds != nil {
		hAds, err = rawMessageToValue(poi.HotelAds)
		if err != nil {
			return nil, err
		}
	}
	if poi.AdditionalInfo != nil {
		additionalInfo, err = rawMessageToValue(poi.AdditionalInfo)
		if err != nil {
			return nil, err
		}
	}
	if poi.OpeningHours != nil {
		var o []*poi_v1.OpeningHour
		err := json.Unmarshal(poi.OpeningHours, &o)
		if err != nil {
			return nil, err
		}
		openingHours = append(openingHours, o...)
	}
	if poi.PeopleAlsoSearch != nil {
		peopleAlsoSearch, err = rawMessageToValue(poi.PeopleAlsoSearch)
		if err != nil {
			return nil, err
		}
	}
	if poi.GasPrices != nil {
		gasPrices, err = rawMessageToValue(poi.GasPrices)
		if err != nil {
			return nil, err
		}
	}
	if poi.PlacesTags != nil {
		placesTags, err = rawMessageToValue(poi.PlacesTags)
		if err != nil {
			return nil, err
		}
	}
	if poi.ReviewsTags != nil {
		reviewTags, err = rawMessageToValue(poi.ReviewsTags)
		if err != nil {
			return nil, err
		}
	}
	if poi.SimilarHotelsNearby != nil {
		similarHotelsNearby, err = rawMessageToValue(poi.SimilarHotelsNearby)
		if err != nil {
			return nil, err
		}
	}
	if poi.HotelReviewSummary != nil {
		hotelReviewSummary, err = rawMessageToValue(poi.HotelReviewSummary)
		if err != nil {
			return nil, err
		}
	}
	if poi.PopularTimesHistogram != nil {
		popularTimesHistogram, err = rawMessageToValue(poi.PopularTimesHistogram)
		if err != nil {
			return nil, err
		}
	}
	if poi.QuestionsAndAnswers != nil {
		questionsAndAnswers, err = rawMessageToValue(poi.QuestionsAndAnswers)
		if err != nil {
			return nil, err
		}
	}
	if poi.WebResults != nil {
		webResults, err = rawMessageToValue(poi.WebResults)
		if err != nil {
			return nil, err
		}
	}
	if poi.TableReservationLinks != nil {
		tableReservationLinks, err = rawMessageToValue(poi.TableReservationLinks)
		if err != nil {
			return nil, err
		}
	}
	if poi.BookingLinks != nil {
		bookingLinks, err = rawMessageToValue(poi.BookingLinks)
		if err != nil {
			return nil, err
		}
	}
	if poi.Reviews != nil {
		reviews, err = rawMessageToValue(poi.Reviews)
		if err != nil {
			return nil, err
		}
	}
	if poi.RestaurantData != nil {
		restaurantData, err = rawMessageToValue(poi.RestaurantData)
		if err != nil {
			return nil, err
		}
	}
	if poi.OwnerUpdates != nil {
		ownerUpdates, err = rawMessageToValue(poi.OwnerUpdates)
		if err != nil {
			return nil, err
		}
	}
	if poi.OrderBy != nil {
		orderBy, err = rawMessageToValue(poi.OrderBy)
		if err != nil {
			return nil, err
		}
	}
	if poi.UserPlaceNote != nil {
		userPlaceNote, err = rawMessageToValue(poi.UserPlaceNote)
		if err != nil {
			return nil, err
		}
	}
	if poi.UpdatesFromCustomers != nil {
		updatesFromCustomers, err = rawMessageToValue(poi.UpdatesFromCustomers)
		if err != nil {
			return nil, err
		}
	}

	return &poi_v1.Poi{
		Id:                      poi.ID,
		SearchString:            poi.SearchString.String,
		Rank:                    poi.Rank.Int32,
		SearchPageUrl:           poi.SearchPageUrl.String,
		IsAdvertisement:         poi.IsAdvertisement.Bool,
		Title:                   poi.Title.String,
		SubTitle:                poi.SubTitle.String,
		Price:                   poi.Price.String,
		CategoryName:            poi.CategoryName.String,
		Address:                 poi.Address.String,
		Neighborhood:            poi.Neighborhood.String,
		Street:                  poi.Street.String,
		City:                    poi.City.String,
		PostalCode:              poi.PostalCode.String,
		State:                   poi.State.String,
		CountryCode:             poi.CountryCode.String,
		Website:                 poi.Website.String,
		Phone:                   poi.Phone.String,
		PhoneUnformatted:        poi.PhoneUnformatted.String,
		ClaimThisBusiness:       poi.ClaimThisBusiness.Bool,
		LocationLat:             poi.LocationLat.Float64,
		LocationLng:             poi.LocationLng.Float64,
		TotalScore:              poi.TotalScore.Float64,
		PermanentlyClosed:       poi.PermanentlyClosed.Bool,
		TemporarilyClosed:       poi.TemporarilyClosed.Bool,
		PlaceId:                 poi.PlaceID.String,
		Categories:              poi.Categories,
		Fid:                     poi.Fid.String,
		Cid:                     poi.Cid.String,
		ReviewsCount:            poi.ReviewsCount.Int32,
		ImagesCount:             poi.ImagesCount.Int32,
		ImageCategories:         poi.ImageCategories,
		ScrapedAt:               poi.ScrapedAt.Time.String(),
		GoogleFoodUrl:           poi.GoogleFoodUrl.String,
		HotelAds:                hAds,
		OpeningHours:            openingHours,
		PeopleAlsoSearch:        peopleAlsoSearch,
		PlacesTags:              placesTags,
		ReviewsTags:             reviewTags,
		AdditionalInfo:          additionalInfo,
		GasPrices:               gasPrices,
		Url:                     poi.Url.String,
		ImageUrl:                poi.ImageUrl.String,
		Kgmid:                   poi.Kgmid.String,
		Geom:                    poi.Geom,
		H3Index:                 poi.H3Index.String,
		SearchPageLoadedUrl:     poi.SearchPageLoadedUrl.String,
		Description:             poi.Description.String,
		LocatedIn:               poi.LocatedIn.String,
		PlusCode:                poi.PlusCode.String,
		Menu:                    poi.Menu.String,
		ReserveTableUrl:         poi.ReserveTableUrl.String,
		HotelStars:              poi.HotelStars.String,
		HotelDescription:        poi.HotelDescription.String,
		CheckInDate:             poi.CheckInDate.String,
		CheckOutDate:            poi.CheckOutDate.String,
		SimilarHotelsNearby:     similarHotelsNearby,
		HotelReviewSummary:      hotelReviewSummary,
		PopularTimesLiveText:    poi.PopularTimesLiveText.String,
		PopularTimesLivePercent: poi.PopularTimesLivePercent.Int32,
		PopularTimesHistogram:   popularTimesHistogram,
		QuestionsAndAnswers:     questionsAndAnswers,
		UpdatesFromCustomers:    updatesFromCustomers,
		WebResults:              webResults,
		ParentPlaceUrl:          poi.ParentPlaceUrl.String,
		TableReservationLinks:   tableReservationLinks,
		BookingLinks:            bookingLinks,
		OrderBy:                 orderBy,
		Images:                  poi.Images.String,
		UserPlaceNote:           userPlaceNote,
		ImageUrls:               poi.ImageUrls,
		Reviews:                 reviews,
		RestaurantData:          restaurantData,
		OwnerUpdates:            ownerUpdates,
	}, nil
}

func (p *PoiService) ListPOIByH3Cells(in *poi_v1.ListPOIsByH3CellsRequest, stream poi_v1.PoiService_ListPOIByH3CellsServer) error {
	indexes := in.GetParentCells()

	// Validate h3 indexes
	for _, index := range indexes {
		c := h3.IndexFromString(index)
		cell := h3.Cell(c)
		if !cell.IsValid() {
			return fmt.Errorf("invalid h3 index; %s", index)

		}
		if cell.Resolution() >= DATABASE_RESOLUTION {
			return fmt.Errorf("h3 index resolution is too small(higher than 9 is not allowed); %s", index)
		}
	}

	res, err := p.Database.Queries.ListPOIsByH3Cells(stream.Context(), sqlc_db.ListPOIsByH3CellsParams{
		Column1: DATABASE_RESOLUTION,
		Column2: indexes,
	})
	if err != nil {
		return err
	}
	const batchSize = 150 // Adjust the batch size as needed
	var pois []*poi_v1.Poi
	for i, poi := range res {
		poi, err := p.toPOI(poi)
		if err != nil {
			fmt.Printf("error converting poi: %v", err)
			continue
		}
		pois = append(pois, poi)
		if (i+1)%batchSize == 0 {
			if err := stream.Send(&poi_v1.ListPOIResponse{Pois: pois}); err != nil {
				return err
			}
			pois = nil // Reset the batch
		}
	}
	// Send remaining POIs
	if len(pois) > 0 {
		if err := stream.Send(&poi_v1.ListPOIResponse{Pois: pois}); err != nil {
			return err
		}
	}
	return nil
}

func (p *PoiService) ListPOIInBox(ctx context.Context, in *poi_v1.ListPOIInBoxRequest) (*poi_v1.ListPOIResponse, error) {
	res, err := p.Database.Queries.ListPOIInBox(ctx, sqlc_db.ListPOIInBoxParams{
		Column1: in.GetMinX(),
		Column2: in.GetMinY(),
		Column3: in.GetMaxX(),
		Column4: in.GetMaxY(),
	})
	if err != nil {
		return nil, err
	}

	pois, err := p.toPOIs(res)
	if err != nil {
		return nil, err
	}

	return &poi_v1.ListPOIResponse{
		Pois: pois,
	}, nil
}

func (p *PoiService) ListPOIInBoxWithCategorySearch(ctx context.Context, in *poi_v1.ListPOIInBoxWithCategorySearchRequest) (*poi_v1.ListPOIResponse, error) {
	res, err := p.Database.Queries.ListPOIInBoxWithCategoryH3(
		ctx,
		sqlc_db.ListPOIInBoxWithCategoryH3Params{
			Column1: in.GetMinX(),
			Column2: in.GetMinY(),
			Column3: in.GetMaxX(),
			Column4: in.GetMaxY(),
			Column5: pgtype.Text{
				String: in.GetCategorySubstring(),
				Valid:  true,
			},
		})
	if err != nil {
		return nil, err
	}

	pois, err := p.toPOIs(res)
	if err != nil {
		return nil, err
	}

	return &poi_v1.ListPOIResponse{
		Pois: pois,
	}, nil
}

func (p *PoiService) ListPOIAlongRoute(ctx context.Context, in *poi_v1.ListPOIAlongRouteRequest) (*poi_v1.ListPOIResponse, error) {
	res, err := p.Database.Queries.ListPOIAlongRouteH3(context.Background(), sqlc_db.ListPOIAlongRouteH3Params{
		Column1: in.GetALat(),
		Column2: in.GetALon(),
		Column3: in.GetBLat(),
		Column4: in.GetBLon(),
		Column5: float64(in.GetBuffer()),
	})
	if err != nil {
		return nil, err
	}

	pois, err := p.toPOIs(res)
	if err != nil {
		return nil, err
	}
	return &poi_v1.ListPOIResponse{
		Pois: pois,
	}, nil
}

func (p *PoiService) ListPOIAlongRouteWithCategorySearch(ctx context.Context, in *poi_v1.ListPOIAlongRouteWithCategoryRequest) (*poi_v1.ListPOIResponse, error) {
	res, err := p.Database.Queries.ListPOIAlongRouteWithCategoryH3(context.Background(), sqlc_db.ListPOIAlongRouteWithCategoryH3Params{
		Column1: in.GetALat(),
		Column2: in.GetALon(),
		Column3: in.GetBLat(),
		Column4: in.GetBLon(),
		Column5: float64(in.GetBuffer()),
		Column6: pgtype.Text{
			String: in.GetCategorySubstring(),
			Valid:  true,
		},
	})
	if err != nil {
		return nil, err
	}

	pois, err := p.toPOIs(res)
	if err != nil {
		return nil, err
	}

	return &poi_v1.ListPOIResponse{
		Pois: pois,
	}, nil
}
