package apify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"apify-poi-data/internal/models"
)

const (
	// TODO: Make the Actor ID more generic and modifiable.
	RunTaskURL         = "https://api.apify.com/v2/actor-tasks/%s/runs?maxItems=%d" // RunTaskURL is the URL for running the Apify API
	BaseRunTaskURL     = "https://api.apify.com/v2/actor-tasks/%s/runs"             // BaseURL is the base URL for the Apify API
	PollingURL         = "https://api.apify.com/v2/actor-runs/%s"                   // PollingURL is the URL for polling the Apify API
	GetDatasetItemsURL = "https://api.apify.com/v2/actor-runs/%s/dataset/items"     // GetDatasetItemsURL is the URL for getting the dataset items from the Apify API
)

type Poll struct {
	Data chan []byte
	Err  chan error
}

type POIResponse struct {
	Data chan []models.POI
	Err  chan error
}

type Client struct {
	client           *http.Client
	actorExtractorID string
	actorScraperID   string
	key              string
}

// NewClient creates a new Apify client.
func NewClient(key, actorExtractorID, actorScraperID string) *Client {
	return &Client{
		client:           &http.Client{},
		actorExtractorID: actorExtractorID,
		actorScraperID:   actorScraperID,
		key:              key,
	}
}

// GetDataset gets the dataset from the Apify API.
// returns an array of items.
func (c *Client) GetDataset(id string) ([]byte, error) {
	completeURL := fmt.Sprintf(GetDatasetItemsURL, id)
	req, err := http.NewRequest("GET", completeURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.key))

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

// pollingWithBackoff polls the Apify API with backoff.
// The backoff is doubled each time the polling fails.
// The backoff is capped at 1 minute.
func (c *Client) pollingWithBackoff(id string, p *Poll, shouldBackoff bool) {
	// Polling the Apify API
	completeURL := fmt.Sprintf(PollingURL, id)
	backoff := time.Second

	for {
		req, err := http.NewRequest("GET", completeURL, nil)
		if err != nil {
			p.Err <- err
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.key))

		resp, err := c.client.Do(req)
		if err != nil {
			p.Err <- err
			return
		}

		// Read the response body
		respBody, err := readResponseBody(resp)
		if err != nil {
			p.Err <- err
			return
		}

		var response ResponseRunInfo
		if err := json.Unmarshal(respBody, &response); err != nil {
			p.Err <- err
			return
		}

		switch response.Data.Status {
		case "SUCCEEDED":
			// Get the dataset
			dataset, err := c.GetDataset(response.Data.ID)
			fmt.Printf("Successfully retrieved dataset; STATUS=%s\n", response.Data.Status)
			if err != nil {
				p.Err <- err
			}
			p.Data <- dataset
			return
		case "ABORTED":
			p.Err <- fmt.Errorf("run aborted")
			return
		case "FAILED":
			p.Err <- fmt.Errorf("run failed")
			return
		default:
			// Do nothing
		}

		time.Sleep(backoff)
		if shouldBackoff {
			backoff *= 2
			if backoff > time.Minute {
				backoff = time.Minute
			}
		}
	}
}

func (c *Client) newRequest(method, url string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.key))
	return req, nil
}

// TripAdvisorPOIs extracts POIs from the Apify API.
// Required that a task is created in the Apify API, or via the console. The task ID is required to run the task.
func (c *Client) TripAdvisorPOIs(payload models.TripAdvisorInput, maxResults int, backoff bool) POIResponse {
	completeURL := fmt.Sprintf(RunTaskURL, c.actorExtractorID, maxResults)
	resp := POIResponse{
		Data: make(chan []models.POI, 1),
		Err:  make(chan error, 1),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		resp.Err <- err
		return resp
	}

	req, err := c.newRequest("POST", completeURL, body)
	if err != nil {
		resp.Err <- err
		return resp
	}

	r, err := c.client.Do(req)
	if err != nil {
		resp.Err <- err
		return resp
	}

	// Read the response body
	respBody, err := readResponseBody(r)
	if err != nil {
		resp.Err <- err
		return resp
	}

	var unmarshaledResponse RunInitiated
	// Unmarshal the response body
	if err := json.Unmarshal(respBody, &unmarshaledResponse); err != nil {
		var errResponse map[string]any
		if err := json.Unmarshal(respBody, &errResponse); err != nil {
			resp.Err <- err
			return resp
		}
		resp.Err <- fmt.Errorf("error response: %v", errResponse)
	}

	p := &Poll{
		Data: make(chan []byte, 1),
		Err:  make(chan error, 1),
	}

	go c.pollingWithBackoff(unmarshaledResponse.Data.ID, p, backoff)

	go func() {
		select {
		case data := <-p.Data:
			fmt.Println("Data received")
			poilist, err := models.ParsePOIsFromJSON(data)
			if err != nil {
				resp.Err <- err
				return
			}
			resp.Data <- poilist
		case err := <-p.Err:
			resp.Err <- err
			return
		}
	}()

	return resp
}

// ExtractPOIs extracts POIs from the Apify API.
// Required that a task is created in the Apify API, or via the console. The task ID is required to run the task.
func (c *Client) ExtractPOIs(payload models.InputPayloadMaps, maxResults int, backoff bool) POIResponse {

	completeURL := fmt.Sprintf(RunTaskURL, c.actorExtractorID, maxResults)
	resp := POIResponse{
		Data: make(chan []models.POI, 1),
		Err:  make(chan error, 1),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		resp.Err <- err
		return resp
	}

	// print body to see what is being sent
	fmt.Println(string(body))

	req, err := c.newRequest("POST", completeURL, body)
	if err != nil {
		resp.Err <- err
		return resp
	}
	r, err := c.client.Do(req)
	if err != nil {
		resp.Err <- err
		return resp
	}

	// Read the response body
	respBody, err := readResponseBody(r)
	if err != nil {
		resp.Err <- err
		return resp
	}

	var unmarshaledResponse RunInitiated
	// Unmarshal the response body
	if err := json.Unmarshal(respBody, &unmarshaledResponse); err != nil {
		var errResponse map[string]any
		if err := json.Unmarshal(respBody, &errResponse); err != nil {
			resp.Err <- err
			return resp
		}
		resp.Err <- fmt.Errorf("error response: %v", errResponse)
	}

	p := &Poll{
		Data: make(chan []byte, 1),
		Err:  make(chan error, 1),
	}

	go c.pollingWithBackoff(unmarshaledResponse.Data.ID, p, backoff)

	go func() {
		select {
		case data := <-p.Data:
			fmt.Println("Data received")
			poilist, err := models.ParsePOIsFromJSON(data)
			if err != nil {
				resp.Err <- err
				return
			}
			resp.Data <- poilist
		case err := <-p.Err:
			resp.Err <- err
			return
		}
	}()

	return resp
}

// ExtractPOIs extracts POIs from the Apify API.
// Required that a task is created in the Apify API, or via the console. The task ID is required to run the task.
func (c *Client) ScrapePOIs(payload models.ScraperInputPayloadMaps, backoff bool) POIResponse {

	completeURL := fmt.Sprintf(BaseRunTaskURL, c.actorScraperID)
	resp := POIResponse{
		Data: make(chan []models.POI, 1),
		Err:  make(chan error, 1),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		resp.Err <- err
		return resp
	}

	// print body to see what is being sent
	fmt.Println(string(body))

	req, err := c.newRequest("POST", completeURL, body)
	if err != nil {
		resp.Err <- err
		return resp
	}
	r, err := c.client.Do(req)
	if err != nil {
		resp.Err <- err
		return resp
	}

	// Read the response body
	respBody, err := readResponseBody(r)
	if err != nil {
		resp.Err <- err
		return resp
	}

	var unmarshaledResponse RunInitiated
	// Unmarshal the response body
	if err := json.Unmarshal(respBody, &unmarshaledResponse); err != nil {
		var errResponse map[string]any
		if err := json.Unmarshal(respBody, &errResponse); err != nil {
			resp.Err <- err
			return resp
		}
		resp.Err <- fmt.Errorf("error response: %v", errResponse)
	}

	p := &Poll{
		Data: make(chan []byte, 1),
		Err:  make(chan error, 1),
	}

	go c.pollingWithBackoff(unmarshaledResponse.Data.ID, p, backoff)

	go func() {
		select {
		case data := <-p.Data:
			fmt.Println("Data received")
			poilist, err := models.ParsePOIsFromJSON(data)
			if err != nil {
				resp.Err <- err
				return
			}
			resp.Data <- poilist
		case err := <-p.Err:
			resp.Err <- err
			return
		}
	}()

	return resp
}

// readResponseBody reads the response body and checks the status code.
// If the status code is not OK or Created, it parses and returns the error with the error body.
// If successful, it returns the response body as a byte slice.
func readResponseBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errResponse map[string]any
		if err := json.Unmarshal(body, &errResponse); err != nil {
			return nil, fmt.Errorf("unexpected status code: %d, error reading body: %v", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("unexpected status code: %d, error response: %v", resp.StatusCode, errResponse)
	}

	return body, nil
}
