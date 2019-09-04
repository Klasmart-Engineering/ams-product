package services

import (
	"bitbucket.org/calmisland/go-server-configs/configs"
	"github.com/calmisland/go-errors"
)

type productConfig struct {
	CacheExpireMinutes int32     `json:"cacheExpireMinutes"`
	DownloadBaseURL    string    `json:"downloadBaseUrl"`
	Signing            *signInfo `json:"signing"`
}

var (
	config productConfig
)

// InitializeFromConfigs InitializeFromConfigs
func InitializeFromConfigs() error {
	err := configs.LoadConfig("product", &config, true)
	if err != nil {
		return errors.Wrap(err, "Failed to read the product configuration file")
	}

	if config.CacheExpireMinutes == 0 {
		config.CacheExpireMinutes = 15
	}

	return initSignedUrls()
}
