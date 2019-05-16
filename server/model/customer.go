package model

import (
	"database/sql"
	"fmt"
	"strings"

	"scholacantorum.org/orders/db"
)

// CustomerID is the unique identifier of a customer.
type CustomerID int

// Customer represents a customer to whom Schola sells things.  See
// db/schema.sql for details.
type Customer struct {
	ID       CustomerID `json:"id"`
	StripeID string     `json:"stripeID"`
	MemberID int        `json:"memberID"`
	Name     string     `json:"name"`
	Email    string     `json:"email"`
	Address  string     `json:"address"`
	City     string     `json:"city"`
	State    string     `json:"state"`
	Zip      string     `json:"zip"`
	Phone    string     `json:"phone"`
}

// customerColumns is the list of columns of the customer table.
var customerColumns = `id, stripeID, memberID, name, email, address, city, state, zip, phone`

// scanCustomer scans a customer table row.
func scanCustomer(scanner interface{ Scan(...interface{}) error }, c *Customer) error {
	var (
		stripeID sql.NullString
		err      error
	)
	err = scanner.Scan(&c.ID, &stripeID, &c.MemberID, &c.Name, &c.Email, &c.Address, &c.City, &c.State, &c.Zip, &c.Phone)
	if err != nil {
		return err
	}
	c.StripeID = stripeID.String
	return nil
}

// Save saves a customer to the database.
func (c *Customer) Save(tx *sql.Tx) {
	var (
		q        strings.Builder
		res      sql.Result
		nid      int64
		err      error
		stripeID = sql.NullString{String: c.StripeID, Valid: c.StripeID != ""}
	)
	q.WriteString(`INSERT OR REPLACE INTO customer (`)
	q.WriteString(customerColumns)
	q.WriteString(`) VALUES (?,?,?,?,?,?,?,?,?,?)`)
	res, err = tx.Exec(q.String(), db.ID(c.ID), stripeID, c.MemberID,
		c.Name, c.Email, c.Address, c.City, c.State, c.Zip, c.Phone)
	panicOnError(err)
	if c.ID == 0 {
		if nid, err = res.LastInsertId(); err != nil {
			panic(err)
		}
		c.ID = CustomerID(nid)
	}
}

// Delete deletes a customer.
func (c *Customer) Delete(tx *sql.Tx) {
	panicOnNoRows(tx.Exec(`DELETE FROM customer WHERE id=?`, c.ID))
}

// FetchCustomer returns the customer with the specified ID.  It returns nil if
// no such customer exists.
func FetchCustomer(tx *sql.Tx, id CustomerID) (c *Customer) {
	var (
		q   strings.Builder
		err error
	)
	c = new(Customer)
	q.WriteString(`SELECT `)
	q.WriteString(customerColumns)
	q.WriteString(` FROM customer WHERE id=?`)
	switch err = scanCustomer(tx.QueryRow(q.String(), id), c); err {
	case nil:
		return c
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}
}

// FetchCustomerByNameEmailMemberID returns the customer matching the specified
// name, email address, and member ID value.  (Note that name and email are case
// insensitive.)  See db/schema.sql for details of the matching algorithm.  It
// returns nil if no such customer exists.
func FetchCustomerByNameEmailMemberID(tx *sql.Tx, name, email string, memberID int) (c *Customer) {
	var (
		q    strings.Builder
		args []interface{}
		err  error
	)
	c = new(Customer)
	q.WriteString(`SELECT `)
	q.WriteString(customerColumns)
	q.WriteString(` FROM customer`)
	switch {
	case name != "" && email != "":
		q.WriteString(` WHERE ((name=?1 AND email=?2) OR (name='' AND email=?2) OR (name=?1 AND email='')) AND memberID=?3`)
		args = []interface{}{name, email, memberID}
	case name != "":
		q.WriteString(` WHERE name=? AND memberID=?`)
		args = []interface{}{name, memberID}
	case email != "":
		q.WriteString(` WHERE email=? AND memberID=?`)
		args = []interface{}{email, memberID}
	default:
		q.WriteString(` WHERE name='' AND email=''`)
	}
	switch err = scanCustomer(tx.QueryRow(q.String(), args...), c); err {
	case nil:
		return c
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}
}

// FetchCustomers returns the set of customers meeting the specified criteria.
// Customers are returned in alphabetical order by name.  Offset and limit, when
// nonzero, limit the calls to the specified page of results.
func FetchCustomers(tx *sql.Tx, offset, limit int, criteria string, args ...interface{}) (list []*Customer) {
	var (
		q    strings.Builder
		rows *sql.Rows
		err  error
	)
	q.WriteString(`SELECT `)
	q.WriteString(customerColumns)
	q.WriteString(` FROM customer`)
	if criteria != "" {
		q.WriteString(` WHERE `)
		q.WriteString(criteria)
	}
	q.WriteString(` ORDER BY name`)
	if limit != 0 {
		fmt.Fprintf(&q, ` LIMIT %d`, limit)
	}
	if offset != 0 {
		fmt.Fprintf(&q, ` OFFSET %d`, offset)
	}
	rows, err = tx.Query(q.String(), args...)
	panicOnError(err)
	for rows.Next() {
		var c Customer
		panicOnError(scanCustomer(rows, &c))
		list = append(list, &c)
	}
	panicOnError(rows.Err())
	return list
}
