// Copyright (c) 2014 Canonical Ltd.
// Licensed under the GPLv3, see the COPYING file for details.

package textsecure

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
	log "github.com/sirupsen/logrus"
)

// Config holds application configuration settings
type Config struct {
	Tel                       string              `yaml:"tel"` // Our telephone number
	UUID                      string              `yaml:"uuid" default:"notset"`
	Server                    string              `yaml:"server"`                    // The TextSecure server URL
	RootCA                    string              `yaml:"rootCA"`                    // The TLS signing certificate of the server we connect to
	ProxyServer               string              `yaml:"proxy"`                     // HTTP Proxy URL if one is being used
	VerificationType          string              `yaml:"verificationType"`          // Code verification method during registration (SMS/VOICE/DEV)
	StorageDir                string              `yaml:"storageDir"`                // Directory for the persistent storage
	UnencryptedStorage        bool                `yaml:"unencryptedStorage"`        // Whether to store plaintext keys and session state (only for development)
	StoragePassword           string              `yaml:"storagePassword"`           // Password to the storage
	LogLevel                  string              `yaml:"loglevel"`                  // Verbosity of the logging messages
	UserAgent                 string              `yaml:"userAgent"`                 // Override for the default HTTP User Agent header field
	AlwaysTrustPeerID         bool                `yaml:"alwaysTrustPeerID"`         // Workaround until proper handling of peer reregistering with new ID.
	AccountCapabilities       AccountCapabilities `yaml:"accountCapabilities"`       // Account Attrributes are used in order to track the support of different function for signal
	DiscoverableByPhoneNumber bool                `yaml:"discoverableByPhoneNumber"` // If the user should be found by his phone number
	ProfileKey                []byte              `yaml:"profileKey"`                // The profile key is used in many places to encrypt the avatar, name etc and also in groupsv2 context
	Name                      string              `yaml:"name"`
}

// TODO: some race conditions to be solved
func checkUUID(cfg *Config) *Config {
	defer func(cfg *Config) *Config {
		log.Debugln("[textsecure] missing my uuid defer")
		recover()
		UUID, err := GetMyUUID()
		if err != nil {
			log.Debugln("[textsecure] missing my uuid", err)
			return cfg
		}
		cfg.UUID = UUID
		return cfg
	}(cfg)
	if cfg.UUID == "notset" {
		log.Debugln("[textsecure] missing my uuid notset")
	}
	return cfg

}

// ReadConfig reads a YAML config file
func ReadConfig(fileName string) (*Config, error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = yaml.Unmarshal(b, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// WriteConfig saves a config to a file
func WriteConfig(filename string, cfg *Config) error {
	b, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0600)
}

// loadConfig gets the config via the client and makes sure
// that for unset values sane defaults are used
func loadConfig() (*Config, error) {
	cfg, err := client.GetConfig()

	if err != nil {
		return nil, err
	}

	if cfg.Server == "" {
		cfg.Server = "https://textsecure-service.whispersystems.org:443"
	}

	if cfg.VerificationType == "" {
		cfg.VerificationType = "sms"
	}

	if cfg.StorageDir == "" {
		cfg.StorageDir = ".storage"
	}

	return cfg, nil
}
