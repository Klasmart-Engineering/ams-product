package globals

import (
	"errors"

	"bitbucket.org/calmisland/go-server-requests/tokens/accesstokens"
)

var (
	// AccessTokenValidator is the access token validator.
	AccessTokenValidator accesstokens.Validator
)

// Verify verifies if all variables have been properly set.
func Verify() {
	if AccessTokenValidator == nil {
		panic(errors.New("The access token validator has not been set"))
	}
}
