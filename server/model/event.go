package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"scholacantorum.org/orders/db"
)

// EventID is the unique identifier of an event.
type EventID int

// Event represents an event to which tickets are sold.  See db/schema.sql for
// details.
type Event struct {
	ID        EventID   `json:"id"`
	MembersID int       `json:"membersID"`
	Name      string    `json:"name"`
	Start     time.Time `json:"start"`
	Capacity  int       `json:"capacity"`
}

// eventColumns is the list of columns of the event table.
var eventColumns = `id, membersID, name, start, capacity`

// scanEvent scans an event table row.
func scanEvent(scanner interface{ Scan(...interface{}) error }, e *Event) error {
	var (
		membersID db.ID
		err       error
	)
	err = scanner.Scan(&e.ID, &membersID, &e.Name, (*db.Time)(&e.Start), &e.Capacity)
	if err != nil {
		return err
	}
	e.MembersID = int(membersID)
	return nil
}

// Save saves an event to the database.
func (e *Event) Save(tx *sql.Tx) {
	var (
		q   strings.Builder
		res sql.Result
		nid int64
		err error
	)
	q.WriteString(`INSERT OR REPLACE INTO event (`)
	q.WriteString(eventColumns)
	q.WriteString(`) VALUES (?,?,?,?,?)`)
	res, err = tx.Exec(q.String(), db.ID(e.ID), db.ID(e.MembersID), e.Name, db.Time(e.Start), e.Capacity)
	panicOnError(err)
	if e.ID == 0 {
		if nid, err = res.LastInsertId(); err != nil {
			panic(err)
		}
		e.ID = EventID(nid)
	}
}

// Delete deletes an event.
func (e *Event) Delete(tx *sql.Tx) {
	panicOnNoRows(tx.Exec(`DELETE FROM event WHERE id=?`, e.ID))
}

// FetchEvent returns the event with the specified ID.  It returns nil if no
// such event exists.
func FetchEvent(tx *sql.Tx, id EventID) (e *Event) {
	var (
		q   strings.Builder
		err error
	)
	e = new(Event)
	q.WriteString(`SELECT `)
	q.WriteString(eventColumns)
	q.WriteString(` FROM event WHERE id=?`)
	switch err = scanEvent(tx.QueryRow(q.String(), id), e); err {
	case nil:
		return e
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}
}

// FetchEventWithMembersID returns the event with the specified MembersID.  It
// returns nil if no such event exists.
func FetchEventWithMembersID(tx *sql.Tx, membersID int) (e *Event) {
	var (
		q   strings.Builder
		err error
	)
	e = new(Event)
	q.WriteString(`SELECT `)
	q.WriteString(eventColumns)
	q.WriteString(` FROM event WHERE membersID=?`)
	switch err = scanEvent(tx.QueryRow(q.String(), membersID), e); err {
	case nil:
		return e
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}
}

// FetchEvents calls the supplied function for each event meeting the specified
// criteria.  Events are returned in chronological order by start time.  Offset
// and limit, when nonzero, limit the calls to the specified page of results.
func FetchEvents(tx *sql.Tx, offset, limit int, fn func(*Event), criteria string, args ...interface{}) {
	var (
		q    strings.Builder
		e    Event
		rows *sql.Rows
		err  error
	)
	q.WriteString(`SELECT `)
	q.WriteString(eventColumns)
	q.WriteString(` FROM event`)
	if criteria != "" {
		q.WriteString(` WHERE `)
		q.WriteString(criteria)
	}
	q.WriteString(` ORDER BY start`)
	if limit != 0 {
		fmt.Fprintf(&q, ` LIMIT %d`, limit)
	}
	if offset != 0 {
		fmt.Fprintf(&q, ` OFFSET %d`, offset)
	}
	rows, err = tx.Query(q.String(), args...)
	panicOnError(err)
	for rows.Next() {
		panicOnError(scanEvent(rows, &e))
		fn(&e)
	}
	panicOnError(rows.Err())
}
