package payapi

import (
	"io"
	"net/http"

	"github.com/mailru/easyjson/jwriter"
	"scholacantorum.org/orders/config"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/stripe"
)

// CreateCustomer handles POST /payapi/customer requests.
//
// Parameters:
//     auth:  authorization token
//     name:  new customer name
//     email:  new customer email
//     card:  new customer card source
//
// Returns:
//     If successful: status 200, and JSON object containing
//         customer:  Stripe customer ID
//         method:  Stripe payment method
//         description:  card description
//     If card error: status 400, and plain text error message
func CreateCustomer(tx db.Tx, w http.ResponseWriter, r *http.Request) {
	if r.FormValue("auth") != config.Get("galaAPIKey") {
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		return
	}
	name, email, card := r.FormValue("name"), r.FormValue("email"), r.FormValue("card")
	customer, pmtmeth, desc, problem := stripe.CreateCustomer(name, email, card)
	if problem != "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, problem)
		return
	}
	if customer == "" {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	jw := new(jwriter.Writer)
	jw.RawString(`{"customer":`)
	jw.String(customer)
	jw.RawString(`,"method":`)
	jw.String(pmtmeth)
	jw.RawString(`,"description":`)
	jw.String(desc)
	jw.RawByte('}')
	jw.DumpTo(w)
}

// UpdateCustomer handles POST /payapi/customer/${id} requests.
//
// Parameters:
//     auth:  authorization token
//     name:  new customer name
//     email:  new customer email
func UpdateCustomer(tx db.Tx, w http.ResponseWriter, r *http.Request, customerID string) {
	var ok bool
	var desc string

	if r.FormValue("auth") != config.Get("galaAPIKey") {
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		return
	}
	name, email, card := r.FormValue("name"), r.FormValue("email"), r.FormValue("card")
	if ok, desc = stripe.UpdateCustomer(customerID, name, email, card); !ok {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	if desc == "" {
		w.WriteHeader(http.StatusNoContent)
	} else {
		io.WriteString(w, desc)
	}
}
