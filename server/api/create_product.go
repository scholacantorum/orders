package api

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/rothskeller/json"

	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
)

// CreateProduct handles POST /api/product requests.
func CreateProduct(tx db.Tx, w http.ResponseWriter, r *http.Request) {
	var (
		session   *model.Session
		product   *model.Product
		out       []byte
		err       error
		seenEvent = map[model.EventID]bool{}
		seenPrio0 bool
	)
	if session = auth.GetSession(tx, w, r, model.PrivSetupOrders); session == nil {
		return
	}
	if product, err = parseCreateProduct(r.Body); err != nil {
		BadRequestError(tx, w, err.Error())
		return
	}
	if product.ID == "" || product.Name == "" || product.ShortName == "" || product.Type == "" || product.TicketCount < 0 {
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
	for _, pe := range product.Events {
		if pe.Event == nil {
			BadRequestError(tx, w, "invalid event")
			return
		}
		if seenEvent[pe.Event.ID] {
			BadRequestError(tx, w, "duplicate event")
			return
		}
		if tx.FetchEvent(pe.Event.ID) == nil {
			BadRequestError(tx, w, "nonexistent event")
			return
		}
		seenEvent[pe.Event.ID] = true
		if pe.Priority == 0 {
			if seenPrio0 {
				BadRequestError(tx, w, "multiple priority 0 events")
				return
			}
			seenPrio0 = true
		}
	}
	for i, sku := range product.SKUs {
		switch sku.Source {
		case model.OrderFromPublic, model.OrderFromMembers, model.OrderFromGala, model.OrderFromOffice, model.OrderInPerson:
		default:
			BadRequestError(tx, w, "invalid SKU source")
			return
		}
		if sku.Price < 0 || (!sku.SalesStart.IsZero() && !sku.SalesEnd.IsZero() && !sku.SalesEnd.After(sku.SalesStart)) {
			BadRequestError(tx, w, "invalid SKU parameters")
			return
		}
		for j := 0; j < i; j++ {
			prev := product.SKUs[j]
			if prev.Source == sku.Source && prev.Flags == sku.Flags && prev.Coupon == sku.Coupon &&
				overlappingDates(prev, sku) {
				BadRequestError(tx, w, "overlapping SKUs")
				return
			}
		}
	}
	tx.SaveProduct(product)
	commit(tx)
	out = emitProduct(product)
	log.Printf("%s CREATE PRODUCT %s", session.Username, out)
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// parseCreateProduct reads the product details from the request body.
func parseCreateProduct(r io.Reader) (p *model.Product, err error) {
	var jr = json.NewReader(r)

	p = new(model.Product)
	err = jr.Read(json.ObjectHandler(func(key string) json.Handlers {
		switch key {
		case "id":
			return json.StringHandler(func(s string) { p.ID = model.ProductID(s) })
		case "series":
			return json.StringHandler(func(s string) { p.Series = s })
		case "name":
			return json.StringHandler(func(s string) { p.Name = s })
		case "shortname":
			return json.StringHandler(func(s string) { p.ShortName = s })
		case "type":
			return json.StringHandler(func(s string) { p.Type = model.ProductType(s) })
		case "receipt":
			return json.StringHandler(func(s string) { p.Receipt = s })
		case "ticketCount":
			return json.IntHandler(func(i int) { p.TicketCount = i })
		case "ticketClass":
			return json.StringHandler(func(s string) { p.TicketClass = s })
		case "skus":
			return json.ArrayHandler(func() json.Handlers {
				var sku model.SKU
				p.SKUs = append(p.SKUs, &sku)
				return parseCreateProductSKU(&sku)
			})
		case "events":
			return json.ArrayHandler(func() json.Handlers {
				p.Events = append(p.Events, model.ProductEvent{})
				return parseCreateProductEvent(&p.Events[len(p.Events)-1])
			})
		default:
			return json.RejectHandler()
		}
	}))
	return p, err
}
func parseCreateProductSKU(sku *model.SKU) json.Handlers {
	return json.ObjectHandler(func(key string) json.Handlers {
		switch key {
		case "source":
			return json.StringHandler(func(s string) { sku.Source = model.OrderSource(s) })
		case "coupon":
			return json.StringHandler(func(s string) { sku.Coupon = s })
		case "salesStart":
			return json.TimeHandler(func(t time.Time) { sku.SalesStart = t })
		case "salesEnd":
			return json.TimeHandler(func(t time.Time) { sku.SalesEnd = t })
		case "price":
			return json.IntHandler(func(i int) { sku.Price = i })
		default:
			return json.RejectHandler()
		}
	})
}
func parseCreateProductEvent(pe *model.ProductEvent) json.Handlers {
	return json.ObjectHandler(func(key string) json.Handlers {
		switch key {
		case "event":
			return json.StringHandler(func(s string) { pe.Event = &model.Event{ID: model.EventID(s)} })
		case "priority":
			return json.IntHandler(func(i int) { pe.Priority = i })
		default:
			return json.RejectHandler()
		}
	})
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

// emitProduct generates the JSON representation of a product.
func emitProduct(p *model.Product) []byte {
	var (
		buf bytes.Buffer
		jw  = json.NewWriter(&buf)
	)
	jw.Object(func() {
		jw.Prop("id", string(p.ID))
		jw.Prop("series", p.Series)
		jw.Prop("name", p.Name)
		jw.Prop("shortname", p.ShortName)
		jw.Prop("type", string(p.Type))
		jw.Prop("receipt", p.Receipt)
		if p.TicketCount != 0 {
			jw.Prop("ticketCount", p.TicketCount)
		}
		if p.TicketClass != "" {
			jw.Prop("ticketClass", p.TicketClass)
		}
		jw.Prop("skus", func() {
			jw.Array(func() {
				for _, sku := range p.SKUs {
					jw.Object(func() {
						jw.Prop("source", sku.Source)
						if sku.Coupon != "" {
							jw.Prop("coupon", sku.Coupon)
						}
						if !sku.SalesStart.IsZero() {
							jw.Prop("salesStart", sku.SalesStart.Format(time.RFC3339))
						}
						if !sku.SalesEnd.IsZero() {
							jw.Prop("salesEnd", sku.SalesEnd.Format(time.RFC3339))
						}
						jw.Prop("price", sku.Price)
					})
				}
			})
		})
		if len(p.Events) != 0 {
			jw.Prop("events", func() {
				jw.Array(func() {
					for _, pe := range p.Events {
						jw.Object(func() {
							jw.Prop("event", string(pe.Event.ID))
							jw.Prop("priority", pe.Priority)
						})
					}
				})
			})
		}
	})
	jw.Close()
	return buf.Bytes()
}
