package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

const (
	httpPort   = "PORTS.HTTP.Port"
	grpcPort   = "PORTS.GRPC.Port"
	healthPort = "PORTS.Health.Port"
	dbPort     = "PORTS.Database.Port"
)

const (
	apifyKey            = "APIFY.KEY"
	apifyExtractorActor = "APIFY.ACTOR.EXTRACTOR.ID"
	apifyScraperActor   = "APIFY.ACTOR.SCRAPER.ID"
)

const (
	dbUser      = "DATABASE.User"
	dbPassword  = "DATABASE.Password"
	dbHost      = "DATABASE.Host"
	dbName      = "DATABASE.Name"
	dbMigration = "DATABASE.Migrations.Path"
	dbVersion   = "DATABASE.Version"
	dbURL       = "DATABASE.URL"
)

type Conf interface {
	Validate() error
}

type Config struct {
	Database Postgres `mapstructure:"database"`
	Ports    Ports    `mapstructure:"ports"`
	Apify    Apify    `mapstructure:"apify"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Validate() error {
	if err := c.Database.Validate(); err != nil {
		return err
	}
	if err := c.Ports.Validate(); err != nil {
		return err
	}
	if err := c.Apify.Validate(); err != nil {
		return err
	}
	return nil
}

func LoadDatabaseConfig() *Config {
	cfg := NewConfig()

	viper, err := viperSetup()
	if err != nil {
		log.Fatalf("Error setting up viper: %v", err)
	}

	getValues(viper, cfg)

	err = cfg.Database.Validate()
	if err != nil {
		log.Fatalf("Error validating config: %v", err)
		panic(err)
	}

	if cfg.Ports.DatabasePort == 0 {
		log.Fatalf("Database Port is required")
	}

	return cfg
}

func LoadConfig() *Config {
	cfg := NewConfig()

	viper, err := viperSetup()
	if err != nil {
		log.Fatalf("Error setting up viper: %v", err)
	}

	getValues(viper, cfg)

	err = cfg.Validate()
	if err != nil {
		log.Fatalf("Error validating config: %v", err)
		panic(err)
	}

	return cfg
}

func viperSetup() (*viper.Viper, error) {
	// Get config path from environment variable, fallback to default
	root := viper.New()
	root.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	root.AutomaticEnv()

	root.SetDefault(httpPort, 8080)
	root.SetDefault(grpcPort, 50051)
	root.SetDefault(healthPort, 8081)
	root.SetDefault(dbPort, 5432)

	root.SetDefault(apifyKey, "")
	root.SetDefault(apifyExtractorActor, "")
	root.SetDefault(apifyScraperActor, "")

	root.SetDefault(dbUser, "postgres")
	root.SetDefault(dbPassword, "postgres")
	root.SetDefault(dbHost, "localhost")
	root.SetDefault(dbName, "POIRawData")
	root.SetDefault(dbMigration, "db/migrations")
	root.SetDefault(dbVersion, 1)
	root.SetDefault(dbURL, "")

	return root, nil
}

func getValues(root *viper.Viper, cfg *Config) {

	cfg.Database.User = root.GetString(dbUser)
	cfg.Database.Password = root.GetString(dbPassword)
	cfg.Database.Host = root.GetString(dbHost)
	cfg.Database.DatabaseName = root.GetString(dbName)
	cfg.Database.MigrationsPath = root.GetString(dbMigration)
	cfg.Database.DatabaseVersion = root.GetUint(dbVersion)

	cfg.Ports.DatabasePort = root.GetInt(dbPort)
	cfg.Ports.GRPCPort = root.GetInt(grpcPort)
	cfg.Ports.HTTPPort = root.GetInt(httpPort)
	cfg.Ports.HealthPort = root.GetInt(healthPort)

	cfg.Apify.Key = root.GetString(apifyKey)
	cfg.Apify.ActorExtractorID = root.GetString(apifyExtractorActor)
	cfg.Apify.ActorScraperID = root.GetString(apifyScraperActor)

	cfg.Database.URL = fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Ports.DatabasePort,
		cfg.Database.DatabaseName,
	)
}
