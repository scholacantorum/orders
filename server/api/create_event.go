package api

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/rothskeller/json"

	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
)

// CreateEvent handles POST /api/event requests.
func CreateEvent(tx db.Tx, w http.ResponseWriter, r *http.Request) {
	var (
		session *model.Session
		event   *model.Event
		out     []byte
		err     error
	)
	if session = auth.GetSession(tx, w, r, model.PrivSetupOrders); session == nil {
		return
	}
	if event, err = parseCreateEvent(r.Body); err != nil {
		BadRequestError(tx, w, err.Error())
		return
	}
	if event.ID == "" || event.MembersID < 0 || event.Name == "" || event.Start.IsZero() || event.Capacity < 0 {
		BadRequestError(tx, w, "invalid parameters")
		return
	}
	if tx.FetchEvent(event.ID) != nil {
		BadRequestError(tx, w, "duplicate event ID")
		return
	}
	if event.MembersID != 0 && tx.FetchEventByMembersID(event.MembersID) != nil {
		BadRequestError(tx, w, "membersID already in use")
		return
	}
	tx.SaveEvent(event)
	commit(tx)
	out = emitCreatedEvent(event)
	log.Printf("%s CREATE EVENT %s", session.Username, out)
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func parseCreateEvent(r io.Reader) (e *model.Event, err error) {
	var (
		jr = json.NewReader(r)
	)
	e = new(model.Event)
	err = jr.Read(json.ObjectHandler(func(key string) json.Handlers {
		switch key {
		case "id":
			return json.StringHandler(func(s string) { e.ID = model.EventID(s) })
		case "membersID":
			return json.IntHandler(func(i int) { e.MembersID = i })
		case "name":
			return json.StringHandler(func(s string) { e.Name = s })
		case "series":
			return json.StringHandler(func(s string) { e.Series = s })
		case "start":
			return json.TimeHandler(func(t time.Time) { e.Start = t })
		case "capacity":
			return json.IntHandler(func(i int) { e.Capacity = i })
		default:
			return json.RejectHandler()
		}
	}))
	return e, err
}

func emitCreatedEvent(e *model.Event) []byte {
	var (
		buf bytes.Buffer
		jw  = json.NewWriter(&buf)
	)
	jw.Object(func() {
		jw.Prop("id", string(e.ID))
		if e.MembersID != 0 {
			jw.Prop("membersID", e.MembersID)
		}
		jw.Prop("name", e.Name)
		jw.Prop("series", e.Series)
		jw.Prop("start", e.Start.Format(time.RFC3339))
		if e.Capacity != 0 {
			jw.Prop("capacity", e.Capacity)
		}
	})
	jw.Close()
	return buf.Bytes()
}
