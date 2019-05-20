// Package model contains the data model types and constants for the Schola
// Cantorum ordering system.
package model

import (
	"time"
)

type (
	Event struct {
		ID        EventID   `json:"id"`
		MembersID int       `json:"membersID"`
		Name      string    `json:"name"`
		Start     time.Time `json:"start"`
		Capacity  int       `json:"capacity"`
	}
	EventID string
	Order   struct {
		ID       OrderID
		Name     string
		Email    string
		Address  string
		City     string
		State    string
		Zip      string
		Phone    string
		Customer string
		Member   int
		Created  time.Time
		Flags    OrderFlags
		Note     string
		Coupon   string
		Repeat   time.Time
		Lines    []*OrderLine
		Payments []*Payment
	}
	OrderFlags byte
	OrderID    int
	OrderLine  struct {
		ID           OrderLineID
		Product      *Product
		Quantity     int
		Token        string
		Price        int
		Tickets      []*Ticket
		PaymentLines []*PaymentLine
	}
	OrderLineID int
	Payment     struct {
		ID      PaymentID
		Method  string
		Stripe  string
		Created time.Time
		Flags   PaymentFlags
		Lines   []*PaymentLine
	}
	PaymentFlags byte
	PaymentID    int
	PaymentLine  struct {
		Payment *Payment
		Amount  int
	}
	Privilege uint8
	Product   struct {
		ID          ProductID   `json:"id"`
		Name        string      `json:"name"`
		Type        ProductType `json:"type"`
		Receipt     string      `json:"receipt"`
		TicketCount int         `json:"ticketCount"`
		TicketClass string      `json:"ticketClass"`
		SKUs        []*SKU      `json:"skus"`
		Events      []*Event    `json:"events"`
	}
	ProductID   string
	ProductType string
	Session     struct {
		Token      string
		Username   string
		Expires    time.Time
		Member     int
		Privileges Privilege
	}
	SKU struct {
		Coupon      string    `json:"coupon"`
		SalesStart  time.Time `json:"salesStart"`
		SalesEnd    time.Time `json:"salesEnd"`
		MembersOnly bool      `json:"membersOnly"`
		Price       int       `json:"price"`
	}
	Ticket struct {
		ID    TicketID
		Event *Event
		Used  time.Time
	}
	TicketID int
)

const (
	// PrivSetup is the privilege needed to create, modify, or delete
	// products, SKUs, and events.
	PrivSetup Privilege = 1 << iota

	// PrivAnalyze allows read-only access to all data in the system.
	PrivAnalyze

	// PrivHandleOrders allows recording offline orders and making notes on
	// orders.
	PrivHandleOrders

	// PrivSell allows making and recording in-person sales.  Note that this
	// includes refunding those sales within a few minutes of recording
	// them.
	PrivSell

	// PrivAdmit allows recording ticket usage.
	PrivAdmit
)
