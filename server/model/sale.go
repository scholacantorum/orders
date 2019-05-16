package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"scholacantorum.org/orders/db"
)

// SaleID is the unique identifier of a sale.
type SaleID int

// Sale represents a sale to a customer.  See db/schema.sql for details.
type Sale struct {
	ID        SaleID      `json:"id"`
	StripeID  string      `json:"stripeID"`
	Customer  CustomerID  `json:"customer"`
	Source    string      `json:"source"`
	Timestamp time.Time   `json:"timestamp"`
	Payment   string      `json:"payment"`
	Note      string      `json:"note"`
	Lines     []*SaleLine `json:"lines"`
}

// saleColumns is the list of columns of the sale table.
var saleColumns = `id, stripeID, customer, source, timestamp, payment, note`

// scanSale scans a sale table row.
func scanSale(scanner interface{ Scan(...interface{}) error }, s *Sale) error {
	var (
		stripeID sql.NullString
		err      error
	)
	err = scanner.Scan(&s.ID, &stripeID, &s.Customer, &s.Source, (*db.Time)(&s.Timestamp), &s.Payment, &s.Note)
	if err != nil {
		return err
	}
	s.StripeID = stripeID.String
	return nil
}

// Save saves a sale to the database.  Note that it does not save the Lines
// array; sale lines have to be saved individually.
func (s *Sale) Save(tx *sql.Tx) {
	var (
		q        strings.Builder
		res      sql.Result
		nid      int64
		err      error
		stripeID = sql.NullString{String: s.StripeID, Valid: s.StripeID != ""}
	)
	q.WriteString(`INSERT OR REPLACE INTO sale (`)
	q.WriteString(saleColumns)
	q.WriteString(`) VALUES (?,?,?,?,?,?,?)`)
	res, err = tx.Exec(q.String(), db.ID(s.ID), stripeID, s.Customer, s.Source, db.Time(s.Timestamp), s.Payment, s.Note)
	panicOnError(err)
	if s.ID == 0 {
		if nid, err = res.LastInsertId(); err != nil {
			panic(err)
		}
		s.ID = SaleID(nid)
	}
}

// Delete deletes a sale, including all of its lines.
func (s *Sale) Delete(tx *sql.Tx) {
	panicOnNoRows(tx.Exec(`DELETE FROM sale WHERE id=?`, s.ID))
}

// FetchSale returns the sale with the specified ID.  It returns nil if no
// such sale exists.  It fills in the Lines array if lines is true.
func FetchSale(tx *sql.Tx, id SaleID, lines bool) (s *Sale) {
	var (
		q   strings.Builder
		err error
	)
	s = new(Sale)
	q.WriteString(`SELECT `)
	q.WriteString(saleColumns)
	q.WriteString(` FROM sale WHERE id=?`)
	switch err = scanSale(tx.QueryRow(q.String(), id), s); err {
	case nil:
		if lines {
			s.Lines = FetchSaleLinesForSale(tx, s.ID)
		}
		return s
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}
}

// FetchSales returns the list of sales meeting the specified criteria.  Sales
// are returned in chronological order by timestamp.  Offset and limit, when
// nonzero, limit the calls to the specified page of results.  If lines is true,
// the Lines array is populated for each returned sale.
func FetchSales(tx *sql.Tx, offset, limit int, lines bool, criteria string, args ...interface{}) (sales []*Sale) {
	var (
		q    strings.Builder
		rows *sql.Rows
		err  error
	)
	q.WriteString(`SELECT `)
	q.WriteString(saleColumns)
	q.WriteString(` FROM sale`)
	if criteria != "" {
		q.WriteString(` WHERE `)
		q.WriteString(criteria)
	}
	q.WriteString(` ORDER BY timestamp`)
	if limit != 0 {
		fmt.Fprintf(&q, ` LIMIT %d`, limit)
	}
	if offset != 0 {
		fmt.Fprintf(&q, ` OFFSET %d`, offset)
	}
	rows, err = tx.Query(q.String(), args...)
	panicOnError(err)
	for rows.Next() {
		var s Sale
		panicOnError(scanSale(rows, &s))
		sales = append(sales, &s)
	}
	panicOnError(rows.Err())
	if lines {
		for _, s := range sales {
			s.Lines = FetchSaleLinesForSale(tx, s.ID)
		}
	}
	return sales
}
