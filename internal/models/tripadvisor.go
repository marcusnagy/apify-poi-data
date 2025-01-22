package models

// POI is a common interface satisfied by Hotel, Restaurant, Attraction, etc.
type POI interface {
	GetID() string
	GetType() string
	GetName() string
}

// BasePOI holds the most common fields across all POI types.
type BasePOI struct {
	ID            string   `json:"id"`
	Type          string   `json:"type"`
	Category      string   `json:"category"`
	Subcategories []string `json:"subcategories,omitempty"`
	Name          string   `json:"name"`
	Description   *string  `json:"description,omitempty"`
	Image         *string  `json:"image,omitempty"`
	PhotoCount    *int     `json:"photoCount,omitempty"`

	// Basic ranking info
	RankingPosition    *int     `json:"rankingPosition,omitempty"`
	RankingString      *string  `json:"rankingString,omitempty"`
	RankingDenominator *string  `json:"rankingDenominator,omitempty"`
	Rating             *float64 `json:"rating,omitempty"`
	RawRanking         *float64 `json:"rawRanking,omitempty"`

	Phone         *string     `json:"phone,omitempty"`
	Address       *string     `json:"address,omitempty"`
	AddressObj    *AddressObj `json:"addressObj,omitempty"`
	LocalName     *string     `json:"localName,omitempty"`
	LocalAddress  *string     `json:"localAddress,omitempty"`
	LocalLangCode *string     `json:"localLangCode,omitempty"`
	Email         *string     `json:"email,omitempty"`

	Latitude       float64 `json:"latitude,omitempty"`
	Longitude      float64 `json:"longitude,omitempty"`
	LocationString *string `json:"locationString,omitempty"`

	WebURL  *string `json:"webUrl,omitempty"`
	Website *string `json:"website,omitempty"`

	// "rankingDenominator" often appears as a string (e.g. "177"), you might parse it to int.

	NeighborhoodLocations []NeighborhoodLocation `json:"neighborhoodLocations,omitempty"`
	NearestMetroStations  []NearestMetroStation  `json:"nearestMetroStations,omitempty"`
	AncestorLocations     []AncestorLocation     `json:"ancestorLocations,omitempty"`

	RatingHistogram *RatingHistogram `json:"ratingHistogram,omitempty"`

	// Booking / OfferGroup appear mostly for attractions, sometimes hotels, etc.
	Booking    *Booking    `json:"booking,omitempty"`
	OfferGroup *OfferGroup `json:"offerGroup,omitempty"`

	// Subtype often appears for attractions or restaurants
	Subtype []string `json:"subtype,omitempty"`

	// Photos is a list of image URLs
	Photos []string `json:"photos,omitempty"`

	TravelerChoiceAward *string `json:"travelerChoiceAward,omitempty"`
	Input               *string `json:"input,omitempty"` // e.g. "Gothenburg"

	ReviewTags []ReviewTag `json:"reviewTags,omitempty"`

	// Flags
	IsNearbyResult bool    `json:"isNearbyResult,omitempty"`
	IsClosed       bool    `json:"isClosed,omitempty"`
	IsLongClosed   bool    `json:"isLongClosed,omitempty"`
	OpenNowText    *string `json:"openNowText,omitempty"`

	// For certain POIs that have "hours" info
	Hours            *Hours  `json:"hours,omitempty"`
	MenuWebURL       *string `json:"menuWebUrl,omitempty"`
	OwnersTopReasons any     `json:"ownersTopReasons,omitempty"` // can store as interface{} or define a struct
}

// Satisfy the POI interface on BasePOI
func (b *BasePOI) GetID() string {
	return b.ID
}

func (b *BasePOI) GetType() string {
	return b.Type
}

func (b *BasePOI) GetName() string {
	return b.Name
}

// AddressObj is for "addressObj"
//
//	e.g.  "addressObj": {
//	        "street1": "Kungstorget 9",
//	        "street2": "",
//	        "city": "Gothenburg",
//	        "state": null,
//	        "country": "Sweden",
//	        "postalcode": "411 17"
//	      }
type AddressObj struct {
	Street1    *string `json:"street1,omitempty"`
	Street2    *string `json:"street2,omitempty"`
	City       *string `json:"city,omitempty"`
	State      *string `json:"state,omitempty"`
	Country    *string `json:"country,omitempty"`
	Postalcode *string `json:"postalcode,omitempty"`
}

// NeighborhoodLocation and NearestMetroStation are lists of objects
type NeighborhoodLocation struct {
	ID   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

type NearestMetroStation struct {
	Name         *string  `json:"name,omitempty"`
	LocalName    *string  `json:"localName,omitempty"`
	Address      *string  `json:"address,omitempty"`
	LocalAddress *string  `json:"localAddress,omitempty"`
	Lines        []string `json:"lines,omitempty"`
	Latitude     *float64 `json:"latitude,omitempty"`
	Longitude    *float64 `json:"longitude,omitempty"`
	Distance     *float64 `json:"distance,omitempty"`
}

// AncestorLocation is for the array like:
// "ancestorLocations": [
//
//	{"id": "189894","name": "Gothenburg","abbreviation": null,"subcategory": "City"}
//	...
//
// ]
type AncestorLocation struct {
	ID           *string `json:"id,omitempty"`
	Name         *string `json:"name,omitempty"`
	Abbreviation *string `json:"abbreviation,omitempty"`
	Subcategory  *string `json:"subcategory,omitempty"`
}

// RatingHistogram has "count1", "count2", ...
type RatingHistogram struct {
	Count1 int `json:"count1"`
	Count2 int `json:"count2"`
	Count3 int `json:"count3"`
	Count4 int `json:"count4"`
	Count5 int `json:"count5"`
}

// Booking might show up on attractions/hotels
type Booking struct {
	Provider *string `json:"provider,omitempty"`
	URL      *string `json:"url,omitempty"`
}

// OfferGroup often includes lowestPrice, offerList, etc.
type OfferGroup struct {
	LowestPrice *string `json:"lowestPrice,omitempty"`
	OfferList   []Offer `json:"offerList,omitempty"`
}

type Offer struct {
	URL             *string `json:"url,omitempty"`
	Price           *string `json:"price,omitempty"`
	RoundedUpPrice  *string `json:"roundedUpPrice,omitempty"`
	OfferType       *string `json:"offerType,omitempty"`
	Title           *string `json:"title,omitempty"`
	ProductCode     *string `json:"productCode,omitempty"`
	Partner         *string `json:"partner,omitempty"`
	ImageURL        *string `json:"imageUrl,omitempty"`
	Description     *string `json:"description,omitempty"`
	PrimaryCategory *string `json:"primaryCategory,omitempty"`
}

// ReviewTag has "text" and "reviews"
type ReviewTag struct {
	Text    string `json:"text"`
	Reviews int    `json:"reviews"`
}

// Hours represents something like:
//
//	"hours": {
//	  "weekRanges": [
//	      [], [ { open:540, close:1320 } ], ...
//	  ],
//	  "timezone": "Pacific/Auckland"
//	}
type Hours struct {
	WeekRanges [][]TimeRange `json:"weekRanges,omitempty"`
	Timezone   *string       `json:"timezone,omitempty"`
}

// TimeRange is for "open" and "close" in minutes from midnight, plus "openHours" string
type TimeRange struct {
	Open       int    `json:"open,omitempty"`
	OpenHours  string `json:"openHours,omitempty"`
	Close      int    `json:"close,omitempty"`
	CloseHours string `json:"closeHours,omitempty"`
}

// ---------- HOTEL ----------
type Hotel struct {
	BasePOI

	HotelClass    *string `json:"hotelClass,omitempty"`
	NumberOfRooms *int    `json:"numberOfRooms,omitempty"`

	// Some hotels show "amenities"
	Amenities []string `json:"amenities,omitempty"`

	// Standard "priceLevel" or "priceRange" for hotels?
	PriceLevel *string `json:"priceLevel,omitempty"`
	PriceRange *string `json:"priceRange,omitempty"`

	CheckInDate  *string `json:"checkInDate,omitempty"`
	CheckOutDate *string `json:"checkOutDate,omitempty"`

	// roomTips is an array of objects
	RoomTips []RoomTip `json:"roomTips,omitempty"`
}

// RoomTip objects
type RoomTip struct {
	User        interface{} `json:"user,omitempty"` // seems to be null in samples
	Type        *string     `json:"type,omitempty"`
	Text        *string     `json:"text,omitempty"`
	Rating      *string     `json:"rating,omitempty"`
	ReviewID    *string     `json:"reviewId,omitempty"`
	ID          *string     `json:"id,omitempty"`
	CreatedTime *string     `json:"createdTime,omitempty"`
}

// ---------- RESTAURANT ----------
type Restaurant struct {
	BasePOI

	// e.g. "cuisines": [ "European", "Fusion", "New Zealand" ]
	Cuisines []string `json:"cuisines,omitempty"`

	// "dietaryRestrictions": ["Vegetarian friendly", "Vegan options", ...]
	DietaryRestrictions []string `json:"dietaryRestrictions,omitempty"`

	// "establishmentTypes": ["Restaurants"]
	EstablishmentTypes []string `json:"establishmentTypes,omitempty"`

	// "features": ["Reservations","Outdoor Seating","Seating","Serves Alcohol",...]
	Features []string `json:"features,omitempty"`

	// "mealTypes": ["Lunch","Dinner","Brunch","Drinks"]
	MealTypes []string `json:"mealTypes,omitempty"`

	PriceLevel *string `json:"priceLevel,omitempty"`
	PriceRange *string `json:"priceRange,omitempty"`

	// e.g. "isClaimedIcon", "isClaimedText", "orderOnline"
	IsClaimedIcon *bool       `json:"isClaimedIcon,omitempty"`
	IsClaimedText *string     `json:"isClaimedText,omitempty"`
	OrderOnline   interface{} `json:"orderOnline,omitempty"` // often an empty list
}

// ---------- ATTRACTION ----------
type Attraction struct {
	BasePOI
	// For attractions, you often see "booking", "offerGroup", "subtype"
	// but we already have them in BasePOI.
}

// TripAdvisorInput represents the input payload for the Tripadvisor Scraper,
// matching the schema described in your OpenAPI components.
type TripAdvisorInput struct {
	Query                   string    `json:"query,omitempty"`
	StartURLs               []URLItem `json:"startUrls,omitempty"`
	MaxItemsPerQuery        int       `json:"maxItemsPerQuery,omitempty"`
	IncludeTags             bool      `json:"includeTags,omitempty"`
	IncludeNearbyResults    bool      `json:"includeNearbyResults,omitempty"`
	IncludeAttractions      bool      `json:"includeAttractions,omitempty"`
	IncludeRestaurants      bool      `json:"includeRestaurants,omitempty"`
	IncludeHotels           bool      `json:"includeHotels,omitempty"`
	IncludeVacationRentals  bool      `json:"includeVacationRentals,omitempty"`
	CheckInDate             string    `json:"checkInDate,omitempty"`
	CheckOutDate            string    `json:"checkOutDate,omitempty"`
	IncludePriceOffers      bool      `json:"includePriceOffers,omitempty"`
	IncludeAiReviewsSummary bool      `json:"includeAiReviewsSummary,omitempty"`
	Language                string    `json:"language,omitempty"`
	Currency                string    `json:"currency,omitempty"`
}

// URLItem represents a single URL object within StartURLs.
type URLItem struct {
	URL string `json:"url"`
}
