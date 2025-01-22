package converter

import (
	"encoding/json"
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"apify-poi-data/internal/models"
	maps_v1 "apify-poi-data/proto/apify/maps/v1"
)

func customGeolocationToJson(geolocationRequest *maps_v1.CustomGeolocation) (json.RawMessage, error) {
	type GeoLocation struct {
		Type        string           `json:"type"`
		Coordinates [][][][2]float64 `json:"coordinates"`
	}

	var coordinates [][][][2]float64
	switch geolocationRequest.GetType() {
	case maps_v1.CustomGeolocation_MULTIPOLYGON:
		for _, polygon := range geolocationRequest.GetPolygons() {
			var polygonCoordinates [][2]float64
			for _, coord := range polygon.GetCoordinates() {
				polygonCoordinates = append(
					polygonCoordinates,
					[2]float64{
						float64(coord.GetLongitude()),
						float64(coord.GetLatitude()),
					},
				)
			}
			coordinates = append(coordinates, [][][2]float64{polygonCoordinates})
		}
	default:
		return nil, fmt.Errorf("unknown type: %v", geolocationRequest.GetType())
	}

	// Should not insert custom geo locaiton if there are no coordinates
	if len(coordinates) == 0 {
		return nil, nil
	}

	titler := cases.Title(language.English)

	geoLocation := GeoLocation{
		Type:        titler.String(strings.ToLower(geolocationRequest.GetType().String())),
		Coordinates: coordinates,
	}

	return json.Marshal(geoLocation)
}

func SearchRequestToInputPayloadMaps(request *maps_v1.SearchRequest) (models.InputPayloadMaps, error) {
	customGeolocationJson, err := customGeolocationToJson(request.GetCustomGeolocation())
	if err != nil {
		return models.InputPayloadMaps{}, err
	}

	out := models.InputPayloadMaps{
		SearchStringsArray:        request.SearchStringsArray,
		LocationQuery:             request.GetLocationQuery(),
		MaxCrawledPlacesPerSearch: int(request.GetMaxCrawledPlacesPerSearch()),
		Language:                  request.GetLanguage(),
		CountryCode:               request.GetCountryCode(),
		City:                      request.GetCity(),
		PostalCode:                request.GetPostalCode(),
		SkipClosedPlaces:          request.GetSkipClosedPlaces(),
		PlaceMinimumStars:         request.GetPlacesMinimumStars(),
	}

	if customGeolocationJson != nil {
		out.CustomGeolocation = customGeolocationJson
	}
	return out, nil
}

func SearchRequestToInputPayloadMapsScraper(request *maps_v1.ScraperRequest) (models.ScraperInputPayloadMaps, error) {
	customGeolocationJson, err := customGeolocationToJson(request.GetCustomGeolocation())
	if err != nil {
		return models.ScraperInputPayloadMaps{}, err
	}

	out := models.ScraperInputPayloadMaps{
		SearchStringsArray:             request.SearchStringsArray,
		LocationQuery:                  request.GetLocationQuery(),
		MaxCrawledPlacesPerSearch:      int(request.GetMaxCrawledPlacesPerSearch()),
		Language:                       request.GetLanguage(),
		MaxImages:                      int(request.GetMaxImages()),
		ScrapeImageAuthors:             request.GetScrapeImageAuthors(),
		OnlyDataFromSearchPage:         request.GetOnlyDataFromSearchPage(),
		IncludeWebResults:              request.GetIncludeWebResults(),
		ScrapeDirectories:              request.GetScrapeDirectories(),
		ScrapeTableReservationProvider: request.GetScrapeTableReservationProvider(),
		MaxReviews:                     int(request.GetMaxReviews()),
		ReviewsStartDate:               request.GetReviewsStartDate(),
		ReviewsSort:                    request.GetReviewsSort(),
		ReviewsFilterString:            request.GetReviewsFilterString(),
		ReviewsOrigin:                  request.GetReviewsOrigin(),
		ScrapeReviewsPersonalData:      request.GetScrapeReviewsPersonalData(),
		MaxQuestions:                   int(request.GetMaxQuestions()),
		Zoom:                           int(request.GetZoom()),
		CountryCode:                    request.GetCountryCode(),
		City:                           request.GetCity(),
		State:                          request.GetState(),
		County:                         request.GetCounty(),
		PostalCode:                     request.GetPostalCode(),
		CategoryFilterWords:            request.GetCategoryFilterWords(),
		SearchMatching:                 request.GetSearchMatching(),
		PlaceMinimumStars:              request.GetPlaceMinimumStars(),
		SkipClosedPlaces:               request.GetSkipClosedPlaces(),
		Website:                        request.GetWebsite(),
		AllPlacesNoSearchAction:        strings.ToLower(request.GetAllPlacesNoSearchAction().String()),
	}

	if customGeolocationJson != nil {
		out.CustomGeolocation = customGeolocationJson
	}
	return out, nil
}
