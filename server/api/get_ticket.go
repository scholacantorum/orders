package api

import (
	"net/http"
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
		val int
		max int
	}
	var (
		session    *model.Session
		order      *model.Order
		ticket     bool
		eventmatch bool
		jw         json.Writer
		classes    = map[string]*class{}
	)
	// Must have PrivSell to use this API.
	if session = auth.GetSession(tx, w, r, model.PrivSell); session == nil {
		return
	}
	// Make sure the requisite event exists.
	if tx.FetchEvent(eventID) == nil {
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
	commit(tx)
	// Search the order lines looking for tickets to the desired event.
	for _, ol := range order.Lines {
		for _, e := range ol.Product.Events {
			ticket = true
			if e.ID == eventID {
				eventmatch = true
				var cname = ol.Product.TicketClass
				var cl = classes[cname]
				if cl == nil {
					cl = new(class)
					classes[cname] = cl
				}
				cl.val += ol.Quantity / len(ol.Product.Events)
				for _, t := range ol.Tickets {
					if t.Used.IsZero() {
						cl.max++
					}
				}
			}
		}
	}
	for cname, cl := range classes {
		if cl.val > cl.max {
			cl.val = cl.max
		}
		if cl.max == 0 {
			delete(classes, cname)
		}
	}
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
			jw.Object(func() {
				for cname, class := range classes {
					jw.Prop(cname, func() {
						jw.Object(func() {
							jw.Prop("val", class.val)
							jw.Prop("max", class.max)
						})
					})
				}
			})
		})
	})
	jw.Close()
}
