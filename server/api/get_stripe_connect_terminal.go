package api

import (
	"net/http"

	"github.com/rothskeller/json"

	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
	"scholacantorum.org/orders/stripe"
)

// GetStripeConnectTerminal gets a connection token from Stripe allowing a
// terminal to connect to our Stripe account.
func GetStripeConnectTerminal(tx db.Tx, w http.ResponseWriter, r *http.Request) {
	var (
		session *model.Session
		jw      json.Writer
		token   string
	)
	// Getting a connection token requires PrivInPersonSales.
	if session = auth.GetSession(tx, w, r, model.PrivInPersonSales); session == nil {
		return
	}
	commit(tx)
	if token = stripe.GetConnectionToken(); token == "" {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	jw = json.NewWriter(w)
	jw.String(token)
	jw.Close()
}
