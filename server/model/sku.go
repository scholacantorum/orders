package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"scholacantorum.org/orders/db"
)

// SKUID is the unique identifier of a SKU.
type SKUID int

// SKU represents a SKU to which tickets are sold.  See db/schema.sql for
// details.
type SKU struct {
	ID          SKUID     `json:"id"`
	StripeID    string    `json:"stripeID"`
	Product     ProductID `json:"product"`
	Coupon      string    `json:"coupon"`
	SalesStart  time.Time `json:"salesStart"`
	SalesEnd    time.Time `json:"salesEnd"`
	MembersOnly bool      `json:"membersOnly"`
	Price       int       `json:"price"`
}

// skuColumns is the list of columns of the sku table.
var skuColumns = `id, stripeID, product, coupon, sales_start, sales_end, members_only, price`

// scanSKU scans a SKU table row.
func scanSKU(scanner interface{ Scan(...interface{}) error }, s *SKU) error {
	return scanner.Scan(&s.ID, &s.StripeID, &s.Product, &s.Coupon,
		(*db.Time)(&s.SalesStart), (*db.Time)(&s.SalesEnd),
		&s.MembersOnly, &s.Price)
}

// Save saves a SKU to the database.
func (s *SKU) Save(tx *sql.Tx) {
	var (
		q   strings.Builder
		res sql.Result
		nid int64
		err error
	)
	q.WriteString(`INSERT OR REPLACE INTO sku (`)
	q.WriteString(skuColumns)
	q.WriteString(`) VALUES (?,?,?,?,?,?,?,?)`)
	res, err = tx.Exec(q.String(), db.ID(s.ID), s.StripeID, s.Product,
		s.Coupon, db.Time(s.SalesStart), db.Time(s.SalesEnd),
		s.MembersOnly, s.Price)
	panicOnError(err)
	if s.ID == 0 {
		if nid, err = res.LastInsertId(); err != nil {
			panic(err)
		}
		s.ID = SKUID(nid)
	}
}

// Delete deletes a SKU.
func (s *SKU) Delete(tx *sql.Tx) {
	panicOnNoRows(tx.Exec(`DELETE FROM sku WHERE id=?`, s.ID))
}

// FetchSKU returns the sku with the specified ID.  It returns nil if no
// such sku exists.
func FetchSKU(tx *sql.Tx, id SKUID) (s *SKU) {
	var (
		q   strings.Builder
		err error
	)
	s = new(SKU)
	q.WriteString(`SELECT `)
	q.WriteString(skuColumns)
	q.WriteString(` FROM sku WHERE id=?`)
	switch err = scanSKU(tx.QueryRow(q.String(), id), s); err {
	case nil:
		return s
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}
}

// FetchSKUs returns the list of SKUs meeting the specified criteria.  SKUs are
// returned in ID order.  Offset and limit, when nonzero, limit the calls to the
// specified page of results.
func FetchSKUs(tx *sql.Tx, offset, limit int, criteria string, args ...interface{}) (skus []*SKU) {
	var (
		q    strings.Builder
		rows *sql.Rows
		err  error
	)
	q.WriteString(`SELECT `)
	q.WriteString(skuColumns)
	q.WriteString(` FROM sku`)
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
		var s SKU
		panicOnError(scanSKU(rows, &s))
		skus = append(skus, &s)
	}
	panicOnError(rows.Err())
	return skus
}

// FetchSKUsForProduct returns the list of SKUs for the product with the
// specified ID.
func FetchSKUsForProduct(tx *sql.Tx, id ProductID) []*SKU {
	return FetchSKUs(tx, 0, 0, `product=?`, id)
}
