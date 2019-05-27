package api

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/rothskeller/json"

	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
)

// GetTicket is the API used by the scanner app to get information about the
// ticket(s) on the specified order to the specified event.
func GetTicket(tx db.Tx, w http.ResponseWriter, r *http.Request, eventID model.EventID, token string) {
	type class struct {
		val  int
		max  int
		used int
	}
	var (
		session     *model.Session
		order       *model.Order
		event       *model.Event
		freeClasses []string
		ticket      bool
		eventmatch  bool
		jw          json.Writer
		cnames      []string
		classes     = map[string]*class{}
	)
	// Must have PrivSell to use this API.
	if session = auth.GetSession(tx, w, r, model.PrivSell); session == nil {
		return
	}
	// Make sure the requisite event exists.
	if event = tx.FetchEvent(eventID); event == nil {
		NotFoundError(tx, w)
		return
	}
	// Get the requested order.  It could be either an order number or an
	// order token.
	if oid, err := strconv.Atoi(token); err == nil {
		order = tx.FetchOrder(model.OrderID(oid))
	} else {
		order = tx.FetchOrderByToken(token)
	}
	if order == nil {
		sendError(tx, w, "Order not found")
		return
	}
	// Also search for any free ticket classes.
	for _, p := range tx.FetchProductsByEvent(event) {
		for _, sku := range p.SKUs {
			if sku.Price != 0 || sku.Coupon != "" || sku.Quantity != 1 || sku.MembersOnly {
				continue
			}
			// Note deliberately ignoring SalesStart..SalesEnd
			// range.  Free student tickets are usually set up with
			// a range that never matches so they don't appear as
			// explicitly orderable.
			freeClasses = append(freeClasses, p.TicketClass)
		}
	}
	commit(tx)
	// Search the order lines looking for tickets to the desired event.
	for _, ol := range order.Lines {
		for _, pe := range ol.Product.Events {
			ticket = true
			if pe.Event.ID == eventID {
				eventmatch = true
				var cname = ol.Product.TicketClass
				var cl = classes[cname]
				if cl == nil {
					cl = new(class)
					classes[cname] = cl
				}
				cl.val += ol.Quantity
				for _, t := range ol.Tickets {
					if t.Used.IsZero() {
						cl.max++
					} else {
						cl.used++
					}
				}
			}
		}
	}
	// Normalize the availability map.
	for cname, cl := range classes {
		if cl.val > cl.max {
			cl.val = cl.max
		}
		if cl.max == 0 {
			delete(classes, cname)
		}
	}
	// Add in the free ticket classes if they aren't already there.  But
	// don't do so if the class map is empty; that means the ticket is
	// invalid and we don't want to make it look valid.
	if len(classes) != 0 {
		for _, f := range freeClasses {
			if cl := classes[f]; cl != nil {
				if cl.max < 6 {
					cl.max = 6
				}
			} else {
				classes[f] = &class{max: 6}
			}
		}
	}
	// Sort the classes by name.
	for cl := range classes {
		cnames = append(cnames, cl)
	}
	sort.Strings(cnames)
	// Send the results.
	w.Header().Set("Content-Type", "application/json")
	jw = json.NewWriter(w)
	jw.Object(func() {
		jw.Prop("id", int(order.ID))
		if order.Name != "" {
			jw.Prop("name", order.Name)
		}
		if !ticket {
			jw.Prop("error", "Not a ticket order")
		} else if !eventmatch {
			jw.Prop("error", "Wrong event")
		} else if len(classes) == 0 {
			jw.Prop("error", "Ticket already used")
		}
		jw.Prop("classes", func() {
			jw.Array(func() {
				for _, cname := range cnames {
					class := classes[cname]
					jw.Object(func() {
						jw.Prop("name", cname)
						jw.Prop("val", class.val)
						jw.Prop("max", class.max)
						jw.Prop("used", class.used)
					})
				}
			})
		})
	})
	jw.Close()
}
