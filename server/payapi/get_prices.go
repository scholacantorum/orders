package payapi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rothskeller/json"

	"scholacantorum.org/orders/api"
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

// GetPrices returns the prices and availability of one or more products.  It is
// used to drive payment forms.
func GetPrices(tx db.Tx, w http.ResponseWriter, r *http.Request) {
	var (
		source      model.OrderSource
		session     *model.Session
		productIDs  []string
		coupon      string
		couponMatch bool
		pdata       []*getPricesData
		product     *model.Product
		masterSKU   *model.SKU
		jw          json.Writer
	)
	// Get the request source and authorization.
	switch source = model.OrderSource(r.FormValue("source")); source {
	case "":
		source = model.OrderFromPublic
	case model.OrderFromPublic:
		// no-op
	case model.OrderFromMembers:
		if session = auth.GetSessionMembersAuth(tx, w, r, r.FormValue("auth")); session == nil {
			return
		}
	default:
		api.BadRequestError(tx, w, "invalid source")
		return
	}
	// Read the request details from the request.
	productIDs = r.Form["p"]
	if coupon = r.FormValue("coupon"); coupon == "" {
		couponMatch = true
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
			if !api.MatchingSKU(s, coupon, source, true) {
				continue
			}
			if s.Coupon == coupon {
				couponMatch = true
			}
			sku = api.BetterSKU(s, sku)
		}
		if sku == nil {
			// No applicable SKUs; ignore product.
			continue
		}
		// Generate the product data to return.
		if !api.ProductHasCapacity(tx, product) {
			pd.message = "This event is sold out."
		} else if pd.message = noSalesMessage(sku); pd.message == "" {
			pd.price = sku.Price
		}
		pdata = append(pdata, &pd)
		// The masterSKU determines the message to be shown in lieu of
		// the purchase button if none of the products are available for
		// sale.
		masterSKU = api.BetterSKU(sku, masterSKU)
	}
	api.Commit(tx)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jw = json.NewWriter(w)
	if len(pdata) == 0 {
		// No products available for sale, or even with messages to
		// display, so we return a null.
		jw.Null()
	} else if message := noSalesMessage(masterSKU); message != "" {
		// No products available for sale, but we do have a message to
		// display.
		jw.String(message)
	} else {
		// Return the product data.
		emitGetPrices(jw, session, couponMatch, pdata)
	}
	jw.Close()
}

// noSalesMessage returns the string describing why the SKU isn't available for
// sale, or an empty string if it is available.
func noSalesMessage(sku *model.SKU) string {
	if sku.InSalesRange(time.Now()) < 0 {
		return fmt.Sprintf("Sales start on %s.", sku.SalesStart.Format("January\u00A02"))
	}
	return ""
}

// emitGetPrices writes the JSON response.
func emitGetPrices(jw json.Writer, session *model.Session, couponMatch bool, pdata []*getPricesData) {
	jw.Object(func() {
		if session != nil && session.Name != "" {
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
