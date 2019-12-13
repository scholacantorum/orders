package posapi

import (
	"net/http"

	"github.com/mailru/easyjson/jwriter"

	"scholacantorum.org/orders/api"
	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/model"
	"scholacantorum.org/orders/stripe"
)

// GetStripeConnectTerminal gets a connection token from Stripe allowing a
// terminal to connect to our Stripe account.
func GetStripeConnectTerminal(r *api.Request) error {
	var (
		token string
		jw    jwriter.Writer
	)
	r.Tx.Commit()
	// Getting a connection token requires PrivInPersonSales.
	if r.Privileges&model.PrivInPersonSales == 0 {
		return auth.Forbidden
	}
	if token = stripe.GetConnectionToken(); token == "" {
		return api.HTTPError(http.StatusInternalServerError, "500 Internal Server Error")
	}
	r.Header().Set("Content-Type", "application/json")
	jw.String(token)
	jw.DumpTo(r)
	return nil
}
