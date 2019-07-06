package proxyconfig

import (
	"fmt"
	"io/ioutil"

	"github.com/prometheus/prometheus/config"

	"github.com/jacksontj/promxy/pkg/servergroup"

	yaml "gopkg.in/yaml.v2"
)

var DefaultPromxyConfig = PromxyConfig{
	BoundaryTimeWorkaround: true,
}

// ConfigFromFile loads a config file at path
func ConfigFromFile(path string) (*Config, error) {
	// load the config file
	cfg := &Config{
		PromConfig:   config.DefaultConfig,
		PromxyConfig: DefaultPromxyConfig,
	}
	configBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Error loading config: %v", err)
	}
	err = yaml.Unmarshal([]byte(configBytes), &cfg)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling config: %v", err)
	}

	return cfg, nil
}

// Config is the entire config file. This includes both the Prometheus Config
// as well as the Promxy config. This is done by "inline-ing" the promxy
// config into the prometheus config under the "promxy" key
type Config struct {
	// Prometheus configs -- this includes configurations for
	// recording rules, alerting rules, etc.
	PromConfig config.Config `yaml:",inline"`

	// Promxy specific configuration -- under its own namespace
	PromxyConfig `yaml:"promxy"`
}

// PromxyConfig is the configuration for Promxy itself
type PromxyConfig struct {
	// Config for each of the server groups promxy is configured to aggregate
	ServerGroups []*servergroup.Config `yaml:"server_groups"`

	// BoundaryTimeWorkaround enables a workaround to prometheus' internal boundary
	// times being un-supported serverside. Newer versions of prometheus (2.11+)
	// don't need this workaround enabled as it was worked around server-side.
	BoundaryTimeWorkaround bool `yaml:"boundary_time_workaround"`
}
