package models

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	AllPOIsSearch = "all_places_no_search"
)

// POI is a Point of Interest
type shortType struct {
	Type         string `json:"type,omitempty"`         // Unique for Tripadvisor
	Kgmid        string `json:"kgmid,omitempty"`        // Unique for Google Maps
	SearchString string `json:"searchString,omitempty"` // Unique for Google Maps

}

// ParsePOIsFromJSON reads the array of POIs, then decides how to decode each object.
func ParsePOIsFromJSON(data []byte) ([]POI, error) {
	// First unmarshal into a slice of RawMessages
	var rawItems []json.RawMessage
	if err := json.Unmarshal(data, &rawItems); err != nil {
		return nil, fmt.Errorf("error unmarshaling top-level array: %w", err)
	}

	results := make([]POI, 0, len(rawItems))

	for _, raw := range rawItems {
		// Quick parse just the "type"
		var st shortType
		if err := json.Unmarshal(raw, &st); err != nil {
			return nil, fmt.Errorf("error reading type field: %w", err)
		}

		if strings.ToLower(st.SearchString) == AllPOIsSearch {
			fmt.Printf("Parsing POI for PlaceScraper; all_places_no_search\n")
			var g PlaceScraper
			if err := json.Unmarshal(raw, &g); err != nil {
				return nil, fmt.Errorf("error unmarshaling GoogleMaps: %w", err)
			}
			results = append(results, &g)
		} else if st.Type != "" && st.Kgmid == "" {
			fmt.Printf("Parsing POI for Tripadvisor; type=%s\n", st.Type)
			switch st.Type {
			case "HOTEL":
				var h Hotel
				if err := json.Unmarshal(raw, &h); err != nil {
					return nil, fmt.Errorf("error unmarshaling HOTEL: %w", err)
				}
				results = append(results, &h)

			case "RESTAURANT":
				var r Restaurant
				if err := json.Unmarshal(raw, &r); err != nil {
					return nil, fmt.Errorf("error unmarshaling RESTAURANT: %w", err)
				}
				results = append(results, &r)

			case "ATTRACTION":
				var a Attraction
				if err := json.Unmarshal(raw, &a); err != nil {
					return nil, fmt.Errorf("error unmarshaling ATTRACTION: %w", err)
				}
				results = append(results, &a)

			default:
				// either skip or store in a generic type
				// for now, we'll just skip
				fmt.Printf("Skipping unrecognized type: %s\n", st.Type)
			}
		} else if st.Type == "" && st.Kgmid != "" {
			fmt.Printf("Parsing POI for Google Maps; kgmid=%s\n", st.Kgmid)
			var g Place
			if err := json.Unmarshal(raw, &g); err != nil {
				return nil, fmt.Errorf("error unmarshaling GoogleMaps: %w", err)
			}
			results = append(results, &g)
		}
	}

	return results, nil
}

// Helper: safely dereference pointers
func ValueOrZero(ptr *int) int {
	if ptr == nil {
		return 0
	}
	return *ptr
}

func ValueOrEmpty(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
