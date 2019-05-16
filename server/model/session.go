package model

import (
	"database/sql"
	"time"

	"scholacantorum.org/orders/db"
)

// Session represents an active login session.  See db/schema.sql for details.
type Session struct {
	Token      string
	Username   string
	Expires    time.Time
	Privileges Privilege
}

// Create creates a session in the database.
func (s *Session) Create(tx *sql.Tx) {
	panicOnNoRows(tx.Exec(
		`INSERT INTO session (token, username, expires, privileges) VALUES (?,?,?,?)`,
		s.Token, s.Username, db.Time(s.Expires), s.Privileges))
}

// Delete deletes a session.
func (s *Session) Delete(tx *sql.Tx) {
	panicOnNoRows(tx.Exec(`DELETE FROM session WHERE token=?`, s.Token))
}

// ExpireSessions deletes all expired sessions.
func ExpireSessions(tx *sql.Tx) {
	panicOnExecError(tx.Exec(`DELETE FROM session WHERE expires<?`, time.Now().Unix()))
}

// FetchSession returns the session with the specified Token.  It returns nil if
// no such session exists.
func FetchSession(tx *sql.Tx, token string) (s *Session) {
	var err error

	s = new(Session)
	err = tx.QueryRow(`SELECT username, expires, privileges FROM session WHERE token=?`, token).
		Scan(&s.Username, (*db.Time)(&s.Expires), &s.Privileges)
	switch err {
	case nil:
		s.Token = token
		return s
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}
}
