package apify

import (
	"time"

	"apify-poi-data/internal/models"
)

/*
RunInitiated: Information about the run initiation.
	- Data: The main object containing all the run initiation information.
	  - ActID: Identifier for the act associated with the run.
	  - BuildID: Identifier for the build used in the run.
	  - BuildNumber: Build number used in the run.
	  - ContainerURL: URL of the container used in the run.
	  - DefaultDatasetID: Identifier for the default dataset used in the run.
	  - DefaultKeyValueStoreID: Identifier for the default key-value store used in the run.
	  - DefaultRequestQueueID: Identifier for the default request queue used in the run.
	  - FinishedAt: Timestamp when the run finished.
	  - ID: Unique identifier for the run.
	  - Meta: Meta information about the run, such as origin and user agent.
	  - Options: Options used for the run, such as build, timeout, memory, and disk usage.
	  - PricingInfo: Pricing information related to the run, including pricing model, unit price, and margin percentage.
	  - StartedAt: Timestamp when the run started.
	  - Stats: Various statistics about the run, including input length, reboot count, and resource usage.
	  - Status: Current status of the run.
	  - Usage: Detailed usage information, including compute units, dataset reads/writes, key-value store operations, and data transfer.
	  - UsageTotalUSD: Total usage cost in USD.
	  - UsageUSD: Detailed usage cost in USD.
	  - UserID: Identifier for the user who initiated the run.
*/
// RunInitiated represents the JSON response from the Apify API.
type RunInitiated struct {
	Data struct {
		ActID                  string      `json:"actId"`
		BuildID                string      `json:"buildId"`
		BuildNumber            string      `json:"buildNumber"`
		ContainerURL           string      `json:"containerUrl"`
		DefaultDatasetID       string      `json:"defaultDatasetId"`
		DefaultKeyValueStoreID string      `json:"defaultKeyValueStoreId"`
		DefaultRequestQueueID  string      `json:"defaultRequestQueueId"`
		FinishedAt             *time.Time  `json:"finishedAt"`
		ID                     string      `json:"id"`
		Meta                   Meta        `json:"meta"`
		Options                Options     `json:"options"`
		PricingInfo            PricingInfo `json:"pricingInfo"`
		StartedAt              time.Time   `json:"startedAt"`
		Stats                  Stats       `json:"stats"`
		Status                 string      `json:"status"`
		Usage                  Usage       `json:"usage"`
		UsageTotalUSD          float64     `json:"usageTotalUsd"`
		UsageUSD               UsageUSD    `json:"usageUsd"`
		UserID                 string      `json:"userId"`
	} `json:"data"`
}

type Meta struct {
	Origin    string `json:"origin"`
	UserAgent string `json:"userAgent"`
}

type Options struct {
	Build        string `json:"build"`
	DiskMbytes   int    `json:"diskMbytes"`
	MaxItems     int    `json:"maxItems"`
	MemoryMbytes int    `json:"memoryMbytes"`
	TimeoutSecs  int    `json:"timeoutSecs"`
}

type PricingInfo struct {
	ApifyMarginPercentage float64   `json:"apifyMarginPercentage"`
	CreatedAt             time.Time `json:"createdAt"`
	NotifiedAboutChangeAt time.Time `json:"notifiedAboutChangeAt"`
	PricePerUnitUSD       float64   `json:"pricePerUnitUsd"`
	PricingModel          string    `json:"pricingModel"`
	ReasonForChange       string    `json:"reasonForChange"`
	StartedAt             time.Time `json:"startedAt"`
	UnitName              string    `json:"unitName"`
}

type Stats struct {
	ComputeUnits   int `json:"computeUnits"`
	InputBodyLen   int `json:"inputBodyLen"`
	RebootCount    int `json:"rebootCount"`
	RestartCount   int `json:"restartCount"`
	ResurrectCount int `json:"resurrectCount"`
}

type Usage struct {
	ActorComputeUnits              int `json:"ACTOR_COMPUTE_UNITS"`
	DatasetReads                   int `json:"DATASET_READS"`
	DatasetWrites                  int `json:"DATASET_WRITES"`
	DataTransferExternalGBytes     int `json:"DATA_TRANSFER_EXTERNAL_GBYTES"`
	DataTransferInternalGBytes     int `json:"DATA_TRANSFER_INTERNAL_GBYTES"`
	KeyValueStoreLists             int `json:"KEY_VALUE_STORE_LISTS"`
	KeyValueStoreReads             int `json:"KEY_VALUE_STORE_READS"`
	KeyValueStoreWrites            int `json:"KEY_VALUE_STORE_WRITES"`
	ProxyResidentialTransferGBytes int `json:"PROXY_RESIDENTIAL_TRANSFER_GBYTES"`
	ProxySerps                     int `json:"PROXY_SERPS"`
	RequestQueueReads              int `json:"REQUEST_QUEUE_READS"`
	RequestQueueWrites             int `json:"REQUEST_QUEUE_WRITES"`
}

type UsageUSD struct {
	ActorComputeUnits              float64 `json:"ACTOR_COMPUTE_UNITS"`
	DatasetReads                   float64 `json:"DATASET_READS"`
	DatasetWrites                  float64 `json:"DATASET_WRITES"`
	DataTransferExternalGBytes     float64 `json:"DATA_TRANSFER_EXTERNAL_GBYTES"`
	DataTransferInternalGBytes     float64 `json:"DATA_TRANSFER_INTERNAL_GBYTES"`
	KeyValueStoreLists             float64 `json:"KEY_VALUE_STORE_LISTS"`
	KeyValueStoreReads             float64 `json:"KEY_VALUE_STORE_READS"`
	KeyValueStoreWrites            float64 `json:"KEY_VALUE_STORE_WRITES"`
	ProxyResidentialTransferGBytes float64 `json:"PROXY_RESIDENTIAL_TRANSFER_GBYTES"`
	ProxySerps                     float64 `json:"PROXY_SERPS"`
	RequestQueueReads              float64 `json:"REQUEST_QUEUE_READS"`
	RequestQueueWrites             float64 `json:"REQUEST_QUEUE_WRITES"`
}

/*
ResponseRunInfo represents the JSON response from the Apify API when retrieving information about a specific run.
It contains detailed information about the run, including its status, timing, options used, and various statistics.

- Data: The main object containing all the run information.
  - ID: Unique identifier for the run.
  - ActID: Identifier for the act associated with the run.
  - UserID: Identifier for the user who initiated the run.
  - StartedAt: Timestamp when the run started.
  - FinishedAt: Timestamp when the run finished.
  - Status: Current status of the run.
  - StatusMessage: Message describing the current status.
  - IsStatusMessageTerminal: Boolean indicating if the status message is terminal.
  - Meta: Meta information about the run, such as origin and user agent.
  - Stats: Various statistics about the run, including input length, reboot count, and resource usage.
  - Options: Options used for the run, such as build, timeout, memory, and disk usage.
  - BuildID: Identifier for the build used in the run.
  - ExitCode: Exit code of the run, if available.
  - DefaultKeyValueStoreID: Identifier for the default key-value store used in the run.
  - DefaultDatasetID: Identifier for the default dataset used in the run.
  - DefaultRequestQueueID: Identifier for the default request queue used in the run.
  - PricingInfo: Pricing information related to the run, including pricing model, unit price, and margin percentage.
  - BuildNumber: Build number used in the run.
  - ContainerUrl: URL of the container used in the run.
  - Usage: Detailed usage information, including compute units, dataset reads/writes, key-value store operations, and data transfer.
*/

// ResponseRunInfo represents the JSON response from the Apify API.
type ResponseRunInfo struct {
	Data GetData `json:"data"`
}

// Data represents the "data" object in the JSON response. When getting the run info, this is the object that contains the run information.
type GetData struct {
	ID                      string         `json:"id"`
	ActID                   string         `json:"actId"`
	UserID                  string         `json:"userId"`
	StartedAt               string         `json:"startedAt"`
	FinishedAt              string         `json:"finishedAt"`
	Status                  string         `json:"status"`
	StatusMessage           string         `json:"statusMessage"`
	IsStatusMessageTerminal *bool          `json:"isStatusMessageTerminal"`
	Meta                    GetMeta        `json:"meta"`
	Stats                   GetStats       `json:"stats"`
	Options                 GetOptions     `json:"options"`
	BuildID                 string         `json:"buildId"`
	ExitCode                *int           `json:"exitCode"`
	DefaultKeyValueStoreID  string         `json:"defaultKeyValueStoreId"`
	DefaultDatasetID        string         `json:"defaultDatasetId"`
	DefaultRequestQueueID   string         `json:"defaultRequestQueueId"`
	PricingInfo             GetPricingInfo `json:"pricingInfo"`
	BuildNumber             string         `json:"buildNumber"`
	ContainerUrl            string         `json:"containerUrl"`
	Usage                   GetUsage       `json:"usage"`
}

// GetMeta represents the "meta" field under "data". It meta information when getting run info.
type GetMeta struct {
	Origin    string `json:"origin"`
	UserAgent string `json:"userAgent"`
}

// GetStats represents the "stats" field under "data". It contains statistics about the run when getting run info.
type GetStats struct {
	InputBodyLen    int     `json:"inputBodyLen"`
	RebootCount     int     `json:"rebootCount"`
	RestartCount    int     `json:"restartCount"`
	DurationMillis  int     `json:"durationMillis"`
	ResurrectCount  int     `json:"resurrectCount"`
	RunTimeSecs     float64 `json:"runTimeSecs"`
	Metamorph       int     `json:"metamorph"`
	ComputeUnits    float64 `json:"computeUnits"`
	MemAvgBytes     float64 `json:"memAvgBytes"`
	MemMaxBytes     int     `json:"memMaxBytes"`
	MemCurrentBytes int     `json:"memCurrentBytes"`
	CpuAvgUsage     float64 `json:"cpuAvgUsage"`
	CpuMaxUsage     float64 `json:"cpuMaxUsage"`
	CpuCurrentUsage float64 `json:"cpuCurrentUsage"`
	NetRxBytes      int     `json:"netRxBytes"`
	NetTxBytes      int     `json:"netTxBytes"`
}

// GetOptions represents the "options" field under "data". It contains the options used when getting run info.
type GetOptions struct {
	Build        string `json:"build"`
	TimeoutSecs  int    `json:"timeoutSecs"`
	MemoryMbytes int    `json:"memoryMbytes"`
	MaxItems     int    `json:"maxItems"`
	DiskMbytes   int    `json:"diskMbytes"`
}

// GetPricingInfo represents the "pricingInfo" field under "data". It contains pricing information when getting run info.
type GetPricingInfo struct {
	PricingModel          string  `json:"pricingModel"`
	ReasonForChange       string  `json:"reasonForChange"`
	UnitName              string  `json:"unitName"`
	PricePerUnitUsd       float64 `json:"pricePerUnitUsd"`
	CreatedAt             string  `json:"createdAt"`
	StartedAt             string  `json:"startedAt"`
	ApifyMarginPercentage float64 `json:"apifyMarginPercentage"`
	NotifiedAboutChangeAt string  `json:"notifiedAboutChangeAt"`
}

// GetUsage represents the "usage" field under "data". It contains usage information when getting run info.
type GetUsage struct {
	ActorComputeUnits              float64 `json:"ACTOR_COMPUTE_UNITS"`
	DatasetReads                   int     `json:"DATASET_READS"`
	DatasetWrites                  int     `json:"DATASET_WRITES"`
	KeyValueStoreReads             int     `json:"KEY_VALUE_STORE_READS"`
	KeyValueStoreWrites            int     `json:"KEY_VALUE_STORE_WRITES"`
	KeyValueStoreLists             int     `json:"KEY_VALUE_STORE_LISTS"`
	RequestQueueReads              int     `json:"REQUEST_QUEUE_READS"`
	RequestQueueWrites             int     `json:"REQUEST_QUEUE_WRITES"`
	DataTransferInternalGBytes     float64 `json:"DATA_TRANSFER_INTERNAL_GBYTES"`
	DataTransferExternalGBytes     float64 `json:"DATA_TRANSFER_EXTERNAL_GBYTES"`
	ProxyResidentialTransferGBytes int     `json:"PROXY_RESIDENTIAL_TRANSFER_GBYTES"`
	ProxySerps                     int     `json:"PROXY_SERPS"`
}

/*
Root represents the top-level array of places.
Place holds the main fields for each place/object in the array.
AdditionalInfo represents the nested map under "additionalInfo".
Location holds latitude and longitude.
OpeningHour holds a single day/hours entry under "openingHours".
*/

// Root represents the top-level array of places.
type Root []models.Place
