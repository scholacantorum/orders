package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
)

// GetTicket is the API used by the scanner app to get information about the
// ticket(s) on the specified order to the specified event.
func GetTicket(tx db.Tx, w http.ResponseWriter, r *http.Request, eventID model.EventID, token string) {
	type (
		Class struct {
			Val int `json:"val"`
			Max int `json:"max"`
		}
		Response struct {
			ID      model.OrderID     `json:"id"`
			Name    string            `json:"name"`
			Error   string            `json:"error,omitempty"`
			Classes map[string]*Class `json:"classes"`
		}
	)
	var (
		session    *model.Session
		order      *model.Order
		response   Response
		ticket     bool
		eventmatch bool
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
	response.ID = order.ID
	response.Name = order.Name
	response.Classes = make(map[string]*Class)
	// Search the order lines looking for tickets to the desired event.
	for _, ol := range order.Lines {
		for _, e := range ol.Product.Events {
			ticket = true
			if e.ID == eventID {
				eventmatch = true
				var cname = ol.Product.TicketClass
				var cl = response.Classes[cname]
				if cl == nil {
					cl = new(Class)
					response.Classes[cname] = cl
				}
				cl.Val += ol.Quantity / len(ol.Product.Events)
				for _, t := range ol.Tickets {
					if t.Used.IsZero() {
						cl.Max++
					}
				}
			}
		}
	}
	for cname, cl := range response.Classes {
		if cl.Val > cl.Max {
			cl.Val = cl.Max
		}
		if cl.Max == 0 {
			delete(response.Classes, cname)
		}
	}
	if !ticket {
		response.Error = "Not a ticket order"
	} else if !eventmatch {
		response.Error = "Wrong event"
	} else if len(response.Classes) == 0 {
		response.Error = "Ticket already used"
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&response)
}
