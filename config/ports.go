package config

import "errors"

type Ports struct {
	GRPCPort     int `mapstructure:"grpc_port"`
	HTTPPort     int `mapstructure:"http_port"`
	HealthPort   int `mapstructure:"health_port"`
	DatabasePort int `mapstructure:"database_port"`
}

func (p *Ports) Validate() error {
	if p.GRPCPort == 0 {
		return errors.New("GRPC Port is required")
	}
	if p.HTTPPort == 0 {
		return errors.New("HTTP Port is required")
	}
	if p.HealthPort == 0 {
		return errors.New("health Port is required")
	}
	if p.DatabasePort == 0 {
		return errors.New("database Port is required")
	}
	return nil
}
