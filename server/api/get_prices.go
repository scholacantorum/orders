package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rothskeller/json"

	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/config"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
)

type getPricesData struct {
	id      model.ProductID
	name    string
	message string
	price   int
}

// GetPrices returns the prices and availability of one or more products, or the
// products giving admission to an event.  It is used to drive the order form
// and the in-person sales app.
func GetPrices(tx db.Tx, w http.ResponseWriter, r *http.Request) {
	var (
		eventID     model.EventID
		session     *model.Session
		privs       model.Privilege
		event       *model.Event
		productIDs  []string
		coupon      string
		couponMatch bool
		pdata       []*getPricesData
		product     *model.Product
		masterSKU   *model.SKU
		jw          json.Writer
	)
	// If the caller specified an event, they have to have in-person sales
	// privilege.
	if eventID = model.EventID(r.FormValue("event")); eventID != "" {
		privs = model.PrivInPersonSales
	}
	// Get current session privileges.
	if auth.HasSession(r) {
		if session = auth.GetSession(tx, w, r, privs); session == nil {
			return
		}
		privs = session.Privileges
	} else if token := r.FormValue("auth"); token != "" {
		if session = auth.GetSessionMembersAuth(tx, w, r, token); session == nil {
			return
		}
		privs = session.Privileges
	}
	// If the caller specified an event, get the list of products offering
	// ticket sales for that event.
	if eventID != "" {
		if event = tx.FetchEvent(eventID); event == nil {
			NotFoundError(tx, w)
			return
		}
		productIDs = getEventProducts(tx, event)
	} else {
		// Read the request details from the request.
		productIDs = r.Form["p"]
		if coupon = r.FormValue("coupon"); coupon == "" {
			couponMatch = true
		}
	}
	// Look up the prices for each product.
	for _, pid := range productIDs {
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
			if !interestingSKU(product, s, privs, coupon, event != nil) {
				continue
			}
			if s.Coupon == coupon {
				couponMatch = true
			}
			sku = betterSKU(s, sku)
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
		masterSKU = betterSKU(sku, masterSKU)
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
		emitGetPrices(jw, session, couponMatch, pdata)
	}
	jw.Close()
}

// getEventProducts returns the product names for the products selling tickets
// to the specified event.  Only products that are targeted at the event, or not
// targeted at any specific event, are included.
func getEventProducts(tx db.Tx, event *model.Event) (productIDs []string) {
	for _, p := range tx.FetchProductsByEvent(event) {
		var targeted, match bool
		for _, pe := range p.Events {
			if pe.Priority == 0 {
				targeted = true
				if pe.Event.ID == event.ID {
					match = true
				}
			}
		}
		if match || !targeted {
			productIDs = append(productIDs, string(p.ID))
		}
	}
	return productIDs
}

// interestingSKU returns true if the SKU's criteria are met, or come close.
// "Close", in this context, means that all criteria are met except for the
// sales range, and one of the following is true:
//   - The caller has PrivInPersonSales
//   - The sales range is in the future
//   - The product is a ticket to an event that starts today or later
func interestingSKU(product *model.Product, sku *model.SKU, privs model.Privilege, coupon string, inPerson bool) bool {
	if sku.Flags&model.SKUMembersOnly != 0 && privs == 0 {
		return false
	}
	if sku.Flags&model.SKUInPerson != 0 && !inPerson {
		return false
	}
	if sku.Coupon != "" && sku.Coupon != coupon {
		return false
	}
	if privs&model.PrivInPersonSales != 0 {
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

// noSalesMessage returns the string describing why the SKU isn't available for
// sale, or an empty string if it is available.
func noSalesMessage(sku *model.SKU, privs model.Privilege) string {
	if privs&model.PrivInPersonSales != 0 {
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
func emitGetPrices(jw json.Writer, session *model.Session, couponMatch bool, pdata []*getPricesData) {
	jw.Object(func() {
		if session.Name != "" {
			// If there's a name in the session, it came from
			// GetSessionMembersAuth, which means it came from the
			// recordings order form loaded in an iframe of the
			// members site, so we'll provide full user detail and
			// also the stripe key.
			jw.Prop("user", func() {
				jw.Object(func() {
					jw.Prop("id", int(session.Member))
					jw.Prop("username", session.Username)
					jw.Prop("name", session.Name)
					jw.Prop("email", session.Email)
					jw.Prop("address", session.Address)
					jw.Prop("city", session.City)
					jw.Prop("state", session.State)
					jw.Prop("zip", session.Zip)
				})
			})
			jw.Prop("stripePublicKey", config.Get("stripePublicKey"))
		}
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
