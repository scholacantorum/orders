package db

import (
	"database/sql"
	"strings"

	"scholacantorum.org/orders/model"
)

// eventColumns is the list of columns of the event table.
var eventColumns = `id, members_id, name, start, capacity, door_sales, free_entries`

// scanEvent scans an event table row.
func (tx Tx) scanEvent(scanner interface{ Scan(...interface{}) error }, e *model.Event) error {
	var (
		membersID ID
		err       error
	)
	err = scanner.Scan(&e.ID, &membersID, &e.Name, (*Time)(&e.Start), &e.Capacity, &e.DoorSales, &e.FreeEntries)
	if err != nil {
		return err
	}
	e.MembersID = int(membersID)
	return nil
}

// SaveEvent saves an event to the database.
func (tx Tx) SaveEvent(e *model.Event) {
	var (
		q strings.Builder
	)
	q.WriteString(`INSERT OR REPLACE INTO event (`)
	q.WriteString(eventColumns)
	q.WriteString(`) VALUES (?,?,?,?,?,?,?)`)
	panicOnExecError(tx.tx.Exec(q.String(), IDStr(e.ID), ID(e.MembersID),
		e.Name, Time(e.Start), e.Capacity, e.DoorSales, e.FreeEntries))
}

// DeleteEvent deletes an event.
func (tx Tx) DeleteEvent(e *model.Event) {
	panicOnNoRows(tx.tx.Exec(`DELETE FROM event WHERE id=?`, e.ID))
}

// FetchEvent returns the event with the specified ID.  It returns nil if no
// such event exists.
func (tx Tx) FetchEvent(id model.EventID) (e *model.Event) {
	var (
		q   strings.Builder
		err error
	)
	e = new(model.Event)
	q.WriteString(`SELECT `)
	q.WriteString(eventColumns)
	q.WriteString(` FROM event WHERE id=?`)
	switch err = tx.scanEvent(tx.tx.QueryRow(q.String(), id), e); err {
	case nil:
		return e
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}
}

// FetchEventByMembersID returns the event with the specified MembersID.  It
// returns nil if no such event exists.
func (tx Tx) FetchEventByMembersID(membersID int) (e *model.Event) {
	var (
		q   strings.Builder
		err error
	)
	e = new(model.Event)
	q.WriteString(`SELECT `)
	q.WriteString(eventColumns)
	q.WriteString(` FROM event WHERE members_id=?`)
	switch err = tx.scanEvent(tx.tx.QueryRow(q.String(), membersID), e); err {
	case nil:
		return e
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}
}
