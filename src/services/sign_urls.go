package services

import (
	"crypto/rsa"
	"strings"

	"bitbucket.org/calmisland/go-server-shared/errors"
	"bitbucket.org/calmisland/go-server-shared/requests/urlsigner"
	"bitbucket.org/calmisland/go-server-shared/services/aws/awscloudfront"
	"bitbucket.org/calmisland/go-server-shared/services/aws/awss3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type signInfo struct {
	AWSCloudFront *signInfoAWSCloudFront `json:"aws_cloudfront"`
	AWSS3         *signInfoAWSS3         `json:"aws_s3"`
}

type signInfoAWSCloudFront struct {
	KeyID          string `json:"keyId"`
	PrivateKey     string `json:"privateKey"`
	PrivateKeyPath string `json:"privateKeyPath"`
}

type signInfoAWSS3 struct {
	Region string `json:"region"`
}

var (
	urlSigner urlsigner.Signer
)

func initSignedUrls() {
	if config.Signing == nil {
		return
	}

	if config.Signing.AWSCloudFront != nil {
		setupAWSCloudFrontSigning(config.Signing.AWSCloudFront)
	} else if config.Signing.AWSS3 != nil {
		setupAWSS3Signing(config.Signing.AWSS3)
	}
}

func setupAWSCloudFrontSigning(signInfo *signInfoAWSCloudFront) {
	if len(signInfo.KeyID) == 0 {
		panic(errors.New("Missing AWS CloudFront key ID for product URL signing"))
	}

	var privateKey *rsa.PrivateKey
	var err error
	if len(signInfo.PrivateKey) > 0 {
		privateKey, err = awscloudfront.LoadPEMPrivKey(strings.NewReader(signInfo.PrivateKey))
	} else if len(signInfo.PrivateKeyPath) > 0 {
		privateKey, err = awscloudfront.LoadPEMPrivKeyFile(signInfo.PrivateKeyPath)
	} else {
		panic(errors.New("Missing AWS CloudFront private key for product URL signing"))
	}

	if err != nil {
		panic(errors.Wrap(err, "Failed to load the private key"))
	} else if privateKey == nil {
		panic(errors.New("Failed to load the private key"))
	}

	urlSigner = awscloudfront.NewURLSigner(signInfo.KeyID, privateKey)
}

func setupAWSS3Signing(signInfo *signInfoAWSS3) {
	if len(signInfo.Region) == 0 {
		panic(errors.New("Missing AWS S3 region for product URL signing"))
	}

	session, err := session.NewSession(&aws.Config{
		Region: aws.String(signInfo.Region),
	})
	if err != nil {
		panic(errors.Wrap(err, "Failed to create a new AWS session"))
	}

	urlSigner = awss3.NewURLSigner(session)
}

func signURL(url string, options urlsigner.SignOptions) (string, error) {
	if urlSigner != nil {
		return urlSigner.SignURL(url, options)
	}

	return url, nil
}
