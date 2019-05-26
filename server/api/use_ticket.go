package api

import (
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/rothskeller/json"

	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
)

// UseTicket is the API used by the scanner app to mark tickets as being used.
func UseTicket(tx db.Tx, w http.ResponseWriter, r *http.Request, eventID model.EventID, token string) {
	var (
		session *model.Session
		order   *model.Order
		event   *model.Event
		jr      *json.Reader
		usage   = map[string]int{}
		err     error
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
		NotFoundError(tx, w)
		return
	}
	// Read the usage parameters from the request body.  It is a (JSON) map
	// from ticket class name to usage count (where the empty string is used
	// as the name for General Admission).
	jr = json.NewReader(r.Body)
	if err = jr.Read(json.ObjectHandler(func(key string) json.Handlers {
		return json.IntHandler(func(i int) { usage[key] = i })
	})); err != nil {
		BadRequestError(tx, w, err.Error())
		return
	}
	if len(usage) == 0 {
		BadRequestError(tx, w, "no ticket classes specified")
		return
	}
	// Walk through each ticket class, marking the requisite tickets used.
	for cname := range usage {
		if !useTicketClass(order, event, cname, usage[cname]) {
			sendError(tx, w, "Ticket already used")
			return
		}
	}
	// Clean up and return success.
	tx.SaveOrder(order)
	commit(tx)
	log.Printf("%s USE TICKETS order:%d event:%s %v", session.Username, order.ID, eventID, usage)
	w.WriteHeader(http.StatusNoContent)
}

// useTicketClass marks the specified number of tickets of the specified order,
// event, and class used.  It returns true if successful, or false if the
// requested number of tickets isn't available.
func useTicketClass(order *model.Order, event *model.Event, cname string, count int) bool {
	var (
		lines []*model.OrderLine
		now   = time.Now()
	)
	// Find the order lines that can be used for this request.
	for _, ol := range order.Lines {
		if ol.Product.TicketClass != cname {
			continue
		}
		var found bool
		for _, e := range ol.Product.Events {
			if e.ID == event.ID {
				found = true
				break
			}
		}
		if !found {
			continue
		}
		lines = append(lines, ol)
	}
	// Sort the lines based on how many events they can be used for; we want
	// to use the most specific tickets first.
	sort.Slice(lines, func(i, j int) bool {
		return len(lines[i].Product.Events) < len(lines[j].Product.Events)
	})
	// Walk through the lines, marking tickets used.
	for _, ol := range lines {
		for _, t := range ol.Tickets {
			if !t.Used.IsZero() {
				continue
			}
			t.Event = event
			t.Used = now
			count--
			if count == 0 {
				return true
			}
		}
	}
	return false
}
