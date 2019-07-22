package services

import (
	"fmt"
	"net/url"
	"time"

	"bitbucket.org/calmisland/go-server-requests/urlsign"
	"bitbucket.org/calmisland/go-server-utils/urlutils"
)

// GetProgramIconURL returns the URL for a specific program icon.
func GetProgramIconURL(productID string) (string, error) {
	productID = url.PathEscape(productID)
	iconFileName := fmt.Sprintf("%s.png", productID)
	iconURL := urlutils.JoinURL(config.DownloadBaseURL, "icons/products", iconFileName)
	return signURL(iconURL, urlsign.SignOptions{
		Expires: time.Now().Add(30 * time.Minute),
	})
}
