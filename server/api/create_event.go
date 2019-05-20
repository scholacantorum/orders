package api

import (
	"encoding/json"
	"log"
	"net/http"

	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
)

// CreateEvent handles POST /api/event requests.
func CreateEvent(tx db.Tx, w http.ResponseWriter, r *http.Request) {
	var (
		session *model.Session
		event   *model.Event
		err     error
	)
	if session = auth.GetSession(tx, w, r, model.PrivSetup); session == nil {
		return
	}
	if err = json.NewDecoder(r.Body).Decode(&event); err != nil {
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
	log.Printf("%s CREATE EVENT %s", session.Username, toJSON(event))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}
