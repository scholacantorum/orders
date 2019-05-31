package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rothskeller/json"

	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
)

type getPricesData struct {
	id      model.ProductID
	name    string
	message string
	price   int
}

// GetPrices returns the prices and availability of one or more products.  It is
// used to drive the order form.
func GetPrices(tx db.Tx, w http.ResponseWriter, r *http.Request) {
	var (
		session     *model.Session
		privs       model.Privilege
		coupon      string
		couponMatch bool
		pdata       []*getPricesData
		product     *model.Product
		masterSKU   *model.SKU
		jw          json.Writer
	)
	// Get current session privileges.
	if auth.HasSession(r) {
		if session = auth.GetSession(tx, w, r, 0); session == nil {
			return
		}
		privs = session.Privileges
	}
	// Read the request details from the request.
	if coupon = r.FormValue("coupon"); coupon == "" {
		couponMatch = true
	}
	// Look up the prices for each product.
	for _, pid := range r.Form["p"] {
		var (
			sku *model.SKU
			pd  getPricesData
		)
		// Get the product.  Skip nonexistent ones.
		if product = tx.FetchProduct(model.ProductID(pid)); product == nil {
			continue
		}
		pd.id = product.ID
		pd.name = product.ShortName
		// Find the best SKU for this product.
		for _, s := range product.SKUs {
			if !interestingSKU(product, s, privs, coupon) {
				continue
			}
			if s.Coupon == coupon {
				couponMatch = true
			}
			sku = betterSKU2(s, sku)
		}
		if sku == nil {
			// No applicable SKUs; ignore product.
			continue
		}
		// Generate the product data to return.
		if !hasCapacity(tx, product) {
			pd.message = "This event is sold out."
		} else if pd.message = noSalesMessage(sku, privs); pd.message == "" {
			pd.price = sku.Price
		}
		pdata = append(pdata, &pd)
		// The masterSKU determines the message to be shown in lieu of
		// the purchase button if none of the products are available for
		// sale.
		masterSKU = betterSKU2(sku, masterSKU)
	}
	commit(tx)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jw = json.NewWriter(w)
	if len(pdata) == 0 {
		// No products available for sale, or even with messages to
		// display, so we return a null.
		jw.Null()
	} else if message := noSalesMessage(masterSKU, privs); message != "" {
		// No products available for sale, but we do have a message to
		// display.
		jw.String(message)
	} else {
		// Return the product data.
		emitGetPrices(jw, couponMatch, pdata)
	}
	jw.Close()
}

// interestingSKU returns true if the SKU's criteria are met, or come close.
// "Close", in this context, means that all criteria are met except for the
// sales range, and one of the following is true:
//   - The caller has Sell privilege
//   - The sales range is in the future
//   - The product is a ticket to an event that starts today or later
func interestingSKU(product *model.Product, sku *model.SKU, privs model.Privilege, coupon string) bool {
	if sku.MembersOnly && privs == 0 {
		return false
	}
	if sku.Coupon != "" && sku.Coupon != coupon {
		return false
	}
	if privs&model.PrivSell != 0 {
		return true
	}
	var now = time.Now()
	if sku.InSalesRange(now) <= 0 {
		return true
	}
	var today = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	for _, pe := range product.Events {
		if pe.Priority != 0 {
			continue
		}
		if !pe.Event.Start.Before(today) {
			return true
		}
	}
	return false
}

// betterSKU returns whichever SKU has the better fit for the purchase.  It
// assumes that all SKUs passed the interestingSKU check.  If both SKUs are in
// sales range, the cheaper one is selected.  Otherwise, a current range is
// better than one in the future, which is better than one in the past; and for
// both future and past, a range closer to today is better than one further
// away.
func betterSKU2(sku1, sku2 *model.SKU) *model.SKU {
	var isr1, isr2 int
	var now = time.Now()

	if sku2 == nil {
		return sku1
	}
	if sku1 == nil {
		return sku2
	}
	isr1 = sku1.InSalesRange(now)
	isr2 = sku2.InSalesRange(now)
	if isr1 == 0 && isr2 == 0 {
		if sku1.Price < sku2.Price {
			return sku1
		}
		return sku2
	}
	if isr1 == 0 {
		return sku1
	}
	if isr2 == 0 {
		return sku2
	}
	if isr1 < 0 && isr2 > 0 {
		return sku1
	}
	if isr2 < 0 && isr1 > 0 {
		return sku2
	}
	if isr1 < 0 {
		if sku1.SalesStart.Before(sku2.SalesStart) {
			return sku1
		}
		return sku2
	}
	if sku1.SalesEnd.After(sku2.SalesEnd) {
		return sku1
	}
	return sku2
}

// noSalesMessage returns the string describing why the SKU isn't available for
// sale, or an empty string if it is available.
func noSalesMessage(sku *model.SKU, privs model.Privilege) string {
	if privs&model.PrivSell != 0 {
		return ""
	}
	switch sku.InSalesRange(time.Now()) {
	case -1:
		return fmt.Sprintf("Sales start on %s.", sku.SalesStart.Format("January\u00A02"))
	case 0:
		return ""
	case +1:
		return "Tickets available at the door."
	}
	panic("not reachable")
}

// hasCapacity returns false if the product is a ticket for an event that is
// sold out.  It returns true for non-ticket products.
func hasCapacity(tx db.Tx, product *model.Product) bool {
	for _, pe := range product.Events {
		if pe.Priority == 0 {
			if pe.Event.Capacity == 0 {
				return true
			}
			return tx.FetchTicketCount(pe.Event) < pe.Event.Capacity
		}
	}
	return true
}

// emitGetPrices writes the JSON response.
func emitGetPrices(jw json.Writer, couponMatch bool, pdata []*getPricesData) {
	jw.Object(func() {
		jw.Prop("coupon", couponMatch)
		jw.Prop("products", func() {
			jw.Array(func() {
				for _, pd := range pdata {
					jw.Object(func() {
						jw.Prop("id", string(pd.id))
						jw.Prop("name", pd.name)
						if pd.message != "" {
							jw.Prop("message", pd.message)
						} else {
							jw.Prop("price", pd.price)
						}
					})
				}
			})
		})
	})
}
