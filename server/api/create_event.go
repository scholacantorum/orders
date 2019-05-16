package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/model"
)

// CreateEvent handles POST /api/event requests.
func CreateEvent(tx *sql.Tx, w http.ResponseWriter, r *http.Request) {
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
	if event.ID != 0 || event.MembersID < 0 || event.Name == "" || event.Start.IsZero() || event.Capacity < 0 {
		BadRequestError(tx, w, "invalid parameters")
		return
	}
	if event.MembersID != 0 && model.FetchEventWithMembersID(tx, event.MembersID) != nil {
		BadRequestError(tx, w, "membersID already in use")
		return
	}
	event.Save(tx)
	commit(tx)
	log.Printf("%s CREATE EVENT %s", session.Username, toJSON(event))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}
