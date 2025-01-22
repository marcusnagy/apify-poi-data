package converter

import (
	"apify-poi-data/internal/models"
	tripadvisor_v1 "apify-poi-data/proto/apify/tripsadvisor/v1"
)

func convertTripAdvisorURLsToUrls(urls []*tripadvisor_v1.StartUrl) []models.URLItem {
	var result []models.URLItem
	if urls == nil {
		return result
	}
	for _, url := range urls {
		result = append(result, models.URLItem{
			URL: url.GetUrl(),
		})
	}
	return result
}

func SearchRequestToTripAdvisorInputPayload(request *tripadvisor_v1.SearchRequest) models.TripAdvisorInput {

	return models.TripAdvisorInput{
		Query:                   request.GetQuery(),
		StartURLs:               convertTripAdvisorURLsToUrls(request.GetStartUrls()),
		MaxItemsPerQuery:        int(request.GetMaxItemsPerQuery()),
		IncludeTags:             request.GetIncludeTags(),
		IncludeNearbyResults:    request.GetIncludeNearbyResults(),
		IncludeRestaurants:      request.GetIncludeRestaurants(),
		IncludeAttractions:      request.GetIncludeAttractions(),
		IncludeHotels:           request.GetIncludeHotels(),
		IncludeVacationRentals:  request.GetIncludeVacationRentals(),
		CheckInDate:             request.GetCheckInDate(),
		CheckOutDate:            request.GetCheckOutDate(),
		IncludePriceOffers:      request.GetIncludePriceOffers(),
		IncludeAiReviewsSummary: request.GetIncludeAiReviewsSummary(),
		Language:                request.GetLanguage(),
		Currency:                request.GetCurrency(),
	}
}
