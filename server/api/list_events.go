package api

import (
	"net/http"
	"time"

	"github.com/rothskeller/json"

	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
)

type freeEntryData struct {
	Product model.ProductID
	Class   string
}

// ListEvents handles GET /api/event requests.
func ListEvents(tx db.Tx, w http.ResponseWriter, r *http.Request) {
	var (
		session     *model.Session
		events      []*model.Event
		future      bool
		freeEntry   bool
		jw          json.Writer
		freeEntries = map[model.EventID]freeEntryData{}
	)
	// Getting events needs either PrivSetup or PrivSell.  Here we assume
	// that anyone with PrivSetup will also have PrivSell.
	if session = auth.GetSession(tx, w, r, model.PrivSell); session == nil {
		return
	}
	// If future is specified as a query parameter, then we should only
	// return events scheduled in the future.  This is used by the scanner
	// app.  (Actually this includes all events starting at the beginning of
	// the date of the call.)
	future = r.FormValue("future") != ""
	if future {
		events = tx.FetchFutureEvents()
	} else {
		events = tx.FetchEvents()
	}
	// If freeEntry is specified as a query parameter, then we should
	// retrieve the product and ticket class (if any) that allows free
	// entry to the event.
	freeEntry = r.FormValue("freeEntry") != ""
	if freeEntry {
		for _, e := range events {
			freeEntries[e.ID] = getFreeEntry(tx, e)
		}
	}
	commit(tx)
	w.Header().Set("Content-Type", "application/json")
	jw = json.NewWriter(w)
	jw.Array(func() {
		for _, e := range events {
			emitListedEvent(jw, e, freeEntries[e.ID])
		}
	})
	jw.Close()
}

// getFreeEntry looks for tickets to the event that can be ordered at zero cost,
// which will generally be the ticket class that gets free entry (e.g.
// students).  If it finds one, it returns the product ID and ticket class name;
// otherwise, it returns empty strings.
func getFreeEntry(tx db.Tx, event *model.Event) freeEntryData {
	for _, product := range tx.FetchProductsByEvent(event) {
		if product.TicketClass == "" {
			continue
		}
		for _, sku := range product.SKUs {
			if sku.Coupon == "" && !sku.MembersOnly && sku.Quantity == 1 && sku.Price == 0 {
				return freeEntryData{product.ID, product.TicketClass}
			}
			// Note that we're deliberately ignoring SalesStart and
			// SalesEnd.  SalesEnd is usually before at-the-door
			// sales start.
		}
	}
	return freeEntryData{}
}

// emitListedEvent emits the JSON for one event in a list.
func emitListedEvent(jw json.Writer, e *model.Event, fe freeEntryData) {
	jw.Object(func() {
		jw.Prop("id", string(e.ID))
		jw.Prop("name", e.Name)
		jw.Prop("start", e.Start.Format(time.RFC3339))
		if fe.Product != "" {
			jw.Prop("freeEntryProduct", string(fe.Product))
			jw.Prop("freeEntryClass", fe.Class)
		}
	})
}
