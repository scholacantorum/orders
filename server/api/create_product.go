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
	)
	if session = auth.GetSession(tx, w, r, model.PrivSetup); session == nil {
		return
	}
	if product, err = parseCreateProduct(r.Body); err != nil {
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
		if sku.Quantity == 0 {
			sku.Quantity = 1
		}
		if sku.Price < 0 || sku.Quantity < 1 ||
			(!sku.SalesStart.IsZero() && !sku.SalesEnd.IsZero() && !sku.SalesEnd.After(sku.SalesStart)) {
			BadRequestError(tx, w, "invalid SKU parameters")
			return
		}
		for j := 0; j < i; j++ {
			prev := product.SKUs[j]
			if prev.MembersOnly == sku.MembersOnly && prev.Coupon == sku.Coupon && prev.Quantity == sku.Quantity &&
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
		case "name":
			return json.StringHandler(func(s string) { p.Name = s })
		case "type":
			return json.StringHandler(func(s string) { p.Type = model.ProductType(s) })
		case "receipt":
			return json.StringHandler(func(s string) { p.Receipt = s })
		case "ticketName":
			return json.StringHandler(func(s string) { p.TicketName = s })
		case "ticketCount":
			return json.IntHandler(func(i int) { p.TicketCount = i })
		case "ticketClass":
			return json.StringHandler(func(s string) { p.TicketClass = s })
		case "skus":
			return json.ArrayHandler(func() json.Handlers {
				var sku = model.SKU{Quantity: 1}
				p.SKUs = append(p.SKUs, &sku)
				return parseCreateProductSKU(&sku)
			})
		case "events":
			return json.ArrayHandler(func() json.Handlers {
				return json.Handlers{String: func(s string) {
					p.Events = append(p.Events, &model.Event{ID: model.EventID(s)})
				}}
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
		case "coupon":
			return json.StringHandler(func(s string) { sku.Coupon = s })
		case "salesStart":
			return json.TimeHandler(func(t time.Time) { sku.SalesStart = t })
		case "salesEnd":
			return json.TimeHandler(func(t time.Time) { sku.SalesEnd = t })
		case "membersOnly":
			return json.BoolHandler(func(b bool) { sku.MembersOnly = b })
		case "quantity":
			return json.IntHandler(func(i int) { sku.Quantity = i })
		case "price":
			return json.IntHandler(func(i int) { sku.Price = i })
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
		jw.Prop("name", p.Name)
		jw.Prop("type", string(p.Type))
		jw.Prop("receipt", p.Receipt)
		if p.TicketName != "" {
			jw.Prop("ticketName", p.TicketName)
		}
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
						if sku.Coupon != "" {
							jw.Prop("coupon", sku.Coupon)
						}
						if !sku.SalesStart.IsZero() {
							jw.Prop("salesStart", sku.SalesStart.Format(time.RFC3339))
						}
						if !sku.SalesEnd.IsZero() {
							jw.Prop("salesEnd", sku.SalesEnd.Format(time.RFC3339))
						}
						if sku.MembersOnly {
							jw.Prop("membersOnly", true)
						}
						if sku.Quantity != 1 {
							jw.Prop("quantity", sku.Quantity)
						}
						jw.Prop("price", sku.Price)
					})
				}
			})
		})
		if len(p.Events) != 0 {
			jw.Prop("events", func() {
				jw.Array(func() {
					for _, e := range p.Events {
						jw.String(string(e.ID))
					}
				})
			})
		}
	})
	jw.Close()
	return buf.Bytes()
}
