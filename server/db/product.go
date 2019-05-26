package db

import (
	"database/sql"
	"strings"

	"scholacantorum.org/orders/model"
)

// productColumns is the list of columns of the product table.
var productColumns = `id, name, type, receipt, ticket_name, ticket_count, ticket_class`

// scanProduct scans a product table row.
func (tx Tx) scanProduct(scanner interface{ Scan(...interface{}) error }, p *model.Product) error {
	return scanner.Scan(&p.ID, &p.Name, &p.Type, &p.Receipt, &p.TicketName, &p.TicketCount, &p.TicketClass)
}

// SaveProduct saves a product to the database.
func (tx Tx) SaveProduct(p *model.Product) {
	var (
		q strings.Builder
	)
	q.WriteString(`INSERT OR REPLACE INTO product (`)
	q.WriteString(productColumns)
	q.WriteString(`) VALUES (?,?,?,?,?,?,?)`)
	panicOnExecError(tx.tx.Exec(q.String(), p.ID, p.Name, p.Type, p.Receipt, p.TicketName, p.TicketCount, p.TicketClass))
	panicOnExecError(tx.tx.Exec(`DELETE FROM product_event WHERE product=?`, p.ID))
	for _, event := range p.Events {
		panicOnNoRows(tx.tx.Exec(`INSERT INTO product_event (product, event) VALUES (?,?)`, p.ID, event.ID))
	}
	panicOnExecError(tx.tx.Exec(`DELETE FROM sku WHERE product=?`, p.ID))
	for _, sku := range p.SKUs {
		panicOnExecError(tx.tx.Exec(
			`INSERT INTO sku (product, coupon, sales_start, sales_end, members_only, quantity, price) VALUES (?,?,?,?,?,?,?)`,
			p.ID, sku.Coupon, Time(sku.SalesStart), Time(sku.SalesEnd), sku.MembersOnly, sku.Quantity, sku.Price))
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
	rows, err = tx.tx.Query(`SELECT event FROM product_event WHERE product=?`, p.ID)
	panicOnError(err)
	for rows.Next() {
		panicOnError(rows.Scan(&eid))
		var event = tx.FetchEvent(eid)
		p.Events = append(p.Events, event)
	}
	panicOnError(rows.Err())
	q.Reset()
	q.WriteString(`SELECT `)
	q.WriteString(`coupon, sales_start, sales_end, members_only, quantity, price`)
	q.WriteString(` FROM sku WHERE product=?`)
	rows, err = tx.tx.Query(q.String(), p.ID)
	panicOnError(err)
	for rows.Next() {
		var sku model.SKU
		panicOnError(rows.Scan(&sku.Coupon, (*Time)(&sku.SalesStart),
			(*Time)(&sku.SalesEnd), &sku.MembersOnly, &sku.Quantity,
			&sku.Price))
		p.SKUs = append(p.SKUs, &sku)
	}
	panicOnError(rows.Err())
	return p
}
