package ofcapi

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mailru/easyjson"

	"scholacantorum.org/orders/api"
	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/model"
)

// CreateProduct handles POST /ofcapi/product requests.
func CreateProduct(r *api.Request) error {
	var (
		product   *model.Product
		err       error
		seenEvent = map[model.EventID]bool{}
		seenPrio0 bool
	)
	if r.Privileges&model.PrivSetupOrders == 0 {
		return auth.Forbidden
	}
	if product, err = parseCreateProduct(r); err != nil {
		return err
	}
	if r.Tx.FetchProduct(product.ID) != nil {
		return errors.New("duplicate product ID")
	}
	if product.TicketCount > 0 {
		if len(product.Events) == 0 {
			return errors.New("ticket products must have associated events")
		}
	} else {
		if len(product.Events) != 0 {
			return errors.New("only ticket products can have associated events")
		}
	}
	for _, pe := range product.Events {
		if seenEvent[pe.Event.ID] {
			return errors.New("duplicate event")
		}
		if r.Tx.FetchEvent(pe.Event.ID) == nil {
			return errors.New("nonexistent event")
		}
		seenEvent[pe.Event.ID] = true
		if pe.Priority == 0 {
			if seenPrio0 {
				return errors.New("multiple priority 0 events")
			}
			seenPrio0 = true
		}
	}
	for i, sku := range product.SKUs {
		if !sku.SalesStart.IsZero() && !sku.SalesEnd.IsZero() && !sku.SalesEnd.After(sku.SalesStart) {
			return errors.New("invalid SKU sales range")
		}
		for j := 0; j < i; j++ {
			prev := product.SKUs[j]
			if prev.Source == sku.Source && prev.Coupon == sku.Coupon && overlappingDates(prev, sku) {
				return errors.New("overlapping SKUs")
			}
		}
	}
	r.Tx.SaveProduct(product)
	r.Tx.Commit()
	easyjson.MarshalToHTTPResponseWriter(product, r)
	return nil
}

// parseCreateProduct reads the product details from the request body.
func parseCreateProduct(r *api.Request) (p *model.Product, err error) {
	p = new(model.Product)
	if p.ID = model.ProductID(strings.TrimSpace(r.FormValue("id"))); p.ID == "" {
		return nil, errors.New("missing id")
	}
	p.Series = strings.TrimSpace(r.FormValue("series"))
	if p.Name = strings.TrimSpace(r.FormValue("name")); p.Name == "" {
		return nil, errors.New("missing name")
	}
	if p.ShortName = strings.TrimSpace(r.FormValue("shortname")); p.ShortName == "" {
		return nil, errors.New("missing shortname")
	}
	switch p.Type = model.ProductType(strings.TrimSpace(r.FormValue("type"))); p.Type {
	case model.ProdAuctionItem, model.ProdDonation, model.ProdOther, model.ProdRecording, model.ProdSheetMusic,
		model.ProdTicket, model.ProdWardrobe:
		break
	case "":
		return nil, errors.New("missing type")
	default:
		return nil, errors.New("invalid type")
	}
	p.Receipt = strings.TrimSpace(r.FormValue("receipt"))
	if tcs := r.FormValue("ticketCount"); tcs != "" {
		if tc, err := strconv.Atoi(r.FormValue("ticketCount")); err == nil && tc >= 0 {
			p.TicketCount = tc
		} else {
			return nil, errors.New("invalid ticketCount")
		}
	}
	p.TicketClass = strings.TrimSpace(r.FormValue("ticketClass"))
	for i := 0; true; i++ {
		var (
			sku    model.SKU
			prefix = fmt.Sprintf("sku%d.", i)
		)
		if sku.Source = model.OrderSource(r.FormValue(prefix + "source")); sku.Source == "" {
			break
		}
		switch sku.Source {
		case model.OrderFromPublic, model.OrderFromMembers, model.OrderFromGala, model.OrderFromOffice, model.OrderInPerson:
			break
		default:
			return nil, errors.New("invalid SKU source")
		}
		sku.Coupon = strings.ToUpper(strings.TrimSpace(r.FormValue(prefix + "coupon")))
		if sss := r.FormValue(prefix + "salesStart"); sss != "" {
			if ss, err := time.Parse(time.RFC3339, sss); err == nil {
				sku.SalesStart = ss.In(time.Local)
			} else {
				return nil, errors.New("invalid SKU salesStart")
			}
		}
		if ses := r.FormValue(prefix + "salesEnd"); ses != "" {
			if se, err := time.Parse(time.RFC3339, ses); err == nil {
				sku.SalesEnd = se.In(time.Local)
			} else {
				return nil, errors.New("invalid SKU salesEnd")
			}
		}
		if pr, err := strconv.Atoi(r.FormValue(prefix + "price")); err == nil && pr >= 0 {
			sku.Price = pr
		} else {
			return nil, errors.New("invalid SKU price")
		}
		p.SKUs = append(p.SKUs, &sku)
	}
	for i := 0; true; i++ {
		var (
			pe     model.ProductEvent
			prefix = fmt.Sprintf("event%d.", i)
		)
		if eid := model.EventID(strings.TrimSpace(r.FormValue(prefix + "id"))); eid != "" {
			pe.Event = &model.Event{ID: eid}
		} else {
			break
		}
		if pr, err := strconv.Atoi(r.FormValue(prefix + "priority")); err == nil {
			pe.Priority = pr
		} else {
			return nil, errors.New("invalid event priority")
		}
		p.Events = append(p.Events, pe)
	}
	return p, nil
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
