package api

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/rothskeller/json"

	"scholacantorum.org/orders/db"
)

// NotFound is the error returned when a page is not found.
var NotFound = HTTPError(http.StatusNotFound, "404 Not Found")

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

// SendError sends an error message as a JSON object with an "error" key.  For
// convenience, it also rolls back the transaction.
func SendError(tx db.Tx, w http.ResponseWriter, message string) {
	tx.Rollback()
	w.Header().Set("Content-Type", "application/json")
	var jw = json.NewWriter(w)
	jw.Object(func() {
		jw.Prop("error", message)
	})
	jw.Close()
}
func sendError(tx db.Tx, w http.ResponseWriter, message string) {
	tx.Rollback()
	w.Header().Set("Content-Type", "application/json")
	var jw = json.NewWriter(w)
	jw.Object(func() {
		jw.Prop("error", message)
	})
	jw.Close()
}

// Commit commits the transaction.
func Commit(tx db.Tx) {
	tx.Commit()
}
func commit(tx db.Tx) {
	tx.Commit()
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// NewToken generates a random token string.
func NewToken() string {
	tval := rand.Intn(1000000000000)
	return fmt.Sprintf("%04d-%04d-%04d", tval/100000000, tval/10000%10000, tval%10000)
}
