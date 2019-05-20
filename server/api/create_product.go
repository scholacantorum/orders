package api

import (
	"encoding/json"
	"log"
	"net/http"

	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
)

// CreateProduct handles POST /api/product requests.
func CreateProduct(tx db.Tx, w http.ResponseWriter, r *http.Request) {
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
	if product.ID == "" || product.Name == "" || product.Type == "" || product.TicketCount < 0 {
		BadRequestError(tx, w, "invalid parameters")
		return
	}
	if tx.FetchProduct(product.ID) != nil {
		BadRequestError(tx, w, "duplicate product ID")
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
	for _, event := range product.Events {
		if event.ID == "" {
			BadRequestError(tx, w, "invalid event")
			return
		}
		if seenEvent[event.ID] {
			BadRequestError(tx, w, "duplicate event")
			return
		}
		if tx.FetchEvent(event.ID) == nil {
			BadRequestError(tx, w, "nonexistent event")
			return
		}
		seenEvent[event.ID] = true
	}
	for i, sku := range product.SKUs {
		if sku.Price < 0 || (!sku.SalesStart.IsZero() && !sku.SalesEnd.IsZero() && !sku.SalesEnd.After(sku.SalesStart)) {
			BadRequestError(tx, w, "invalid SKU parameters")
			return
		}
		for j := 0; j < i; j++ {
			prev := product.SKUs[j]
			if prev.MembersOnly == sku.MembersOnly && prev.Coupon == sku.Coupon && overlappingDates(prev, sku) {
				BadRequestError(tx, w, "overlapping SKUs")
				return
			}
		}
	}
	tx.SaveProduct(product)
	commit(tx)
	log.Printf("%s CREATE PRODUCT %s", session.Username, toJSON(product))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// overlappingDates returns true if the SalesStart..SalesEnd ranges of the two
// SKUs overlap.
func overlappingDates(sku1, sku2 *model.SKU) bool {
	if sku1.SalesStart.IsZero() {
		if sku1.SalesEnd.IsZero() {
			return true
		}
		return sku2.SalesStart.Before(sku1.SalesEnd)
	}
	if sku1.SalesEnd.IsZero() {
		return sku2.SalesEnd.IsZero() || sku2.SalesEnd.After(sku1.SalesStart)
	}
	if sku2.SalesStart.IsZero() {
		return sku2.SalesEnd.IsZero() || sku2.SalesEnd.After(sku1.SalesStart)
	}
	if sku2.SalesEnd.IsZero() {
		return sku2.SalesStart.Before(sku1.SalesEnd)
	}
	if sku1.SalesStart.Before(sku2.SalesStart) {
		return sku1.SalesEnd.After(sku2.SalesStart)
	}
	return sku2.SalesEnd.After(sku1.SalesStart)
}
