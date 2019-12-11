package db

import (
	"database/sql"
	"strings"

	"scholacantorum.org/orders/model"
)

// orderColumns is the list of columns in the orderT table.
var orderColumns = `id, token, valid, source, name, email, address, city, state, zip, phone, customer, member, created, cnote, onote, in_access, coupon`

// scanOrder scans an orderT table row.
func scanOrder(scanner interface{ Scan(...interface{}) error }, o *model.Order) error {
	return scanner.Scan(&o.ID, &o.Token, &o.Valid, &o.Source, &o.Name, &o.Email,
		&o.Address, &o.City, &o.State, &o.Zip, &o.Phone, &o.Customer,
		&o.Member, (*Time)(&o.Created), &o.CNote, &o.ONote, &o.InAccess,
		&o.Coupon)
}

// FetchOrder returns the order with the specified ID.  It returns nil if no
// such order exists.  This includes fetching all subsidiary objects.
func (tx Tx) FetchOrder(id model.OrderID) (o *model.Order) {
	var (
		q     strings.Builder
		lrows *sql.Rows
		pid   model.ProductID
		prows *sql.Rows
		trows *sql.Rows
		eid   model.EventID
		err   error
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
	prows, err = tx.tx.Query(
		`SELECT id, type, subtype, method, stripe, created, initial, amount FROM payment WHERE orderid=? ORDER BY id`,
		o.ID)
	panicOnError(err)
	for prows.Next() {
		var p model.Payment
		panicOnError(prows.Scan(&p.ID, &p.Type, &p.Subtype, &p.Method, &p.Stripe, (*Time)(&p.Created), &p.Amount))
		o.Payments = append(o.Payments, &p)
	}
	panicOnError(prows.Err())
	lrows, err = tx.tx.Query(
		`SELECT id, product, quantity, price FROM order_line WHERE orderid=? ORDER BY id`, o.ID)
	panicOnError(err)
	for lrows.Next() {
		var ol model.OrderLine
		panicOnError(lrows.Scan(&ol.ID, &pid, &ol.Quantity, &ol.Price))
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
		o.Lines = append(o.Lines, &ol)
	}
	panicOnError(lrows.Err())
	return o
}

// FetchOrderByToken returns the order with the specified token.  It returns nil
// if no such order exists.
func (tx Tx) FetchOrderByToken(token string) *model.Order {
	var (
		oid model.OrderID
		err error
	)
	switch err = tx.tx.QueryRow(`SELECT id FROM orderT WHERE token=?`, token).Scan(&oid); err {
	case nil:
		break
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}
	return tx.FetchOrder(oid)
}

// SaveOrder saves an order to the database.  This includes saving all
// order-specific subsidiary objects.
func (tx Tx) SaveOrder(o *model.Order) {
	var (
		q   strings.Builder
		res sql.Result
		err error
	)
	q.WriteString(`INSERT OR REPLACE INTO orderT (`)
	q.WriteString(orderColumns)
	q.WriteString(`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	res, err = tx.tx.Exec(q.String(), ID(o.ID), o.Token, o.Valid, o.Source,
		o.Name, o.Email, o.Address, o.City, o.State, o.Zip, o.Phone,
		o.Customer, o.Member, Time(o.Created), o.CNote, o.ONote,
		o.InAccess, o.Coupon)
	panicOnError(err)
	if o.ID == 0 {
		o.ID = model.OrderID(lastInsertID(res))
	}
	for i, p := range o.Payments {
		res, err = tx.tx.Exec(
			`INSERT OR REPLACE INTO payment (id, orderid, type, subtype, method, stripe, created, initial, amount) VALUES (?,?,?,?,?,?,?,?,?)`,
			ID(p.ID), o.ID, p.Type, p.Subtype, p.Method, p.Stripe, Time(p.Created), i == 0, p.Amount)
		panicOnError(err)
		if p.ID == 0 {
			p.ID = model.PaymentID(lastInsertID(res))
		}
	}
	for _, ol := range o.Lines {
		res, err = tx.tx.Exec(
			`INSERT OR REPLACE INTO order_line (id, orderid, product, quantity, price) VALUES (?,?,?,?,?)`,
			ID(ol.ID), o.ID, ol.Product.ID, ol.Quantity, ol.Price)
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
	}
}

// DeleteOrder deletes an order from the database.  Generally this is done only
// if the order was not processed successfully.
func (tx Tx) DeleteOrder(o *model.Order) {
	panicOnExecError(tx.tx.Exec(`DELETE FROM payment WHERE orderid=?`, o.ID))
	panicOnNoRows(tx.tx.Exec(`DELETE FROM orderT WHERE id=?`, o.ID))
}
