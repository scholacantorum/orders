package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
)

// CountMetaUsage increments the specified meta-usage counter for the specified
// event.
func CountMetaUsage(tx db.Tx, w http.ResponseWriter, r *http.Request, eventID model.EventID, counter string) {
	var (
		session *model.Session
		event   *model.Event
		result  int
	)
	// Must have PrivSell to use this API.
	if session = auth.GetSession(tx, w, r, model.PrivSell); session == nil {
		return
	}
	// Get the requested event.
	if event = tx.FetchEvent(eventID); event == nil {
		NotFoundError(tx, w)
		return
	}
	switch counter {
	case "door":
		event.DoorSales++
		result = event.DoorSales
	case "free":
		event.FreeEntries++
		result = event.FreeEntries
	default:
		BadRequestError(tx, w, "unknown counter "+counter)
		return
	}
	tx.SaveEvent(event)
	commit(tx)
	log.Printf("%s INCREMENT %s for event:%s %d", session.Username, strings.ToUpper(counter), eventID, result)
	fmt.Fprint(w, result)
}
