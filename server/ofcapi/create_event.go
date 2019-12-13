package ofcapi

import (
	"errors"
	"strconv"
	"time"

	"github.com/mailru/easyjson"

	"scholacantorum.org/orders/api"
	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/model"
)

// CreateEvent handles POST /ofcapi/event requests.
func CreateEvent(r *api.Request) error {
	var (
		event *model.Event
		err   error
	)
	if r.Privileges&model.PrivSetupOrders == 0 {
		return auth.Forbidden
	}
	if event, err = parseCreateEvent(r); err != nil {
		return err
	}
	if r.Tx.FetchEvent(event.ID) != nil {
		return errors.New("duplicate event ID")
	}
	if event.MembersID != 0 && r.Tx.FetchEventByMembersID(event.MembersID) != nil {
		return errors.New("membersID already in use")
	}
	r.Tx.SaveEvent(event)
	r.Tx.Commit()
	easyjson.MarshalToHTTPResponseWriter(event, r)
	return nil
}

func parseCreateEvent(r *api.Request) (e *model.Event, err error) {
	e = new(model.Event)
	if e.ID = model.EventID(r.FormValue("id")); e.ID == "" {
		return nil, errors.New("missing ID")
	}
	if mid, err := strconv.Atoi(r.FormValue("membersID")); err == nil && mid > 0 {
		e.MembersID = mid
	} else {
		return nil, errors.New("missing or invalid membersID")
	}
	if e.Name = r.FormValue("name"); e.Name == "" {
		return nil, errors.New("missing name")
	}
	if e.Series = r.FormValue("series"); e.Series == "" {
		return nil, errors.New("missing series")
	}
	if start, err := time.Parse(time.RFC3339, r.FormValue("start")); err == nil {
		e.Start = start.In(time.Local)
	} else {
		return nil, errors.New("missing or invalid start")
	}
	if caps := r.FormValue("capacity"); caps != "" {
		if cap, err := strconv.Atoi(r.FormValue("capacity")); err == nil && cap >= 0 {
			e.Capacity = cap
		} else {
			return nil, errors.New("invalid capacity")
		}
	}
	return e, nil
}
