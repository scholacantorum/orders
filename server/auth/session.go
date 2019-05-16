package auth

import (
	"database/sql"
	"net/http"

	"scholacantorum.org/orders/model"
)

// GetSession verifies that the HTTP request identifies a valid session, and
// that session has all of the privileges specified in the priv bitmask.  (The
// mask may be zero to check for a valid session without requiring specific
// privileges.)  If so, GetSession returns the session details.  If not,
// GetSession emits an appropriate error, rolls back the transaction, and
// returns nil.
func GetSession(tx *sql.Tx, w http.ResponseWriter, r *http.Request, priv model.Privilege) (session *model.Session) {
	// TODO
	return &model.Session{Username: "guest", Privileges: 0xFF}
}

// HasSession returns whether the HTTP request identifies a session.  (The
// session is not necessarily valid.)
func HasSession(r *http.Request) bool {
	// TODO
	return true
}
