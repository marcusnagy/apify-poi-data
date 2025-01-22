package models

import "encoding/json"

type InputPayloadMaps struct {
	SearchStringsArray        []string            `json:"searchStringsArray,omitempty"`
	LocationQuery             string              `json:"locationQuery,omitempty"`
	MaxCrawledPlacesPerSearch int                 `json:"maxCrawledPlacesPerSearch,omitempty"`
	Language                  string              `json:"language,omitempty"`
	CountryCode               string              `json:"countryCode,omitempty"`
	City                      string              `json:"city,omitempty"`
	State                     string              `json:"state,omitempty"`
	PostalCode                string              `json:"postalCode,omitempty"`
	SkipClosedPlaces          bool                `json:"skipClosedPlaces,omitempty"`
	PlaceMinimumStars         string              `json:"placeMinimumStars,omitempty"`
	StartUrls                 []map[string]string `json:"startUrls,omitempty"`
	CustomGeolocation         json.RawMessage     `json:"customGeolocation,omitempty"`
}

type ScraperInputPayloadMaps struct {
	SearchStringsArray             []string            `json:"searchStringsArray,omitempty"`
	LocationQuery                  string              `json:"locationQuery,omitempty"`
	MaxCrawledPlacesPerSearch      int                 `json:"maxCrawledPlacesPerSearch,omitempty"`
	Language                       string              `json:"language,omitempty"`
	MaxImages                      int                 `json:"maxImages,omitempty"`
	ScrapeImageAuthors             bool                `json:"scrapeImageAuthors,omitempty"`
	OnlyDataFromSearchPage         bool                `json:"onlyDataFromSearchPage,omitempty"`
	IncludeWebResults              bool                `json:"includeWebResults,omitempty"`
	ScrapeDirectories              bool                `json:"scrapeDirectories,omitempty"`
	ScrapeTableReservationProvider bool                `json:"scrapeTableReservationProvider,omitempty"`
	MaxReviews                     int                 `json:"maxReviews,omitempty"`
	ReviewsStartDate               string              `json:"reviewsStartDate,omitempty"`
	ReviewsSort                    string              `json:"reviewsSort,omitempty"`
	ReviewsFilterString            string              `json:"reviewsFilterString,omitempty"`
	ReviewsOrigin                  string              `json:"reviewsOrigin,omitempty"`
	ScrapeReviewsPersonalData      bool                `json:"scrapeReviewsPersonalData,omitempty"`
	MaxQuestions                   int                 `json:"maxQuestions,omitempty"`
	Zoom                           int                 `json:"zoom,omitempty"`
	CountryCode                    string              `json:"countryCode,omitempty"`
	City                           string              `json:"city,omitempty"`
	State                          string              `json:"state,omitempty"`
	County                         string              `json:"county,omitempty"`
	PostalCode                     string              `json:"postalCode,omitempty"`
	CustomGeolocation              json.RawMessage     `json:"customGeolocation,omitempty"`
	CategoryFilterWords            []string            `json:"categoryFilterWords,omitempty"`
	SearchMatching                 string              `json:"searchMatching,omitempty"`
	PlaceMinimumStars              string              `json:"placeMinimumStars,omitempty"`
	SkipClosedPlaces               bool                `json:"skipClosedPlaces,omitempty"`
	Website                        string              `json:"website,omitempty"`
	StartUrls                      []map[string]string `json:"startUrls,omitempty"`
	AllPlacesNoSearchAction        string              `json:"allPlacesNoSearchAction,omitempty"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type OpeningHour struct {
	Day   string `json:"day"`
	Hours string `json:"hours"`
}

type Place struct {
	AdditionalInfo    json.RawMessage `json:"additionalInfo"`
	Address           string          `json:"address"`
	Categories        []string        `json:"categories"`
	CategoryName      string          `json:"categoryName"`
	CID               string          `json:"cid"`
	City              string          `json:"city"`
	ClaimThisBusiness bool            `json:"claimThisBusiness"`
	CountryCode       string          `json:"countryCode"`
	FID               string          `json:"fid"`
	GasPrices         json.RawMessage `json:"gasPrices"`
	GoogleFoodURL     *string         `json:"googleFoodUrl"` // nil if <nil>
	HotelAds          json.RawMessage `json:"hotelAds"`
	ImageCategories   []string        `json:"imageCategories"`
	ImageURL          string          `json:"imageUrl"`
	ImagesCount       int             `json:"imagesCount"`
	IsAdvertisement   bool            `json:"isAdvertisement"`
	Kgmid             string          `json:"kgmid"`
	Location          Location        `json:"location"`
	Neighborhood      *string         `json:"neighborhood"` // nil if <nil>
	OpeningHours      json.RawMessage `json:"openingHours"`
	PeopleAlsoSearch  json.RawMessage `json:"peopleAlsoSearch"`
	PermanentlyClosed bool            `json:"permanentlyClosed"`
	Phone             *string         `json:"phone"`
	PhoneUnformatted  *string         `json:"phoneUnformatted"`
	PlaceID           string          `json:"placeId"`
	PlacesTags        json.RawMessage `json:"placesTags"`
	PostalCode        string          `json:"postalCode"`
	Price             *string         `json:"price"` // nil if <nil>
	Rank              int             `json:"rank"`
	ReviewsCount      int             `json:"reviewsCount"`
	ReviewsTags       json.RawMessage `json:"reviewsTags"`
	ScrapedAt         string          `json:"scrapedAt"` // Or time.Time if you parse dates
	SearchPageURL     string          `json:"searchPageUrl"`
	SearchString      string          `json:"searchString"`
	State             *string         `json:"state"` // nil if <nil>
	Street            string          `json:"street"`
	SubTitle          string          `json:"subTitle"`
	TemporarilyClosed bool            `json:"temporarilyClosed"`
	Title             string          `json:"title"`
	TotalScore        float64         `json:"totalScore"`
	URL               string          `json:"url"`
	Website           *string         `json:"website"`
}

func (p Place) GetID() string {
	return p.PlaceID
}

func (p Place) GetName() string {
	return p.Title
}

func (p Place) GetType() string {
	return "Google Place Extractor"
}

type PopularTimes struct {
	Day              string `json:"day"`
	Hour             int    `json:"hour"`
	OccupancyPercent int    `json:"occupancyPercent"`
}

type ReviewDistribution struct {
	OneStar   int `json:"oneStar"`
	TwoStar   int `json:"twoStar"`
	ThreeStar int `json:"threeStar"`
	FourStar  int `json:"fourStar"`
	FiveStar  int `json:"fiveStar"`
}

type PlaceScraper struct {
	SearchString            string                    `json:"searchString"`
	Rank                    int                       `json:"rank"`
	SearchPageURL           string                    `json:"searchPageUrl"`
	SearchPageLoadedURL     *string                   `json:"searchPageLoadedUrl"`
	IsAdvertisement         bool                      `json:"isAdvertisement"`
	Title                   string                    `json:"title"`
	SubTitle                *string                   `json:"subTitle"`
	Description             *string                   `json:"description"`
	Price                   *string                   `json:"price"`
	CategoryName            string                    `json:"categoryName"`
	Address                 string                    `json:"address"`
	Neighborhood            *string                   `json:"neighborhood"`
	Street                  string                    `json:"street"`
	City                    string                    `json:"city"`
	PostalCode              string                    `json:"postalCode"`
	State                   *string                   `json:"state"`
	CountryCode             string                    `json:"countryCode"`
	Website                 *string                   `json:"website"`
	Phone                   *string                   `json:"phone"`
	PhoneUnformatted        *string                   `json:"phoneUnformatted"`
	ClaimThisBusiness       bool                      `json:"claimThisBusiness"`
	Location                Location                  `json:"location"`
	LocatedIn               *string                   `json:"locatedIn"`
	PlusCode                *string                   `json:"plusCode"`
	Menu                    *string                   `json:"menu"`
	TotalScore              float64                   `json:"totalScore"`
	PermanentlyClosed       bool                      `json:"permanentlyClosed"`
	TemporarilyClosed       bool                      `json:"temporarilyClosed"`
	PlaceID                 string                    `json:"placeId"`
	Categories              []string                  `json:"categories"`
	FID                     string                    `json:"fid"`
	CID                     string                    `json:"cid"`
	ReviewsCount            int                       `json:"reviewsCount"`
	ReviewsDistribution     ReviewDistribution        `json:"reviewsDistribution"`
	ImagesCount             int                       `json:"imagesCount"`
	ImageCategories         []string                  `json:"imageCategories"`
	ScrapedAt               string                    `json:"scrapedAt"`
	ReserveTableURL         *string                   `json:"reserveTableUrl"`
	GoogleFoodURL           *string                   `json:"googleFoodUrl"`
	HotelStars              *string                   `json:"hotelStars"`
	HotelDescription        *string                   `json:"hotelDescription"`
	CheckInDate             *string                   `json:"checkInDate"`
	CheckOutDate            *string                   `json:"checkOutDate"`
	SimilarHotelsNearby     json.RawMessage           `json:"similarHotelsNearby"`
	HotelReviewSummary      json.RawMessage           `json:"hotelReviewSummary"`
	HotelAds                json.RawMessage           `json:"hotelAds"`
	PopularTimesLiveText    *string                   `json:"popularTimesLiveText"`
	PopularTimesLivePercent *int                      `json:"popularTimesLivePercent"`
	PopularTimesHistogram   map[string][]PopularTimes `json:"popularTimesHistogram"`
	OpeningHours            json.RawMessage           `json:"openingHours"`
	PeopleAlsoSearch        json.RawMessage           `json:"peopleAlsoSearch"`
	PlacesTags              json.RawMessage           `json:"placesTags"`
	ReviewsTags             json.RawMessage           `json:"reviewsTags"`
	AdditionalInfo          json.RawMessage           `json:"additionalInfo"`
	GasPrices               json.RawMessage           `json:"gasPrices"`
	QuestionsAndAnswers     json.RawMessage           `json:"questionsAndAnswers"`
	UpdatesFromCustomers    json.RawMessage           `json:"updatesFromCustomers"`
	URL                     string                    `json:"url"`
	ImageURL                string                    `json:"imageUrl"`
	Kgmid                   string                    `json:"kgmid"`
	WebResults              json.RawMessage           `json:"webResults,omitempty"`
	ParentPlaceURL          *string                   `json:"parentPlaceUrl"`
	TableReservationLinks   json.RawMessage           `json:"tableReservationLinks,omitempty"`
	BookingLinks            json.RawMessage           `json:"bookingLinks,omitempty"`
	OrderBy                 json.RawMessage           `json:"orderBy,omitempty"`
	Images                  *string                   `json:"images"`
	ImageUrls               []string                  `json:"imageUrls"`
	Reviews                 json.RawMessage           `json:"reviews"`
	UserPlaceNote           json.RawMessage           `json:"userPlaceNote"`
	RestaurantData          json.RawMessage           `json:"restaurantData"`
	OwnerUpdates            json.RawMessage           `json:"ownerUpdates"`
}

func (p PlaceScraper) GetID() string {
	return p.PlaceID
}

func (p PlaceScraper) GetName() string {
	return p.Title
}

func (p PlaceScraper) GetType() string {
	return "Google Place Scraper"
}
