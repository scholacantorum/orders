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
		session         *model.Session
		events          []*model.Event
		future          bool
		wantFreeEntries bool
		jw              json.Writer
		freeEntries     = map[model.EventID][]string{}
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
	// If freeEntries is specified as a query parameter, then we should
	// retrieve the ticket class(es) that allow free entry to the event, if
	// any.
	wantFreeEntries = r.FormValue("freeEntries") != ""
	if wantFreeEntries {
		for _, e := range events {
			freeEntries[e.ID] = getFreeEntries(tx, e)
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

// getFreeEntries looks for tickets to the event that can be ordered at zero
// cost, which will generally be the ticket classes that gets free entry (e.g.
// students).  If it finds any, it returns a list of the ticket class names;
// otherwise, it returns nil.
func getFreeEntries(tx db.Tx, event *model.Event) (list []string) {
	var seen = map[string]bool{}
	for _, product := range tx.FetchProductsByEvent(event) {
		for _, sku := range product.SKUs {
			if sku.Coupon == "" && !sku.MembersOnly && sku.Price == 0 && !seen[product.TicketClass] {
				seen[product.TicketClass] = true
				list = append(list, product.TicketClass)
			}
			// Note that we're deliberately ignoring SalesStart and
			// SalesEnd.  Student tickets usually are deliberately
			// out of range so they don't show up for explicit sale.
		}
	}
	return list
}

// emitListedEvent emits the JSON for one event in a list.
func emitListedEvent(jw json.Writer, e *model.Event, freeEntries []string) {
	jw.Object(func() {
		jw.Prop("id", string(e.ID))
		jw.Prop("name", e.Name)
		jw.Prop("start", e.Start.Format(time.RFC3339))
		if freeEntries != nil {
			jw.Prop("freeEntries", func() {
				jw.Array(func() {
					for _, fe := range freeEntries {
						jw.String(fe)
					}
				})
			})
		}
	})
}
