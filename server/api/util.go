package api

import (
	"net/http"

	"github.com/rothskeller/json"

	"scholacantorum.org/orders/db"
)

// NotFoundError returns a 404 Not Found.
func NotFoundError(tx db.Tx, w http.ResponseWriter) {
	tx.Rollback()
	http.Error(w, "404 Not Found", http.StatusNotFound)
}

// BadRequestError returns a 400 Bad Request error.
func BadRequestError(tx db.Tx, w http.ResponseWriter, reason string) {
	tx.Rollback()
	if reason != "" {
		http.Error(w, "400 Bad Request: "+reason, http.StatusBadRequest)
	} else {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
	}
}

// ForbiddenError returns a 403 Forbidden error.
func ForbiddenError(tx db.Tx, w http.ResponseWriter) {
	tx.Rollback()
	http.Error(w, "403 Forbidden", http.StatusForbidden)
}

// sendError sends an error message as a JSON object with an "error" key.  For
// convenience, it also rolls back the transaction.
func sendError(tx db.Tx, w http.ResponseWriter, message string) {
	tx.Rollback()
	w.Header().Set("Content-Type", "application/json")
	var jw = json.NewWriter(w)
	jw.Object(func() {
		jw.Prop("error", message)
	})
	jw.Close()
}

// commit commits the transaction.
func commit(tx db.Tx) {
	if err := tx.Commit(); err != nil {
		panic(err)
	}
}
