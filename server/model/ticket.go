package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"scholacantorum.org/orders/db"
)

// TicketID is the unique identifier of a ticket.
type TicketID int

// Ticket represents a ticket to an event.  See db/schema.sql for details.
type Ticket struct {
	ID       TicketID   `json:"id"`
	Token    string     `json:"token"`
	SaleLine SaleLineID `json:"saleLine"`
	Event    EventID    `json:"event"`
	Used     time.Time  `json:"used"`
}

// ticketColumns is the list of columns of the ticket table.
var ticketColumns = `id, token, sale_line, event, used`

// scanTicket scans a ticket table row.
func scanTicket(scanner interface{ Scan(...interface{}) error }, t *Ticket) error {
	return scanner.Scan(&t.ID, &t.Token, &t.SaleLine, (*db.ID)(&t.Event), (*db.Time)(&t.Used))
}

// Save saves a ticket to the database.
func (t *Ticket) Save(tx *sql.Tx) {
	var (
		q   strings.Builder
		res sql.Result
		nid int64
		err error
	)
	q.WriteString(`INSERT OR REPLACE INTO ticket (`)
	q.WriteString(ticketColumns)
	q.WriteString(`) VALUES (?,?,?,?,?)`)
	res, err = tx.Exec(q.String(), db.ID(t.ID), t.Token, t.SaleLine, db.ID(t.Event), db.Time(t.Used))
	panicOnError(err)
	if t.ID == 0 {
		if nid, err = res.LastInsertId(); err != nil {
			panic(err)
		}
		t.ID = TicketID(nid)
	}
}

// Delete deletes a ticket.
func (t *Ticket) Delete(tx *sql.Tx) {
	panicOnNoRows(tx.Exec(`DELETE FROM ticket WHERE id=?`, t.ID))
}

// FetchTicket returns the ticket with the specified ID.  It returns nil if no
// such ticket exists.
func FetchTicket(tx *sql.Tx, id TicketID) (t *Ticket) {
	var (
		q   strings.Builder
		err error
	)
	t = new(Ticket)
	q.WriteString(`SELECT `)
	q.WriteString(ticketColumns)
	q.WriteString(` FROM ticket WHERE id=?`)
	switch err = scanTicket(tx.QueryRow(q.String(), id), t); err {
	case nil:
		return t
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}
}

// FetchTickets calls the supplied function for each ticket meeting the
// specified criteria.  Tickets are returned in ID order.  Offset and limit,
// when nonzero, limit the calls to the specified page of results.
func FetchTickets(tx *sql.Tx, offset, limit int, fn func(*Ticket), criteria string, args ...interface{}) {
	var (
		q    strings.Builder
		t    Ticket
		rows *sql.Rows
		err  error
	)
	q.WriteString(`SELECT `)
	q.WriteString(ticketColumns)
	q.WriteString(` FROM ticket`)
	if criteria != "" {
		q.WriteString(` WHERE `)
		q.WriteString(criteria)
	}
	q.WriteString(` ORDER BY id`)
	if limit != 0 {
		fmt.Fprintf(&q, ` LIMIT %d`, limit)
	}
	if offset != 0 {
		fmt.Fprintf(&q, ` OFFSET %d`, offset)
	}
	rows, err = tx.Query(q.String(), args...)
	panicOnError(err)
	for rows.Next() {
		panicOnError(scanTicket(rows, &t))
		fn(&t)
	}
	panicOnError(rows.Err())
}

// FetchTicketsWithToken returns the list of tickets with the specified token.
func FetchTicketsWithToken(tx *sql.Tx, token string) (tickets []*Ticket) {
	FetchTickets(tx, 0, 0, func(t *Ticket) {
		var copy = *t
		tickets = append(tickets, &copy)
	}, `token=?`, token)
	return tickets
}

// FetchTicketsForEvent returns the list of tickets to the event with the
// specified ID.  Note that this does not include flexible-use tickets that
// *might* be used for the target event.
func FetchTicketsForEvent(tx *sql.Tx, id EventID) (tickets []*Ticket) {
	FetchTickets(tx, 0, 0, func(t *Ticket) {
		var copy = *t
		tickets = append(tickets, &copy)
	}, `event=?`, id)
	return tickets
}
