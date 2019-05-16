package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// NotFoundError returns a 404 Not Found.
func NotFoundError(tx *sql.Tx, w http.ResponseWriter) {
	tx.Rollback()
	http.Error(w, "404 Not Found", http.StatusNotFound)
}

// BadRequestError returns a 400 Bad Request error.
func BadRequestError(tx *sql.Tx, w http.ResponseWriter, reason string) {
	tx.Rollback()
	if reason != "" {
		http.Error(w, "400 Bad Request: "+reason, http.StatusBadRequest)
	} else {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
	}
}

// toJSON renders the supplied object as a JSON string.  This is primarily used
// for audit logging.
func toJSON(v interface{}) string {
	var (
		buf []byte
		err error
	)
	if buf, err = json.Marshal(v); err != nil {
		panic(err)
	}
	return string(buf)
}

// commit commits the transaction.
func commit(tx *sql.Tx) {
	if err := tx.Commit(); err != nil {
		panic(err)
	}
}
