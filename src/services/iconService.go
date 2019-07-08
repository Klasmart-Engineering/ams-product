package services

import (
	"fmt"
	"net/url"
	"time"

	"bitbucket.org/calmisland/go-server-shared/v3/requests/urlsigner"
	"bitbucket.org/calmisland/go-server-shared/v3/utils/urlutils"
)

// GetProgramIconURL returns the URL for a specific program icon.
func GetProgramIconURL(productID string) (string, error) {
	productID = url.PathEscape(productID)
	iconFileName := fmt.Sprintf("%s.png", productID)
	iconURL := urlutils.JoinURL(config.DownloadBaseURL, "icons/products", iconFileName)
	return signURL(iconURL, urlsigner.SignOptions{
		Expires: time.Now().Add(30 * time.Minute),
	})
}
