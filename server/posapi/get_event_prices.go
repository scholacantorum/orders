package posapi

import (
	"github.com/mailru/easyjson/jwriter"

	"scholacantorum.org/orders/api"
	"scholacantorum.org/orders/auth"
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
func GetEventPrices(r *api.Request, eventID model.EventID) error {
	var (
		event      *model.Event
		productIDs []string
		pdata      []*getPricesData
		jw         jwriter.Writer
	)
	if r.Privileges&model.PrivInPersonSales == 0 {
		return auth.Forbidden
	}
	// Get the list of products offering ticket sales for that event.
	if event = r.Tx.FetchEvent(eventID); event == nil {
		return api.NotFound
	}
	productIDs = getEventProducts(r, event)
	// Look up the prices for each product.
	for _, pid := range productIDs {
		var (
			product *model.Product
			sku     *model.SKU
			pd      getPricesData
		)
		// Get the product.  Make sure it has capacity.
		product = r.Tx.FetchProduct(model.ProductID(pid))
		if !api.ProductHasCapacity(r, product) {
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
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if len(pdata) == 0 {
		// No products available for sale, or even with messages to
		// display, so we return a null.
		jw.RawString("null")
	} else {
		// Return the product data.
		emitGetPrices(&jw, pdata)
	}
	jw.DumpTo(r)
	return nil
}

// getEventProducts returns the product names for the products selling tickets
// to the specified event.  Only products that are targeted at the event, or not
// targeted at any specific event, are included.
func getEventProducts(r *api.Request, event *model.Event) (productIDs []string) {
	for _, p := range r.Tx.FetchProductsByEvent(event) {
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
func emitGetPrices(jw *jwriter.Writer, pdata []*getPricesData) {
	jw.RawByte('[')
	for i, pd := range pdata {
		if i != 0 {
			jw.RawByte(',')
		}
		jw.RawString(`{"id":`)
		jw.String(string(pd.id))
		jw.RawString(`,"name":`)
		jw.String(pd.name)
		jw.RawString(`,"price":`)
		jw.Int(pd.price)
		jw.RawString(`,"ticketCount":`)
		jw.Int(pd.ticketCount)
		jw.RawByte('}')
	}
	jw.RawByte(']')
}
