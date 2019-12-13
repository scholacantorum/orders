package db

import (
	"database/sql"
	"strings"
	"time"

	"scholacantorum.org/orders/model"
)

// eventColumns is the list of columns of the event table.
var eventColumns = `id, members_id, name, series, start, capacity`

// scanEvent scans an event table row.
func (tx Tx) scanEvent(scanner interface{ Scan(...interface{}) error }, e *model.Event) error {
	var (
		membersID ID
		err       error
	)
	err = scanner.Scan(&e.ID, &membersID, &e.Name, &e.Series, (*Time)(&e.Start), &e.Capacity)
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
	q.WriteString(`) VALUES (?,?,?,?,?,?)`)
	panicOnExecError(tx.tx.Exec(q.String(), IDStr(e.ID), ID(e.MembersID), e.Name, e.Series, Time(e.Start), e.Capacity))
	tx.audit(model.AuditRecord{Event: e})
}

// DeleteEvent deletes an event.
func (tx Tx) DeleteEvent(e *model.Event) {
	panicOnNoRows(tx.tx.Exec(`DELETE FROM event WHERE id=?`, e.ID))
	tx.audit(model.AuditRecord{Event: &model.Event{ID: e.ID}})
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

// FetchEvents returns a list of all events.
func (tx Tx) FetchEvents() (events []*model.Event) {
	var (
		q    strings.Builder
		rows *sql.Rows
		err  error
	)
	q.WriteString(`SELECT `)
	q.WriteString(eventColumns)
	q.WriteString(` FROM event ORDER BY start`)
	rows, err = tx.tx.Query(q.String())
	panicOnError(err)
	for rows.Next() {
		var e model.Event
		panicOnError(tx.scanEvent(rows, &e))
		events = append(events, &e)
	}
	panicOnError(rows.Err())
	return events
}

// FetchFutureEvents returns a list of future events, in chronological order.
func (tx Tx) FetchFutureEvents() (events []*model.Event) {
	var (
		q    strings.Builder
		rows *sql.Rows
		err  error
	)
	q.WriteString(`SELECT `)
	q.WriteString(eventColumns)
	q.WriteString(` FROM event WHERE start > ? ORDER BY start`)
	rows, err = tx.tx.Query(q.String(), time.Now().Format("2006-01-02"))
	panicOnError(err)
	for rows.Next() {
		var e model.Event
		panicOnError(tx.scanEvent(rows, &e))
		events = append(events, &e)
	}
	panicOnError(rows.Err())
	return events
}

// FetchTicketCount returns the number of tickets allocated to the specified
// event.  These could be either tickets that have been used at the event, or
// tickets that are labeled for that event.
func (tx Tx) FetchTicketCount(event *model.Event) (count int) {
	panicOnError(tx.tx.QueryRow(`SELECT COUNT(*) FROM ticket WHERE event=?`, event.ID).Scan(&count))
	return count
}

// EventOrder is the type returned by FetchEventOrders (q.v.).
type EventOrder struct {
	ID   model.OrderID
	Name string
}

// FetchEventOrders returns a list of order IDs and customer names for will call
// searches.  The list is not sorted.  It contains those orders which are valid,
// have a customer name, and have either tickets used at the target event or
// unused tickets that could be used at it.
func (tx Tx) FetchEventOrders(event *model.Event) (list []EventOrder) {
	var (
		rows *sql.Rows
		err  error
	)
	rows, err = tx.tx.Query(`
SELECT DISTINCT o.id, o.name FROM ordert o, order_line ol, product_event pe, ticket t
WHERE pe.event=?1 AND pe.product=ol.product AND o.id=ol.orderid AND o.name != ''
AND o.valid AND t.order_line=ol.id AND (t.used='' OR t.event=?1)`, event.ID)
	panicOnError(err)
	for rows.Next() {
		var eo EventOrder
		panicOnError(rows.Scan(&eo.ID, &eo.Name))
		list = append(list, eo)
	}
	panicOnError(rows.Err())
	return list
}
