package api

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"scholacantorum.org/orders/config"
	"scholacantorum.org/orders/model"
)

type ssoLogin struct {
	ID                int
	Username          string
	PrivSetupOrders   bool
	PrivViewOrders    bool
	PrivManageOrders  bool
	PrivInPersonSales bool
	PrivScanTickets   bool
}
type loginResponse struct {
	Token             string
	Username          string
	StripePublicKey   string
	PrivSetupOrders   bool
	PrivViewOrders    bool
	PrivManageOrders  bool
	PrivInPersonSales bool
	PrivScanTickets   bool
}

// Login handles POST /api/login requests.
func Login(r *Request) error {
	var (
		password      string
		buf           []byte
		ssoLogin      ssoLogin
		session       model.Session
		resp          *http.Response
		loginResponse loginResponse
		err           error
	)
	if session.Username = strings.TrimSpace(r.FormValue("username")); session.Username == "" {
		return errors.New("missing username")
	}
	if password = r.FormValue("password"); password == "" {
		return errors.New("missing password")
	}
	resp, err = http.PostForm("https://scholacantorummembers.org/api/login/sso", url.Values{
		"username": []string{session.Username},
		"password": []string{password},
	})
	if err != nil {
		log.Printf("ERROR: can't contact members site for SSO: %s", err)
		return HTTPError(http.StatusInternalServerError, "500 SSO server error")
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusUnauthorized:
		return HTTPError(http.StatusUnauthorized, "401 Unauthorized")
	default:
		log.Printf("ERROR: from members site SSO: %s", resp.Status)
		return HTTPError(http.StatusInternalServerError, "500 SSO server error")
	}
	if buf, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Printf("ERROR: bad response from members site SSO: %s", err)
		return HTTPError(http.StatusInternalServerError, "500 SSO server error")
	}
	if err = ssoLogin.UnmarshalJSON(buf); err != nil {
		log.Printf("ERROR: bad JSON from members site SSO: %s", err)
		return HTTPError(http.StatusInternalServerError, "500 SSO server error")
	}
	session.Member = ssoLogin.ID
	if ssoLogin.PrivSetupOrders {
		session.Privileges |= model.PrivSetupOrders
	}
	if ssoLogin.PrivViewOrders {
		session.Privileges |= model.PrivViewOrders
	}
	if ssoLogin.PrivManageOrders {
		session.Privileges |= model.PrivManageOrders
	}
	if ssoLogin.PrivInPersonSales {
		session.Privileges |= model.PrivInPersonSales
	}
	if ssoLogin.PrivScanTickets {
		session.Privileges |= model.PrivScanTickets
	}
	session.Privileges |= model.PrivLogin
	session.Expires = time.Now().Add(3 * time.Hour)
	session.Token = NewToken()
	r.Tx.SaveSession(&session)
	r.Tx.Commit()
	loginResponse.Token = session.Token
	loginResponse.Username = session.Username
	loginResponse.StripePublicKey = config.Get("stripePublicKey")
	loginResponse.PrivSetupOrders = ssoLogin.PrivSetupOrders
	loginResponse.PrivViewOrders = ssoLogin.PrivViewOrders
	loginResponse.PrivManageOrders = ssoLogin.PrivManageOrders
	loginResponse.PrivInPersonSales = ssoLogin.PrivInPersonSales
	loginResponse.PrivScanTickets = ssoLogin.PrivScanTickets
	buf, _ = loginResponse.MarshalJSON()
	r.Write(buf)
	return nil
}
