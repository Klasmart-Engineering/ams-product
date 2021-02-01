package v1

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

// SignInfo is signing information.
type SignInfo struct {
	AWSCloudFront *SignInfoAWSCloudFront `json:"aws_cloudfront"`
	AWSS3         *SignInfoAWSS3         `json:"aws_s3"`
}

// SignInfoAWSCloudFront is information for signing CloudFront requests.
type SignInfoAWSCloudFront struct {
	KeyID          string `json:"keyId"`
	PrivateKey     string `json:"privateKey"`
	PrivateKeyPath string `json:"privateKeyPath"`
}

// SignInfoAWSS3 is information for signing S3 requests.
type SignInfoAWSS3 struct {
	Region string `json:"region"`
}

var (
	urlSigner urlsign.Signer
)

func initSignedUrls() error {
	if productConfig.Signing == nil {
		return nil
	}

	if productConfig.Signing.AWSCloudFront != nil {
		return setupAWSCloudFrontSigning(productConfig.Signing.AWSCloudFront)
	} else if productConfig.Signing.AWSS3 != nil {
		return setupAWSS3Signing(productConfig.Signing.AWSS3)
	}
	return nil
}

func setupAWSCloudFrontSigning(signInfo *SignInfoAWSCloudFront) error {
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

func setupAWSS3Signing(signInfo *SignInfoAWSS3) error {
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
		signResult, err := urlSigner.SignURL(url, options)
		if err != nil {
			return "", err
		}

		return signResult.URL, nil
	}

	return url, nil
}
