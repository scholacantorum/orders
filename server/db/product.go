package db

import (
	"database/sql"
	"strings"

	"scholacantorum.org/orders/model"
)

// productColumns is the list of columns of the product table.
var productColumns = `id, name, shortname, type, receipt, ticket_count, ticket_class`

// scanProduct scans a product table row.
func (tx Tx) scanProduct(scanner interface{ Scan(...interface{}) error }, p *model.Product) error {
	return scanner.Scan(&p.ID, &p.Name, &p.ShortName, &p.Type, &p.Receipt, &p.TicketCount, &p.TicketClass)
}

// SaveProduct saves a product to the database.
func (tx Tx) SaveProduct(p *model.Product) {
	var (
		q strings.Builder
	)
	q.WriteString(`INSERT OR REPLACE INTO product (`)
	q.WriteString(productColumns)
	q.WriteString(`) VALUES (?,?,?,?,?,?,?)`)
	panicOnExecError(tx.tx.Exec(q.String(), p.ID, p.Name, p.ShortName, p.Type, p.Receipt, p.TicketCount, p.TicketClass))
	panicOnExecError(tx.tx.Exec(`DELETE FROM product_event WHERE product=?`, p.ID))
	for _, pe := range p.Events {
		panicOnNoRows(tx.tx.Exec(
			`INSERT INTO product_event (product, event, priority) VALUES (?,?,?)`, p.ID, pe.Event.ID, pe.Priority))
	}
	panicOnExecError(tx.tx.Exec(`DELETE FROM sku WHERE product=?`, p.ID))
	for _, sku := range p.SKUs {
		panicOnExecError(tx.tx.Exec(
			`INSERT INTO sku (product, coupon, sales_start, sales_end, flags, price) VALUES (?,?,?,?,?,?)`,
			p.ID, sku.Coupon, Time(sku.SalesStart), Time(sku.SalesEnd), sku.Flags, sku.Price))
	}
}

// DeleteProduct deletes a product.
func (tx Tx) DeleteProduct(p *model.Product) {
	panicOnNoRows(tx.tx.Exec(`DELETE FROM product WHERE id=?`, p.ID))
}

// FetchProduct returns the product with the specified ID.  It returns nil if no
// such product exists.
func (tx Tx) FetchProduct(id model.ProductID) (p *model.Product) {
	var (
		q    strings.Builder
		rows *sql.Rows
		prio int
		eid  model.EventID
		err  error
	)
	p = new(model.Product)
	q.WriteString(`SELECT `)
	q.WriteString(productColumns)
	q.WriteString(` FROM product WHERE id=?`)
	switch err = tx.scanProduct(tx.tx.QueryRow(q.String(), id), p); err {
	case nil:
		break
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}
	rows, err = tx.tx.Query(`SELECT event, priority FROM product_event WHERE product=?`, p.ID)
	panicOnError(err)
	for rows.Next() {
		panicOnError(rows.Scan(&eid, &prio))
		var event = tx.FetchEvent(eid)
		p.Events = append(p.Events, model.ProductEvent{Event: event, Priority: prio})
	}
	panicOnError(rows.Err())
	rows, err = tx.tx.Query(
		`SELECT coupon, sales_start, sales_end, flags, price FROM sku WHERE product=?`, p.ID)
	panicOnError(err)
	for rows.Next() {
		var sku model.SKU
		panicOnError(rows.Scan(&sku.Coupon, (*Time)(&sku.SalesStart),
			(*Time)(&sku.SalesEnd), &sku.Flags, &sku.Price))
		p.SKUs = append(p.SKUs, &sku)
	}
	panicOnError(rows.Err())
	return p
}

// FetchProductsByEvent returns the set of products that give entry to the
// specified event.  They are returned in priority order.  Their Events lists
// are not complete; they contain only the single entry for the event requested.
func (tx Tx) FetchProductsByEvent(event *model.Event) (products []*model.Product) {
	var (
		prows    *sql.Rows
		srows    *sql.Rows
		erows    *sql.Rows
		eid      model.EventID
		priority int
		err      error
	)
	prows, err = tx.tx.Query(`
SELECT p.id, p.name, p.shortname, p.type, p.receipt, p.ticket_count, p.ticket_class
FROM product p, product_event pe WHERE pe.product=p.id AND pe.event=? ORDER BY pe.priority`, event.ID)
	panicOnError(err)
	for prows.Next() {
		var p model.Product
		panicOnError(tx.scanProduct(prows, &p))
		srows, err = tx.tx.Query(
			`SELECT coupon, sales_start, sales_end, flags, price FROM sku WHERE product=?`, p.ID)
		panicOnError(err)
		for srows.Next() {
			var sku model.SKU
			panicOnError(srows.Scan(&sku.Coupon, (*Time)(&sku.SalesStart),
				(*Time)(&sku.SalesEnd), &sku.Flags, &sku.Price))
			p.SKUs = append(p.SKUs, &sku)
		}
		panicOnError(srows.Err())
		erows, err = tx.tx.Query(`SELECT event, priority FROM product_event WHERE product=?`, p.ID)
		panicOnError(err)
		for erows.Next() {
			panicOnError(erows.Scan(&eid, &priority))
			var event = tx.FetchEvent(eid)
			p.Events = append(p.Events, model.ProductEvent{Event: event, Priority: priority})
		}
		panicOnError(erows.Err())
		products = append(products, &p)
	}
	panicOnError(prows.Err())
	return products
}
