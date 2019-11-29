package posapi

import (
	"log"
	"net/http"

	"github.com/rothskeller/json"

	"scholacantorum.org/orders/api"
	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
)

type getPricesData struct {
	id          model.ProductID
	name        string
	price       int
	ticketCount int
}

// GetEventPrices returns the prices and availability of the products giving
// admission to an event.
func GetEventPrices(tx db.Tx, w http.ResponseWriter, r *http.Request, eventID model.EventID) {
	var (
		event      *model.Event
		productIDs []string
		pdata      []*getPricesData
		jw         json.Writer
	)
	// Get current session privileges.
	if auth.GetSession(tx, w, r, model.PrivInPersonSales) == nil {
		return
	}
	// Get the list of products offering ticket sales for that event.
	if event = tx.FetchEvent(eventID); event == nil {
		log.Printf("ERROR: no such event %q", eventID)
		api.NotFoundError(tx, w)
		return
	}
	productIDs = getEventProducts(tx, event)
	// Look up the prices for each product.
	for _, pid := range productIDs {
		var (
			product *model.Product
			sku     *model.SKU
			pd      getPricesData
		)
		// Get the product.  Make sure it has capacity.
		product = tx.FetchProduct(model.ProductID(pid))
		if !api.ProductHasCapacity(tx, product) {
			continue
		}
		pd.id = product.ID
		pd.name = product.ShortName
		// Find the best SKU for this product.
		for _, s := range product.SKUs {
			if !api.MatchingSKU(s, "", model.OrderInPerson, false) {
				continue
			}
			sku = api.BetterSKU(sku, s)
		}
		if sku == nil {
			// No applicable SKUs; ignore product.
			continue
		}
		// Generate the product data to return.
		pd.price = sku.Price
		pd.ticketCount = product.TicketCount
		pdata = append(pdata, &pd)
	}
	api.Commit(tx)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jw = json.NewWriter(w)
	if len(pdata) == 0 {
		// No products available for sale, or even with messages to
		// display, so we return a null.
		jw.Null()
	} else {
		// Return the product data.
		emitGetPrices(jw, pdata)
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

// emitGetPrices writes the JSON response.
func emitGetPrices(jw json.Writer, pdata []*getPricesData) {
	jw.Object(func() {
		jw.Prop("coupon", true) // needed by iOS app 2019-10-01
		jw.Prop("products", func() {
			jw.Array(func() {
				for _, pd := range pdata {
					jw.Object(func() {
						jw.Prop("id", string(pd.id))
						jw.Prop("name", pd.name)
						jw.Prop("price", pd.price)
						jw.Prop("ticketCount", pd.ticketCount)
					})
				}
			})
		})
	})
}
