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

// CreateSKU handles POST /api/product/$id/sku requests.
func CreateSKU(tx *sql.Tx, w http.ResponseWriter, r *http.Request, productID model.ProductID) {
	var (
		session *model.Session
		product *model.Product
		sku     *model.SKU
		price   int
		coupon  string
		err     error
	)
	if session = auth.GetSession(tx, w, r, model.PrivSetup); session == nil {
		return
	}
	if product = model.FetchProduct(tx, productID); product == nil {
		NotFoundError(tx, w)
		return
	}
	if err = json.NewDecoder(r.Body).Decode(&sku); err != nil {
		BadRequestError(tx, w, err.Error())
		return
	}
	if sku.ID != 0 || sku.StripeID == "" || sku.Price < 0 || (sku.Product != 0 && sku.Product != productID) ||
		(!sku.SalesStart.IsZero() && !sku.SalesEnd.IsZero() && !sku.SalesEnd.After(sku.SalesStart)) {
		BadRequestError(tx, w, "invalid parameters")
		return
	}
	sku.Product = productID
	for _, existing := range model.FetchSKUsForProduct(tx, productID) {
		if existing.MembersOnly == sku.MembersOnly && existing.Coupon == sku.Coupon && overlappingDates(existing, sku) {
			BadRequestError(tx, w, "conflicts with existing SKU")
			return
		}
	}
	if price, coupon = stripe.ValidateSKU(product.StripeID, sku.StripeID); price < 0 {
		BadRequestError(tx, w, "invalid Stripe SKU ID")
		return
	}
	if price == 100 {
		// We use $1 in Stripe for variable-price SKUs, but we use 0 for
		// that in our database.
		price = 0
	}
	if (sku.Price != 0 && sku.Price != price) || sku.Coupon != coupon {
		BadRequestError(tx, w, "price or coupon mismatch with Stripe")
	}
	sku.Price = price
	sku.Save(tx)
	commit(tx)
	log.Printf("%s CREATE SKU %s", session.Username, toJSON(sku))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sku)
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
