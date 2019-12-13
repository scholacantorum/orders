package posapi

import (
	"sort"
	"strings"

	"github.com/mailru/easyjson/jwriter"
	"scholacantorum.org/orders/api"
	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
)

// ListEventOrders returns a list of orders that include tickets valid at the
// specified event.  It is used to support the Will Call feature of the
// at-the-door sales app.
func ListEventOrders(r *api.Request, eventID model.EventID) error {
	var (
		event *model.Event
		list  []db.EventOrder
		jw    jwriter.Writer
	)
	if r.Privileges&model.PrivInPersonSales == 0 {
		return auth.Forbidden
	}
	// Get the event whose orders we're supposed to list.
	if event = r.Tx.FetchEvent(eventID); event == nil {
		return api.NotFound
	}
	list = r.Tx.FetchEventOrders(event)
	r.Tx.Commit()
	for i := range list {
		list[i].Name = lastNameFirst(list[i].Name)
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].Name < list[j].Name {
			return true
		}
		if list[i].Name > list[j].Name {
			return false
		}
		return list[i].ID < list[j].ID
	})
	r.Header().Set("Content-Type", "application/json")
	jw.RawByte('[')
	for i, eo := range list {
		if i != 0 {
			jw.RawByte(',')
		}
		jw.RawString(`{"id":`)
		jw.Int(int(eo.ID))
		jw.RawString(`,"name":`)
		jw.String(eo.Name)
		jw.RawByte('}')
	}
	jw.RawByte(']')
	jw.DumpTo(r)
	return nil
}

func lastNameFirst(name string) string {
	var comma, space int
	var suffix string

	name = strings.TrimSpace(name)
	if comma = strings.LastIndexByte(name, ','); comma >= 0 {
		name, suffix = strings.TrimSpace(name[:comma]), strings.TrimSpace(name[comma+1:])
	}
	if space = strings.LastIndexByte(name, ' '); space >= 0 {
		if suffix != "" {
			return name[space+1:] + ", " + name[:space] + ", " + suffix
		}
		return name[space+1:] + ", " + name[:space]
	}
	if suffix != "" {
		return name + ", " + suffix
	}
	return name
}
