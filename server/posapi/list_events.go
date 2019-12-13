package posapi

import (
	"github.com/mailru/easyjson/jwriter"
	"scholacantorum.org/orders/api"
	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/model"
)

type freeEntryData struct {
	Product model.ProductID
	Class   string
}

// ListEvents handles GET /api/event requests.
func ListEvents(r *api.Request) error {
	var (
		events      []*model.Event
		jw          jwriter.Writer
		freeEntries = map[model.EventID][]string{}
	)
	// Getting events needs PrivSetupOrders, PrivInPersonSales, or
	// PrivScanTickets.  Here we assume that anyone with PrivSetupOrders or
	// PrivInPersonSales will also have PrivScanTickets.
	if r.Privileges&(model.PrivSetupOrders|model.PrivInPersonSales|model.PrivScanTickets) == 0 {
		return auth.Forbidden
	}
	// We only return events scheduled in the future.  (Actually this
	// includes all events starting at the beginning of the date of the
	// call.)
	events = r.Tx.FetchFutureEvents()
	// We also retrieve the ticket class(es) that allow free entry to the
	// event, if any.
	for _, e := range events {
		freeEntries[e.ID] = getFreeEntries(r, e)
	}
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json")
	jw.RawByte('[')
	for i, e := range events {
		if i != 0 {
			jw.RawByte(',')
		}
		emitListedEvent(&jw, e, freeEntries[e.ID])
	}
	jw.RawByte(']')
	jw.DumpTo(r)
	return nil
}

// getFreeEntries looks for tickets to the event that can be ordered at zero
// cost, which will generally be the ticket classes that gets free entry (e.g.
// students).  If it finds any, it returns a list of the ticket class names;
// otherwise, it returns nil.
func getFreeEntries(r *api.Request, event *model.Event) (list []string) {
	var seen = map[string]bool{}
	for _, p := range getFreeClasses(r, event) {
		if !seen[p.TicketClass] {
			list = append(list, p.TicketClass)
		}
	}
	return list
}

// emitListedEvent emits the JSON for one event in a list.
func emitListedEvent(jw *jwriter.Writer, e *model.Event, freeEntries []string) {
	jw.RawString(`{"id":`)
	jw.String(string(e.ID))
	jw.RawString(`,"name":`)
	jw.String(e.Name)
	jw.RawString(`,"start"`)
	jw.Raw(e.Start.MarshalJSON())
	if freeEntries != nil {
		jw.RawString(`,"freeEntries":[`)
		for i, fe := range freeEntries {
			if i != 0 {
				jw.RawByte(',')
			}
			jw.String(fe)
		}
	}
}
