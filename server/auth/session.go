package auth

//go:generate easyjson -lower_camel_case session.go

import (
	"log"
	"net/http"
	"net/url"

	"github.com/mailru/easyjson"

	"scholacantorum.org/orders/api"
	"scholacantorum.org/orders/model"
)

// Forbidden is the error returned when the calling session lacks the privileges
// needed for the call it issued.
var Forbidden = api.HTTPError(http.StatusForbidden, "403 Forbidden")

// ValidateSession decodes the session token in the request, if any, and sets
// the Session and Privileges fields of the request appropriately.  It returns
// an error only if the request contains a session token that is invalid.
func ValidateSession(r *api.Request) error {
	if r.Request.Header.Get("Auth") == "" {
		return nil
	}
	r.Tx.ExpireSessions()
	if r.Session = r.Tx.FetchSession(r.Request.Header.Get("Auth")); r.Session == nil {
		return api.HTTPError(http.StatusUnauthorized, "401 Unauthorized")
	}
	r.Privileges = r.Session.Privileges
	return nil
}

//easyjson:json
type ssoLogin struct {
	ID                int
	Username          string
	Name              string
	Email             string
	Address           string
	City              string
	State             string
	Zip               string
	PrivSetupOrders   bool
	PrivViewOrders    bool
	PrivManageOrders  bool
	PrivInPersonSales bool
	PrivScanTickets   bool
}

// GetSessionMembersAuth verifies that the provided auth token is a valid token
// for someone logged into the members site, and returns a pseudo-session with
// the corresponding data.
func GetSessionMembersAuth(r *api.Request, auth string) (err error) {
	var (
		resp     *http.Response
		ssoLogin ssoLogin
	)
	resp, err = http.Get("http://scholacantorummembers.org/api/login/sso?auth=" + url.QueryEscape(auth))
	if err != nil {
		log.Printf("ERROR: can't contact members site for SSO: %s", err)
		return api.HTTPError(http.StatusInternalServerError, "500 SSO server error")
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusUnauthorized:
		return api.HTTPError(http.StatusUnauthorized, "401 Unauthorized")
	default:
		log.Printf("ERROR: from members site SSO: %s", resp.Status)
		return api.HTTPError(http.StatusInternalServerError, "500 SSO server error")
	}
	if err = easyjson.UnmarshalFromReader(resp.Body, &ssoLogin); err != nil {
		log.Printf("ERROR: bad JSON from members site SSO: %s", err)
		return api.HTTPError(http.StatusInternalServerError, "500 SSO server error")
	}
	r.Session = &model.Session{
		Member:     ssoLogin.ID,
		Username:   ssoLogin.Username,
		Name:       ssoLogin.Name,
		Email:      ssoLogin.Email,
		Address:    ssoLogin.Address,
		City:       ssoLogin.City,
		State:      ssoLogin.State,
		Zip:        ssoLogin.Zip,
		Privileges: model.PrivLogin,
	}
	if ssoLogin.PrivSetupOrders {
		r.Session.Privileges |= model.PrivSetupOrders
	}
	if ssoLogin.PrivViewOrders {
		r.Session.Privileges |= model.PrivViewOrders
	}
	if ssoLogin.PrivManageOrders {
		r.Session.Privileges |= model.PrivManageOrders
	}
	if ssoLogin.PrivInPersonSales {
		r.Session.Privileges |= model.PrivInPersonSales
	}
	if ssoLogin.PrivScanTickets {
		r.Session.Privileges |= model.PrivScanTickets
	}
	r.Privileges = r.Session.Privileges
	return nil
}
