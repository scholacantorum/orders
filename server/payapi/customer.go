package payapi

import (
	"net/http"

	"scholacantorum.org/orders/config"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/stripe"
)

// UpdateCustomer handles POST /payapi/customer/${id} requests.
//
// Parameters:
//     auth:  authorization token
//     name:  new customer name
//     email:  new customer email
func UpdateCustomer(tx db.Tx, w http.ResponseWriter, r *http.Request, customerID string) {
	if r.FormValue("auth") != config.Get("galaAPIKey") {
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		return
	}
	name, email := r.FormValue("name"), r.FormValue("email")
	if !stripe.UpdateCustomer(customerID, name, email) {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
