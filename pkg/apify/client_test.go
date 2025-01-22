package apify

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"

	"apify-poi-data/internal/models"
)

var (
	client *Client
)

func init() {
	if err := godotenv.Load("../../.env"); err != nil {
		panic(err)
	}
	client = NewClient(os.Getenv("APIFY_KEY"), os.Getenv("APIFY_ACTOR_EXTRACTOR_ID"), os.Getenv("APIFY_ACTOR_SCRAPER_ID"))
}

func TestClient(t *testing.T) {
	t.Run("ExtractPOIs", func(t *testing.T) {
		payload := models.InputPayloadMaps{
			SearchStringsArray:        []string{"restaurant", "cafe"},
			Language:                  "en",
			CountryCode:               "se",
			City:                      "Gothenburg",
			PostalCode:                "41105",
			MaxCrawledPlacesPerSearch: 10,
			SkipClosedPlaces:          true,
		}

		log.Println("Extracting POIs...")

		resp := client.ExtractPOIs(payload, 1, false)
		select {
		case data := <-resp.Data:
			log.Printf("Data: %v", data)
			return
		case err := <-resp.Err:
			t.Errorf("Error: %v", err)
			return
		case <-time.After(10 * time.Minute):
			t.Errorf("Timeout")
		}
	})
}
