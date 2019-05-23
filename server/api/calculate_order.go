package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
)

// CalculateOrder handles GET /api/order requests.  It takes a JSON order in its
// request body, and returns the same order with prices filled in.  The order is
// not validated except as necessary to compute prices.
func CalculateOrder(tx db.Tx, w http.ResponseWriter, r *http.Request) {
	var (
		session *model.Session
		order   *model.Order
		privs   model.Privilege
		err     error
	)
	// Get current session privileges.
	if auth.HasSession(r) {
		if session = auth.GetSession(tx, w, r, 0); session == nil {
			return
		}
		privs = session.Privileges
	}
	// Read the order details from the request.
	if err = json.NewDecoder(r.Body).Decode(&order); err != nil {
		BadRequestError(tx, w, err.Error())
		return
	}
	// Resolve the SKUs and fill in the prices.
	if !resolveSKUs(tx, order, privs, false) {
		BadRequestError(tx, w, "invalid products or prices")
		return
	}
	// Send the result back.
	commit(tx)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// resolveSKUs walks through each line of the order, finding the listed product
// and selecting the best SKU of that product, following the rules documented in
// db/schema.sql.  If mustMatch is true, it requires that the line have the
// correct price for that SKU.  If mustMatch is false, it overrides the line
// price, setting it to the price for that SKU. It returns true if everything
// resolved successfully and false otherwise.
func resolveSKUs(tx db.Tx, order *model.Order, privs model.Privilege, mustMatch bool) bool {
	for _, line := range order.Lines {
		var sku *model.SKU

		if line.Product = tx.FetchProduct(line.Product.ID); line.Product == nil {
			return false
		}
		if line.Product.Type == model.ProdAuctionItem || line.Product.Type == model.ProdDonation {
			if line.Price < 1 {
				return false
			}
			continue
		}
		for _, s := range line.Product.SKUs {
			if !matchSKU(order, privs, s) {
				continue
			}
			if sku == nil || betterSKU(s, sku) {
				sku = s
			}
		}
		if sku == nil {
			return false
		}
		if mustMatch && line.Price != sku.Price {
			return false
		}
		line.Price = sku.Price
	}
	return true
}

// matchSKU returns true if all of the criteria for the SKU are met by the order
// being placed.
func matchSKU(order *model.Order, privs model.Privilege, sku *model.SKU) bool {
	if sku.MembersOnly && privs == 0 {
		return false
	}
	if sku.Coupon != "" && !strings.EqualFold(sku.Coupon, order.Coupon) {
		return false
	}
	if privs&model.PrivHandleOrders != 0 {
		return true
	}
	if !sku.SalesStart.IsZero() && sku.SalesStart.After(time.Now()) {
		return false
	}
	if !sku.SalesEnd.IsZero() && sku.SalesEnd.Before(time.Now()) {
		return false
	}
	return true
}

// betterSKU returns whether sku1 is a higher priority match than sku2.
func betterSKU(sku1, sku2 *model.SKU) bool {
	if sku1.MembersOnly && !sku2.MembersOnly {
		return true
	}
	if !sku1.MembersOnly && sku2.MembersOnly {
		return false
	}
	if sku1.Coupon != "" && sku2.Coupon == "" {
		return true
	}
	if sku1.Coupon == "" && sku2.Coupon != "" {
		return false
	}
	if (sku1.SalesStart.IsZero() || !sku1.SalesStart.After(time.Now())) &&
		(sku2.SalesEnd.IsZero() || !sku2.SalesEnd.Before(time.Now())) {
		return true
	}
	return false
}
