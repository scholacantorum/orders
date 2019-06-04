package db

import (
	"database/sql"
	"strings"
	"time"

	"scholacantorum.org/orders/model"
)

// sessionColumns is the list of columns of the session table.
var sessionColumns = `token, username, expires, member, privileges`

// scanSession scans a session table row.
func (tx Tx) scanSession(scanner interface{ Scan(...interface{}) error }, s *model.Session) error {
	return scanner.Scan(&s.Token, &s.Username, (*Time)(&s.Expires), &s.Member, &s.Privileges)
}

// SaveSession saves a session to the database.
func (tx Tx) SaveSession(s *model.Session) {
	var (
		q strings.Builder
	)
	q.WriteString(`INSERT INTO session (`)
	q.WriteString(sessionColumns)
	q.WriteString(`) VALUES (?,?,?,?,?)`)
	panicOnExecError(tx.tx.Exec(q.String(), s.Token, s.Username, Time(s.Expires), s.Member, s.Privileges))
}

// ExpireSessions deletes all sessions that have expired.
func (tx Tx) ExpireSessions() {
	panicOnExecError(tx.tx.Exec(`DELETE FROM session WHERE expires<?`, Time(time.Now())))
}

// FetchSession returns the session with the specified Token.  It returns nil if
// no such session exists.
func (tx Tx) FetchSession(token string) (s *model.Session) {
	var (
		q   strings.Builder
		err error
	)
	s = new(model.Session)
	q.WriteString(`SELECT `)
	q.WriteString(sessionColumns)
	q.WriteString(` FROM session WHERE token=?`)
	switch err = tx.scanSession(tx.tx.QueryRow(q.String(), token), s); err {
	case nil:
		return s
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}
}
