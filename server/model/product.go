package model

import (
	"database/sql"
	"fmt"
	"strings"

	"scholacantorum.org/orders/db"
)

// ProductID is the unique identifier of a product.
type ProductID int

// Product represents a product that Schola sells.  See db/schema.sql for
// details.
type Product struct {
	ID          ProductID `json:"id"`
	StripeID    string    `json:"stripeID"`
	Name        string    `json:"name"`
	TicketCount int       `json:"ticketCount"`
	TicketClass string    `json:"ticketClass"`
	Events      []EventID `json:"events"`
}

// productColumns is the list of columns of the product table.
var productColumns = `id, stripeID, name, ticket_count, ticket_class`

// scanProduct scans a product table row.
func scanProduct(scanner interface{ Scan(...interface{}) error }, p *Product) error {
	return scanner.Scan(&p.ID, &p.StripeID, &p.Name, &p.TicketCount, &p.TicketClass)
}

// Save saves a product to the database.
func (p *Product) Save(tx *sql.Tx) {
	var (
		q   strings.Builder
		res sql.Result
		nid int64
		err error
	)
	q.WriteString(`INSERT OR REPLACE INTO product (`)
	q.WriteString(productColumns)
	q.WriteString(`) VALUES (?,?,?,?,?)`)
	res, err = tx.Exec(q.String(), db.ID(p.ID), p.StripeID, p.Name, p.TicketCount, p.TicketClass)
	panicOnError(err)
	if p.ID == 0 {
		if nid, err = res.LastInsertId(); err != nil {
			panic(err)
		}
		p.ID = ProductID(nid)
	}
	panicOnExecError(tx.Exec(`DELETE FROM product_event WHERE product=?`, p.ID))
	for _, eid := range p.Events {
		panicOnNoRows(tx.Exec(`INSERT INTO product_event (product, event) VALUES (?,?)`, p.ID, eid))
	}
}

// Delete deletes a product.
func (p *Product) Delete(tx *sql.Tx) {
	panicOnNoRows(tx.Exec(`DELETE FROM product WHERE id=?`, p.ID))
}

// FetchProduct returns the product with the specified ID.  It returns nil if no
// such product exists.
func FetchProduct(tx *sql.Tx, id ProductID) (p *Product) {
	var (
		q    strings.Builder
		rows *sql.Rows
		eid  EventID
		err  error
	)
	p = new(Product)
	q.WriteString(`SELECT `)
	q.WriteString(productColumns)
	q.WriteString(` FROM product WHERE id=?`)
	switch err = scanProduct(tx.QueryRow(q.String(), id), p); err {
	case nil:
		break
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}
	rows, err = tx.Query(`SELECT event FROM product_event WHERE product=?`, p.ID)
	panicOnError(err)
	for rows.Next() {
		panicOnError(rows.Scan(&eid))
		p.Events = append(p.Events, eid)
	}
	panicOnError(rows.Err())
	return p
}

// FetchProducts returns the list of products meeting the specified criteria.
// Products are returned in ID order.  Offset and limit, when nonzero, limit the
// calls to the specified page of results.
func FetchProducts(tx *sql.Tx, offset, limit int, criteria string, args ...interface{}) (products []*Product) {
	var (
		q     strings.Builder
		prows *sql.Rows
		erows *sql.Rows
		eid   EventID
		err   error
	)
	q.WriteString(`SELECT `)
	q.WriteString(productColumns)
	q.WriteString(` FROM product`)
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
	prows, err = tx.Query(q.String(), args...)
	panicOnError(err)
	for prows.Next() {
		var p Product
		panicOnError(scanProduct(prows, &p))
		erows, err = tx.Query(`SELECT event FROM product_event WHERE product=?`, p.ID)
		panicOnError(err)
		for erows.Next() {
			panicOnError(erows.Scan(&eid))
			p.Events = append(p.Events, eid)
		}
		panicOnError(erows.Err())
		products = append(products, &p)
	}
	panicOnError(prows.Err())
	return products
}
