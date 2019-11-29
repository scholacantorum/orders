package posapi

import (
	"net/http"
	"sort"
	"strings"

	"github.com/rothskeller/json"

	"scholacantorum.org/orders/api"
	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
)

// ListEventOrders returns a list of orders that include tickets valid at the
// specified event.  It is used to support the Will Call feature of the
// at-the-door sales app.
func ListEventOrders(tx db.Tx, w http.ResponseWriter, r *http.Request, eventID model.EventID) {
	var (
		session *model.Session
		event   *model.Event
		list    []db.EventOrder
		jw      json.Writer
	)
	// Get current session data, if any.
	if session = auth.GetSession(tx, w, r, model.PrivInPersonSales); session == nil {
		return
	}
	// Get the event whose orders we're supposed to list.
	if event = tx.FetchEvent(eventID); event == nil {
		api.NotFoundError(tx, w)
		return
	}
	list = tx.FetchEventOrders(event)
	api.Commit(tx)
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
	w.Header().Set("Content-Type", "application/json")
	jw = json.NewWriter(w)
	jw.Array(func() {
		for _, eo := range list {
			jw.Object(func() {
				jw.Prop("id", int(eo.ID))
				jw.Prop("name", eo.Name)
			})
		}
	})
	jw.Close()
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
