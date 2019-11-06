package middlewares

import (
	"context"

	"github.com/calmisland/go-errors"

	"bitbucket.org/calmisland/go-server-account/accountdatabase"
	"bitbucket.org/calmisland/go-server-auth/authmiddlewares"
	"bitbucket.org/calmisland/go-server-requests/apierrors"
	"bitbucket.org/calmisland/go-server-requests/apirequests"
	"bitbucket.org/calmisland/go-server-requests/apirouter"
	"bitbucket.org/calmisland/go-server-requests/tokens/accesstokens"
)

// ValidateSession is holding session validation for data-reporting with L&P
func ValidateSession(validator accesstokens.Validator, validatorLP accesstokens.Validator, rejectRequest bool) apirouter.MiddlewareFunc {
	if validator == nil || validatorLP == nil {
		panic(errors.New("The validators cannot be nil"))
	}
	authMiddleware := authmiddlewares.ValidateSession(validator, rejectRequest)
	authLPMiddleware := authmiddlewares.ValidateSession(validatorLP, rejectRequest)
	return func(ctx context.Context, req *apirequests.Request, resp *apirequests.Response, next apirouter.NextFunc) error {
		if _, ok := req.GetHeader("X-Access-Token"); ok {
			return authMiddleware(ctx, req, resp, next)
		} else if sessionHeader, ok := req.GetHeader("X-LP-Session"); ok {
			// Here is a hack to reuse ValidateSession method
			req.Headers["X-Access-Token"] = []string{sessionHeader}
			return authLPMiddleware(ctx, req, resp, func() error {
				if req.Session != nil && len(req.Session.ID) > 0 && req.Session.Data != nil && len(req.Session.Data.AccountID) > 0 && len(req.Session.Data.Email) > 0 {
					db, err := accountdatabase.GetDatabase()
					if err != nil {
						return resp.SetServerError(err)
					}
					// This operation is time consuming (1 extra request per endpoint)
					// Add a new layer of security by checking the AccountID
					accInfo, err := db.GetAccountInfo(req.Session.Data.AccountID)
					if err != nil {
						return resp.SetServerError(err)
					}
					if accInfo != nil && accInfo.Email == req.Session.Data.Email {
						return next()
					}
					// If we don't want to reject the request, we just continue directly
					if !rejectRequest {
						return next()
					}
					// If the session token doesn't match criteria, always unauthorized
					return resp.SetClientError(apierrors.ErrorUnauthorized)
				}
				return next()
			})
		} else {
			if rejectRequest {
				return resp.SetClientError(apierrors.ErrorUnauthorized)
			}
			return next()
		}
	}
}
