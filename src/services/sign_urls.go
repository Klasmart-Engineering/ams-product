package services

import (
	"crypto/rsa"
	"strings"

	"bitbucket.org/calmisland/go-server-aws/awscloudfront"
	"bitbucket.org/calmisland/go-server-aws/awss3"
	"bitbucket.org/calmisland/go-server-requests/urlsign"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/calmisland/go-errors"
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
	urlSigner urlsign.Signer
)

func initSignedUrls() error {
	if config.Signing == nil {
		return nil
	}

	if config.Signing.AWSCloudFront != nil {
		return setupAWSCloudFrontSigning(config.Signing.AWSCloudFront)
	} else if config.Signing.AWSS3 != nil {
		return setupAWSS3Signing(config.Signing.AWSS3)
	}
	return nil
}

func setupAWSCloudFrontSigning(signInfo *signInfoAWSCloudFront) error {
	if len(signInfo.KeyID) == 0 {
		return errors.New("Missing AWS CloudFront key ID for product URL signing")
	}

	var privateKey *rsa.PrivateKey
	var err error
	if len(signInfo.PrivateKey) > 0 {
		privateKey, err = awscloudfront.LoadPEMPrivKey(strings.NewReader(signInfo.PrivateKey))
	} else if len(signInfo.PrivateKeyPath) > 0 {
		privateKey, err = awscloudfront.LoadPEMPrivKeyFile(signInfo.PrivateKeyPath)
	} else {
		return errors.New("Missing AWS CloudFront private key for product URL signing")
	}

	if err != nil {
		return errors.Wrap(err, "Failed to load the private key")
	} else if privateKey == nil {
		return errors.New("Failed to load the private key")
	}

	urlSigner = awscloudfront.NewURLSigner(signInfo.KeyID, privateKey)
	return nil
}

func setupAWSS3Signing(signInfo *signInfoAWSS3) error {
	if len(signInfo.Region) == 0 {
		return errors.New("Missing AWS S3 region for product URL signing")
	}

	session, err := session.NewSession(&aws.Config{
		Region: aws.String(signInfo.Region),
	})
	if err != nil {
		return errors.Wrap(err, "Failed to create a new AWS session")
	}

	urlSigner = awss3.NewURLSigner(session)
	return nil
}

func signURL(url string, options urlsign.SignOptions) (string, error) {
	if urlSigner != nil {
		return urlSigner.SignURL(url, options)
	}

	return url, nil
}
