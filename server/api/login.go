package api

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rothskeller/json"

	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
)

// Login handles POST /api/login requests.
func Login(tx db.Tx, w http.ResponseWriter, r *http.Request) {
	var (
		password string
		session  model.Session
		resp     *http.Response
		jw       json.Writer
		err      error
	)
	if session.Username = strings.TrimSpace(r.FormValue("username")); session.Username == "" {
		BadRequestError(tx, w, "missing username")
		return
	}
	if password = r.FormValue("password"); password == "" {
		BadRequestError(tx, w, "missing password")
		return
	}
	resp, err = http.PostForm("https://scholacantorummembers.org/api/login/sso", url.Values{
		"username": []string{session.Username},
		"password": []string{password},
	})
	if err != nil {
		commit(tx)
		log.Printf("ERROR: can't contact members site for SSO: %s", err)
		http.Error(w, "500 SSO server error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusUnauthorized:
		commit(tx)
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		return
	default:
		commit(tx)
		log.Printf("ERROR: from members site SSO: %s", resp.Status)
		http.Error(w, "500 SSO server error", http.StatusInternalServerError)
		return
	}
	err = json.NewReader(resp.Body).Read(json.ObjectHandler(func(key string) json.Handlers {
		switch key {
		case "id":
			return json.IntHandler(func(i int) { session.Member = i })
		case "username":
			return json.IgnoreHandler()
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
		commit(tx)
		log.Printf("ERROR: bad JSON from members site SSO: %s", err)
		http.Error(w, "500 SSO server error", http.StatusInternalServerError)
		return
	}
	session.Privileges |= model.PrivLogin
	session.Expires = time.Now().Add(3 * time.Hour)
	session.Token = newToken()
	tx.SaveSession(&session)
	commit(tx)
	jw = json.NewWriter(w)
	jw.Object(func() {
		jw.Prop("token", session.Token)
		jw.Prop("username", session.Username)
		jw.Prop("privSetupOrders", session.Privileges&model.PrivSetupOrders != 0)
		jw.Prop("privViewOrders", session.Privileges&model.PrivViewOrders != 0)
		jw.Prop("privManageOrders", session.Privileges&model.PrivManageOrders != 0)
		jw.Prop("privInPersonSales", session.Privileges&model.PrivInPersonSales != 0)
		jw.Prop("privScanTickets", session.Privileges&model.PrivScanTickets != 0)
	})
	jw.Close()
}
