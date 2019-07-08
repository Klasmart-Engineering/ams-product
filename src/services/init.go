package services

import (
	"bitbucket.org/calmisland/go-server-shared/v3/configs"
	"bitbucket.org/calmisland/go-server-shared/v3/errors"
)

type curriculumConfig struct {
	CacheExpireMinutes int32     `json:"cacheExpireMinutes"`
	DownloadBaseURL    string    `json:"downloadBaseUrl"`
	Signing            *signInfo `json:"signing"`
}

var (
	config curriculumConfig
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
