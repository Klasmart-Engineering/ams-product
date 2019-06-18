package services

import "bitbucket.org/calmisland/go-server-shared/configs"

type curriculumConfig struct {
	CacheExpireMinutes int32     `json:"cacheExpireMinutes"`
	DownloadBaseURL    string    `json:"downloadBaseUrl"`
	Signing            *signInfo `json:"signing"`
}

var (
	config curriculumConfig
)

func init() {
	err := configs.LoadConfig("product", &config, true)
	if err != nil {
		panic(err)
	}

	if config.CacheExpireMinutes == 0 {
		config.CacheExpireMinutes = 15
	}

	initSignedUrls()
}
