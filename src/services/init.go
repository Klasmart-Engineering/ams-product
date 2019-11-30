package services

import (
	"bitbucket.org/calmisland/go-server-configs/configs"
	"github.com/calmisland/go-errors"
)

// ProductConfig is configuration for products.
type ProductConfig struct {
	DownloadBaseURL string    `json:"downloadBaseUrl"`
	Signing         *SignInfo `json:"signing"`
}

var (
	productConfig ProductConfig
)

// Initialize Initialize
func Initialize(config ProductConfig) error {
	productConfig = config
	return initSignedUrls()
}

// InitializeFromConfigs InitializeFromConfigs
func InitializeFromConfigs() error {
	var config ProductConfig
	err := configs.LoadConfig("product", &config, true)
	if err != nil {
		return errors.Wrap(err, "Failed to read the product configuration file")
	}

	return Initialize(config)
}
