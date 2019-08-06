// Package model contains the data model types and constants for the Schola
// Cantorum ordering system.
package model

import (
	"time"
)

// A ReportDefinition describes what should be in a report.
type ReportDefinition struct {

	// OrderSources is a list of order sources to be included in the report.
	// An empty list means all order sources.
	OrderSources []OrderSource

	// Customer specifies the customer(s) whose tickets should be included
	// in the report.  It is a case-insensitive substring matched against
	// customer name or email address.  An empty string means all customers.
	Customer string

	// CreatedBefore specifies the upper limit of the range of order
	// creation timestamps of orders to include in the report.  If it is
	// the zero time, there is no upper limit.
	CreatedBefore time.Time

	// CreatedAfter specifies the lower limit of the range of order creation
	// timestamps of orders to include in the report.  If it is the zero
	// time, there is no lower limit.
	CreatedAfter time.Time

	// OrderCoupons is a list of order coupons to be included in the report.
	// An empty list means all orders regardless of coupon.  Including an
	// empty string in the list includes orders with no coupon.
	OrderCoupons []string

	// Products is a list of products to be included in the report.  An
	// empty list includes all products.
	Products []ProductID

	// PaymentTypes is a list of payment types (for the initial payment of
	// an order) to be included in the report.  Each payment type is a
	// string of one or two parts, separated by a comma.  An empty list
	// includes all orders regardless of payment type.
	PaymentTypes []string

	// TicketClasses is a list of ticket classes to be included in the
	// report.  An empty list means all ticket classes.
	TicketClasses []string

	// UsedAtEvents is a list of events at which the tickets must have been
	// used in order to be included in the report.  An empty list includes
	// all orders regardless of ticket usage.  An empty string on the list
	// includes orders with unused tickets.
	UsedAtEvents []EventID
}

// ReportResults contains the results of running a report.
type ReportResults struct {

	// OrderCount gives the number of orders matching the report criteria.
	OrderCount int

	// ItemCount gives the number of purchased items matching the report
	// criteria.
	ItemCount int

	// TotalAmount gives the sum of the amounts (in dollars) of the lines
	// matching the report criteria.  Note that this may not be a whole
	// dollar.
	TotalAmount float64

	// Lines gives the matching report lines.  It is nil if no report
	// criteria were given, or if the criteria match too many lines.  It is
	// an empty slice if no purchases match the report criteria.
	Lines []*ReportLine

	// OrderSources gives, for each order source, the number of results that
	// would match all of the other report criteria if that OrderSource were
	// the only one selected.
	OrderSources map[OrderSource]int

	// OrderCoupons gives, for each coupon code, the number of results that
	// would match all of the other report criteria if that OrderCoupon were
	// the only one selected.
	OrderCoupons StringCounts

	// Products gives, for each product, the number of results that would
	// match all of the other report criteria if that product were the only
	// one selected.
	Products []*ReportProductCount

	// PaymentTypes gives, for each payment type, the number of results that
	// would match all of the other report criteria if that payment type
	// were the only one selected.
	PaymentTypes StringCounts

	// TicketClasses gives, for each ticket class, the number of results
	// that would match all of the other report criteria if that TicketClass
	// were the only one selected.
	TicketClasses StringCounts

	// UsedAtEvents gives, for each event, the number of results that would
	// match all of the other report criteria if that UsedAtEvent were the
	// only one selected.  There will be an entry for "" representing unused
	// tickets.
	UsedAtEvents []*ReportEventCount
}

// A ReportLine is one line in a report.
type ReportLine struct {
	OrderID     OrderID
	OrderTime   time.Time
	Name        string
	Email       string
	Quantity    int
	Product     string
	UsedAtEvent EventID
	OrderSource OrderSource
	PaymentType string
	Amount      float64
}

// A ReportProductCount provides the statistical and hierarchical information
// for one product in a report.
type ReportProductCount struct {
	ID     ProductID
	Name   string
	Series string
	Type   ProductType
	Count  int
}

// A ReportEventCount provides the statistical and hierarchical information for
// one event (at which tickets were used) in a report.
type ReportEventCount struct {
	ID     EventID
	Start  time.Time
	Name   string
	Series string
	Count  int
}

// StringCounts is a map from string to integer, with associated methods for
// sorted access.
type StringCounts map[string]int
