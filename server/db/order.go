package db

import (
	"database/sql"
	"strings"

	"scholacantorum.org/orders/model"
)

// orderColumns is the list of columns in the orderT table.
var orderColumns = `id, name, email, address, city, state, zip, phone, customer, member, created, flags, note, coupon, repeat`

// scanOrder scans an orderT table row.
func scanOrder(scanner interface{ Scan(...interface{}) error }, o *model.Order) error {
	return scanner.Scan(&o.ID, &o.Name, &o.Email, &o.Address, &o.City,
		&o.State, &o.Zip, &o.Phone, &o.Customer, &o.Member,
		(*Time)(&o.Created), &o.Flags, &o.Note, &o.Coupon,
		(*Time)(&o.Repeat))
}

// orderLineColumns is the list of columns in the order_line table (not
// including the parent reference and line number).
var orderLineColumns = `linenum, product, quantity, token, price`

// scanOrderLine scans an order line table row.
func scanOrderLine(scanner interface{ Scan(...interface{}) error }, ol *model.OrderLine) error {
	return scanner.Scan(&ol.Product, &ol.Quantity, &ol.Token, &ol.Price)
}

// FetchOrder returns the order with the specified ID.  It returns nil if no
// such order exists.  This includes fetching all subsidiary objects.
func (tx Tx) FetchOrder(id model.OrderID) (o *model.Order) {
	var (
		q     strings.Builder
		lrows *sql.Rows
		pid   model.ProductID
		token sql.NullString
		prows *sql.Rows
		pmid  model.PaymentID
		trows *sql.Rows
		eid   model.EventID
		err   error
		pmids = map[model.PaymentID]*model.Payment{}
	)
	o = new(model.Order)
	q.WriteString(`SELECT `)
	q.WriteString(orderColumns)
	q.WriteString(` FROM orderT WHERE id=?`)
	switch err = scanOrder(tx.tx.QueryRow(q.String(), id), o); err {
	case nil:
		break
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}
	prows, err = tx.tx.Query(`SELECT id, method, stripe, created, flags FROM payment WHERE orderid=? ORDER BY id`, o.ID)
	panicOnError(err)
	for prows.Next() {
		var p model.Payment
		panicOnError(prows.Scan(&p.ID, &p.Method, &p.Stripe, (*Time)(&p.Created), &p.Flags))
		o.Payments = append(o.Payments, &p)
		pmids[p.ID] = &p
	}
	panicOnError(prows.Err())
	lrows, err = tx.tx.Query(`SELECT id, product, quantity, token, price FROM order_line WHERE orderid=? ORDER BY id`, o.ID)
	panicOnError(err)
	for lrows.Next() {
		var ol model.OrderLine
		panicOnError(lrows.Scan(&ol.ID, &pid, &ol.Quantity, &token, &ol.Price))
		if token.Valid {
			ol.Token = token.String
		}
		ol.Product = tx.FetchProduct(pid)
		trows, err = tx.tx.Query(`SELECT id, event, used FROM ticket WHERE order_line=? ORDER BY id`, ol.ID)
		panicOnError(err)
		for trows.Next() {
			var t model.Ticket
			panicOnError(trows.Scan(&t.ID, (*IDStr)(&eid), (*Time)(&t.Used)))
			if eid != "" {
				t.Event = tx.FetchEvent(eid)
			}
			ol.Tickets = append(ol.Tickets, &t)
		}
		panicOnError(trows.Err())
		prows, err = tx.tx.Query(`SELECT payment, amount FROM payment_line WHERE order_line=?`, ol.ID)
		panicOnError(err)
		for prows.Next() {
			var pl model.PaymentLine
			panicOnError(prows.Scan(&pmid, &pl.Amount))
			pl.Payment = pmids[pmid]
			ol.PaymentLines = append(ol.PaymentLines, &pl)
		}
		panicOnError(prows.Err())
		o.Lines = append(o.Lines, &ol)
	}
	panicOnError(lrows.Err())
	return o
}

// FetchOrderByToken returns the order that contains a line with the specified
// token, and the index of that line within that order.  It returns nil, 0 if no
// such order exists.
func (tx Tx) FetchOrderByToken(token string) (o *model.Order, line int) {
	var (
		oid model.OrderID
		err error
	)
	switch err = tx.tx.QueryRow(`SELECT orderid FROM order_line WHERE token=?`, token).Scan(&oid); err {
	case nil:
		break
	case sql.ErrNoRows:
		return nil, 0
	default:
		panic(err)
	}
	o = tx.FetchOrder(oid)
	for i, ol := range o.Lines {
		if ol.Token == token {
			return o, i
		}
	}
	panic("order line not found within order")
}

// SaveOrder saves an order to the database.  This includes saving all
// order-specific subsidiary objects.
func (tx Tx) SaveOrder(o *model.Order) {
	var (
		q     strings.Builder
		res   sql.Result
		token sql.NullString
		err   error
	)
	q.WriteString(`INSERT OR REPLACE INTO orderT (`)
	q.WriteString(orderColumns)
	q.WriteString(`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	res, err = tx.tx.Exec(q.String(), ID(o.ID), o.Name, o.Email, o.Address,
		o.City, o.State, o.Zip, o.Phone, o.Customer, o.Member,
		Time(o.Created), o.Flags, o.Note, o.Coupon, Time(o.Repeat))
	panicOnError(err)
	if o.ID == 0 {
		o.ID = model.OrderID(lastInsertID(res))
	}
	for _, p := range o.Payments {
		res, err = tx.tx.Exec(
			`INSERT OR REPLACE INTO payment (id, orderid, method, stripe, created, flags) VALUES (?,?,?,?,?,?)`,
			ID(p.ID), o.ID, p.Method, p.Stripe, Time(p.Created), p.Flags)
		panicOnError(err)
		if p.ID == 0 {
			p.ID = model.PaymentID(lastInsertID(res))
		}
	}
	for _, ol := range o.Lines {
		token = sql.NullString{String: ol.Token, Valid: ol.Token != ""}
		res, err = tx.tx.Exec(
			`INSERT OR REPLACE INTO order_line (id, orderid, product, quantity, token, price) VALUES (?,?,?,?,?,?)`,
			ID(ol.ID), o.ID, ol.Product.ID, ol.Quantity, token, ol.Price)
		panicOnError(err)
		if ol.ID == 0 {
			ol.ID = model.OrderLineID(lastInsertID(res))
		} else {
			panicOnExecError(tx.tx.Exec(`DELETE FROM ticket WHERE order_line=?`, ol.ID))
		}
		for _, t := range ol.Tickets {
			var eid model.EventID
			if t.Event != nil {
				eid = t.Event.ID
			}
			res, err = tx.tx.Exec(
				`INSERT INTO ticket (id, order_line, event, used) VALUES (?,?,?,?)`,
				ID(t.ID), ol.ID, IDStr(eid), Time(t.Used))
			panicOnError(err)
			if t.ID == 0 {
				t.ID = model.TicketID(lastInsertID(res))
			}
		}
		for _, pl := range ol.PaymentLines {
			panicOnNoRows(tx.tx.Exec(
				`INSERT OR REPLACE INTO payment_line (payment, order_line, amount) VALUES (?,?,?)`,
				pl.Payment.ID, ol.ID, pl.Amount))
		}
	}
}
