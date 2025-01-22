package config

import "errors"

type Postgres struct {
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	User            string `mapstructure:"user"`
	URL             string `mapstructure:"url"`
	Password        string `mapstructure:"password"`
	DatabaseName    string `mapstructure:"name"`
	DatabaseVersion uint   `mapstructure:"database_version"`
	MigrationPath   string `mapstructure:"migration_path"`
}

func (p *Postgres) Validate() error {
	if p.Host == "" {
		return errors.New("database Host is required")
	}
	if p.Port == "" {
		return errors.New("database Port is required")
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
	if p.MigrationPath == "" {
		return errors.New("database Migration Path is required")
	}
	return nil
}
