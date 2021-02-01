package helpers

import (
	"net/url"

	"github.com/labstack/echo/v4"
)

func GetArrayQueryParams(c echo.Context, key string) ([]string, error) {
	queryParams := c.QueryParams()
	ret, err := url.ParseQuery(queryParams.Encode())
	if err != nil {
		return nil, err
	}

	values := ret[key]

	return values, nil
}
