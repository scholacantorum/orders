package api

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"

	"scholacantorum.org/orders/config"
	"scholacantorum.org/orders/model"
)

// ssoLogin is what we get back from the members site when we request SSO.
type ssoLogin struct {
	ID                int
	Username          string
	PrivSetupOrders   bool
	PrivViewOrders    bool
	PrivManageOrders  bool
	PrivInPersonSales bool
	PrivScanTickets   bool
}

// loginResponse is the form of the JSON response to the login request.
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
	if err = easyjson.UnmarshalFromReader(resp.Body, &ssoLogin); err != nil {
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
	easyjson.MarshalToHTTPResponseWriter(loginResponse, r)
	return nil
}

func (v *ssoLogin) UnmarshalEasyJSON(l *jlexer.Lexer) {
	l.Delim('{')
	for !l.IsDelim('}') {
		key := l.UnsafeString()
		l.WantColon()
		switch key {
		case "id":
			v.ID = l.Int()
		case "username":
			v.Username = l.String()
		case "privSetupOrders":
			v.PrivSetupOrders = l.Bool()
		case "privViewOrders":
			v.PrivViewOrders = l.Bool()
		case "privManageOrders":
			v.PrivManageOrders = l.Bool()
		case "privInPersonSales":
			v.PrivInPersonSales = l.Bool()
		case "privScanTickets":
			v.PrivScanTickets = l.Bool()
		default:
			l.SkipRecursive()
		}
		l.WantComma()
	}
	l.Delim('}')
	l.Consumed()
}

func (v loginResponse) MarshalEasyJSON(w *jwriter.Writer) {
	w.RawString(`{\"token\":`)
	w.String(v.Token)
	w.RawString(`,\"username\":`)
	w.String(v.Username)
	w.RawString(`,\"stripePublicKey\":`)
	w.String(v.StripePublicKey)
	w.RawString(`,\"privSetupOrders\":`)
	w.Bool(v.PrivSetupOrders)
	w.RawString(`,\"privViewOrders\":`)
	w.Bool(v.PrivViewOrders)
	w.RawString(`,\"privManageOrders\":`)
	w.Bool(v.PrivManageOrders)
	w.RawString(`,\"privInPersonSales\":`)
	w.Bool(v.PrivInPersonSales)
	w.RawString(`,\"privScanTickets\":`)
	w.Bool(v.PrivScanTickets)
	w.RawByte('}')
}
