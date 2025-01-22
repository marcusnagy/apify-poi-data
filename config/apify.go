package config

import "errors"

type Apify struct {
	Key              string `mapstructure:"key"`
	ActorExtractorID string `mapstructure:"actor_extractor_id"`
	ActorScraperID   string `mapstructure:"actor_scraper_id"`
}

func (a *Apify) Validate() error {
	if a.Key == "" {
		return errors.New("Apify Key is required")
	}
	if a.ActorExtractorID == "" {
		return errors.New("Apify Actor ID is required")
	}
	if a.ActorScraperID == "" {
		return errors.New("Apify Actor ID is required")
	}
	return nil
}
