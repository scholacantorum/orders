package model

import (
	"database/sql"
	"fmt"
	"strings"

	"scholacantorum.org/orders/db"
)

// SaleLineID is the unique identifier of a sale line.
type SaleLineID int

// SaleLine represents one line of a sale.  See db/schema.sql for details.
type SaleLine struct {
	ID     SaleLineID `json:"id"`
	Sale   SaleID     `json:"sale"`
	SKU    SKUID      `json:"sku"`
	Qty    int        `json:"qty"`
	Amount int        `json:"amount"`
}

// saleLineColumns is the list of columns of the sale_line table.
var saleLineColumns = `id, sale, sku, qty, amount`

// scanSaleLine scans a sale line table row.
func scanSaleLine(scanner interface{ Scan(...interface{}) error }, s *SaleLine) error {
	return scanner.Scan(&s.ID, &s.Sale, &s.SKU, &s.Qty, &s.Amount)
}

// Save saves a sale line to the database.
func (s *SaleLine) Save(tx *sql.Tx) {
	var (
		q   strings.Builder
		res sql.Result
		nid int64
		err error
	)
	q.WriteString(`INSERT OR REPLACE INTO sale_line (`)
	q.WriteString(saleLineColumns)
	q.WriteString(`) VALUES (?,?,?,?,?)`)
	res, err = tx.Exec(q.String(), db.ID(s.ID), s.Sale, s.SKU, s.Qty, s.Amount)
	panicOnError(err)
	if s.ID == 0 {
		if nid, err = res.LastInsertId(); err != nil {
			panic(err)
		}
		s.ID = SaleLineID(nid)
	}
}

// Delete deletes a sale line.
func (s *SaleLine) Delete(tx *sql.Tx) {
	panicOnNoRows(tx.Exec(`DELETE FROM sale_line WHERE id=?`, s.ID))
}

// FetchSaleLine returns the sale line with the specified ID.  It returns nil if
// no such sale line exists.
func FetchSaleLine(tx *sql.Tx, id SaleLineID) (s *SaleLine) {
	var (
		q   strings.Builder
		err error
	)
	s = new(SaleLine)
	q.WriteString(`SELECT `)
	q.WriteString(saleLineColumns)
	q.WriteString(` FROM sale_line WHERE id=?`)
	switch err = scanSaleLine(tx.QueryRow(q.String(), id), s); err {
	case nil:
		return s
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}
}

// FetchSaleLines calls the supplied function for each sale line meeting the
// specified criteria.  SaleLines are returned in ID order.  Offset and limit,
// when nonzero, limit the calls to the specified page of results.
func FetchSaleLines(tx *sql.Tx, offset, limit int, fn func(*SaleLine), criteria string, args ...interface{}) {
	var (
		q    strings.Builder
		s    SaleLine
		rows *sql.Rows
		err  error
	)
	q.WriteString(`SELECT `)
	q.WriteString(saleLineColumns)
	q.WriteString(` FROM sale_line`)
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
		panicOnError(scanSaleLine(rows, &s))
		fn(&s)
	}
	panicOnError(rows.Err())
}

// FetchSaleLinesForSale returns the list of sale lines for the sale with the
// specified ID.
func FetchSaleLinesForSale(tx *sql.Tx, id SaleID) (lines []*SaleLine) {
	FetchSaleLines(tx, 0, 0, func(sl *SaleLine) {
		var copy = *sl
		lines = append(lines, &copy)
	}, `sale=?`, id)
	return lines
}
