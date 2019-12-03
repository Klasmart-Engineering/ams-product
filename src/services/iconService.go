package services

import (
	"fmt"
	"net/url"
	"time"

	"bitbucket.org/calmisland/go-server-requests/urlsign"
	"bitbucket.org/calmisland/go-server-utils/urlutils"
)

const (
	iconDownloadURLDuration = 30 * time.Minute
)

// GetContentIconURL returns the URL for a specific content icon.
func GetContentIconURL(contentID string) (string, error) {
	contentID = url.PathEscape(contentID)
	iconFileName := fmt.Sprintf("%s.png", contentID)
	iconURL := urlutils.JoinURL(productConfig.DownloadBaseURL, "icons/contents", iconFileName)
	urlExpireTime := time.Now().Add(iconDownloadURLDuration)

	return signURL(iconURL, urlsign.SignOptions{
		Expires: urlExpireTime,
	})
}

// GetProductIconURL returns the URL for a specific product icon.
func GetProductIconURL(productID string) (string, error) {
	productID = url.PathEscape(productID)
	iconFileName := fmt.Sprintf("%s.png", productID)
	iconURL := urlutils.JoinURL(productConfig.DownloadBaseURL, "icons/products", iconFileName)
	urlExpireTime := time.Now().Add(iconDownloadURLDuration)

	return signURL(iconURL, urlsign.SignOptions{
		Expires: urlExpireTime,
	})
}
