package v1

import (
	"fmt"

	"bitbucket.org/calmisland/go-server-configs/configs"
	"github.com/calmisland/go-errors"
)

// ProductConfig is configuration for products.
type ProductConfig struct {
	DownloadBaseURL string    `json:"downloadBaseUrl"`
	Signing         *SignInfo `json:"signing"`
}

type ProductConfigEnv struct {
	ContentBaseURL          string `env:"CONTENT_BASE_URL"`
	ContentS3Region         string `env:"CONTENT_S3_REGION"`
	ContentCFKeyID          string `env:"CONTENT_CF_KEY_ID"`
	ContentCFPrivateKey     string `env:"CONTENT_CF_PRIVATE_KEY"`
	ContentCFPrivateKeyPath string `env:"CONTENT_CF_PRIVATE_KEY_PATH"`
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

func InitializeFromEnvs() error {
	var configEnv ProductConfigEnv

	err := configs.ReadEnvConfig(&configEnv)
	if err != nil {
		return errors.Wrap(err, "Failed to read the product environmental configuration")
	}

	if len(configEnv.ContentBaseURL) == 0 {
		return fmt.Errorf("failed to read the product environmental configuration: CONTENT_BASE_URL")
	}

	var config ProductConfig
	config.DownloadBaseURL = configEnv.ContentBaseURL

	var signInfo SignInfo
	config.Signing = &signInfo

	if len(configEnv.ContentS3Region) != 0 {
		var s3Info SignInfoAWSS3

		s3Info.Region = configEnv.ContentS3Region

		signInfo.AWSS3 = &s3Info
	} else {
		if len(configEnv.ContentCFKeyID) == 0 {
			return fmt.Errorf("failed to read the product environmental configuration: CONTENT_CF_KEY_ID")
		}

		if len(configEnv.ContentCFPrivateKey) == 0 && len(configEnv.ContentCFPrivateKeyPath) == 0 {
			return fmt.Errorf("failed to read the product environmental configuration: CONTENT_CF_PRIVATE_KEY_PATH or CONTENT_CF_PRIVATE_KEY")
		}

		var cfInfo SignInfoAWSCloudFront

		cfInfo.KeyID = configEnv.ContentCFKeyID
		cfInfo.PrivateKey = configEnv.ContentCFPrivateKey
		cfInfo.PrivateKeyPath = configEnv.ContentCFPrivateKeyPath
	}

	// fmt.Printf("%+v\n", config)

	return Initialize(config)
}
