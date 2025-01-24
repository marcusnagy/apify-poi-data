package config

import "errors"

type Postgres struct {
	Host            string `mapstructure:"host"`
	User            string `mapstructure:"user"`
	URL             string `mapstructure:"url"`
	Password        string `mapstructure:"password"`
	DatabaseName    string `mapstructure:"name"`
	DatabaseVersion uint   `mapstructure:"database_version"`
	MigrationsPath  string `mapstructure:"migrations_path"`
}

func (p *Postgres) Validate() error {
	if p.Host == "" {
		return errors.New("database Host is required")
	}
	if p.User == "" {
		return errors.New("database User is required")
	}
	if p.Password == "" {
		return errors.New("database Password is required")
	}
	if p.DatabaseName == "" {
		return errors.New("database Name is required")
	}
	if p.DatabaseVersion == 0 {
		return errors.New("database Version is required")
	}
	if p.MigrationsPath == "" {
		return errors.New("database Migration Path is required")
	}
	return nil
}
