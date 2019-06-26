package auth

import (
	"log"
	"net/http"
	"net/url"

	"github.com/rothskeller/json"

	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
)

// GetSession verifies that the HTTP request identifies a valid session, and
// that session has all of the privileges specified in the priv bitmask.  (The
// mask may be zero to check for a valid session without requiring specific
// privileges.)  If so, GetSession returns the session details.  If not,
// GetSession emits an appropriate error, rolls back the transaction, and
// returns nil.
func GetSession(tx db.Tx, w http.ResponseWriter, r *http.Request, priv model.Privilege) (session *model.Session) {
	tx.ExpireSessions()
	if session = tx.FetchSession(r.Header.Get("Auth")); session == nil {
		tx.Rollback()
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		return nil
	}
	if session.Privileges&priv != priv {
		tx.Rollback()
		http.Error(w, "403 Forbidden", http.StatusForbidden)
		return nil
	}
	return session
}

// GetSessionMembersAuth verifies that the provided auth token is a valid token
// for someone logged into the members site, and returns a pseudo-session with
// the corresponding data.
func GetSessionMembersAuth(tx db.Tx, w http.ResponseWriter, r *http.Request, auth string) (session *model.Session) {
	var (
		resp *http.Response
		err  error
	)
	resp, err = http.Get("http://scholacantorummembers.org/api/login/sso?auth=" + url.QueryEscape(auth))
	if err != nil {
		tx.Rollback()
		log.Printf("ERROR: can't contact members site for SSO: %s", err)
		http.Error(w, "500 SSO server error", http.StatusInternalServerError)
		return nil
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusUnauthorized:
		tx.Rollback()
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		return nil
	default:
		tx.Rollback()
		log.Printf("ERROR: from members site SSO: %s", resp.Status)
		http.Error(w, "500 SSO server error", http.StatusInternalServerError)
		return nil
	}
	session = new(model.Session)
	err = json.NewReader(resp.Body).Read(json.ObjectHandler(func(key string) json.Handlers {
		switch key {
		case "id":
			return json.IntHandler(func(i int) { session.Member = i })
		case "username":
			return json.StringHandler(func(s string) { session.Username = s })
		case "name":
			return json.StringHandler(func(s string) { session.Name = s })
		case "email":
			return json.StringHandler(func(s string) { session.Email = s })
		case "address":
			return json.StringHandler(func(s string) { session.Address = s })
		case "city":
			return json.StringHandler(func(s string) { session.City = s })
		case "state":
			return json.StringHandler(func(s string) { session.State = s })
		case "zip":
			return json.StringHandler(func(s string) { session.Zip = s })
		case "privSetupOrders":
			return json.BoolHandler(func(b bool) {
				if b {
					session.Privileges |= model.PrivSetupOrders
				}
			})
		case "privViewOrders":
			return json.BoolHandler(func(b bool) {
				if b {
					session.Privileges |= model.PrivViewOrders
				}
			})
		case "privManageOrders":
			return json.BoolHandler(func(b bool) {
				if b {
					session.Privileges |= model.PrivManageOrders
				}
			})
		case "privInPersonSales":
			return json.BoolHandler(func(b bool) {
				if b {
					session.Privileges |= model.PrivInPersonSales
				}
			})
		case "privScanTickets":
			return json.BoolHandler(func(b bool) {
				if b {
					session.Privileges |= model.PrivScanTickets
				}
			})
		default:
			return json.RejectHandler()
		}
	}))
	if err != nil {
		tx.Rollback()
		log.Printf("ERROR: bad JSON from members site SSO: %s", err)
		http.Error(w, "500 SSO server error", http.StatusInternalServerError)
		return nil
	}
	session.Privileges |= model.PrivLogin
	return session
}

// HasSession returns whether the HTTP request identifies a session.  (The
// session is not necessarily valid.)
func HasSession(r *http.Request) bool {
	return r.Header.Get("Auth") != ""
}
