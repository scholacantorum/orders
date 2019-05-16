package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/model"
	"scholacantorum.org/orders/stripe"
)

// CreateProduct handles POST /api/product requests.
func CreateProduct(tx *sql.Tx, w http.ResponseWriter, r *http.Request) {
	var (
		session   *model.Session
		product   *model.Product
		err       error
		seenEvent = map[model.EventID]bool{}
	)
	if session = auth.GetSession(tx, w, r, model.PrivSetup); session == nil {
		return
	}
	if err = json.NewDecoder(r.Body).Decode(&product); err != nil {
		BadRequestError(tx, w, err.Error())
		return
	}
	if product.ID != 0 || product.StripeID == "" || product.Name == "" || product.TicketCount < 0 {
		BadRequestError(tx, w, "invalid parameters")
		return
	}
	if product.TicketCount > 0 {
		if len(product.Events) == 0 {
			BadRequestError(tx, w, "ticket products must have associated events")
			return
		}
	} else {
		if len(product.Events) != 0 {
			BadRequestError(tx, w, "only ticket products can have associated events")
			return
		}
	}
	for _, eid := range product.Events {
		if seenEvent[eid] {
			BadRequestError(tx, w, "duplicate event")
			return
		}
		if model.FetchEvent(tx, eid) == nil {
			BadRequestError(tx, w, "nonexistent event")
			return
		}
		seenEvent[eid] = true
	}
	if !stripe.ValidateProduct(product.StripeID) {
		BadRequestError(tx, w, "invalid Stripe product ID")
		return
	}
	product.Save(tx)
	commit(tx)
	log.Printf("%s CREATE PRODUCT %s", session.Username, toJSON(product))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}
