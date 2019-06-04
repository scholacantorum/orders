package auth

import (
	"net/http"

	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
)

// GetSession verifies that the HTTP request identifies a valid session, and
// that session has all of the privileges specified in the priv bitmask.  (The
// mask may be zero to check for a valid session without requiring specific
// privileges.)  If so, GetSession returns the session details.  If not,
// GetSession emits an appropriate error, rolls back the transaction, and
// returns nil.
func GetSession(tx db.Tx, w http.ResponseWriter, r *http.Request, priv model.Privilege) (session *model.Session) {
	tx.ExpireSessions()
	if session = tx.FetchSession(r.Header.Get("Auth")); session == nil {
		tx.Rollback()
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		return nil
	}
	if session.Privileges&priv != priv {
		tx.Rollback()
		http.Error(w, "403 Forbidden", http.StatusForbidden)
		return nil
	}
	return session
}

// HasSession returns whether the HTTP request identifies a session.  (The
// session is not necessarily valid.)
func HasSession(r *http.Request) bool {
	return r.Header.Get("Auth") != ""
}
