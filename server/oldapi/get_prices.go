package oldapi

import (
	"net/http"

	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
	"scholacantorum.org/orders/payapi"
	"scholacantorum.org/orders/posapi"
)

// GetPrices is two different APIs, one for payment forms and a different one
// for point of sale.  They used to share the same entrypoint.  For backward
// compatibility, this function serves that entrypoint, detects which one is
// intended, and forwards the call to the correct handler.
func GetPrices(tx db.Tx, w http.ResponseWriter, r *http.Request) {
	if eventID := r.FormValue("event"); eventID != "" {
		posapi.GetEventPrices(tx, w, r, model.EventID(eventID))
	} else {
		payapi.GetPrices(tx, w, r)
	}
}
