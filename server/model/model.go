// Package model contains the data model types and constants for the Schola
// Cantorum ordering system.
package model

import (
	"time"
)

type EventID string

type Event struct {
	ID        EventID
	MembersID int
	Name      string
	Start     time.Time
	Capacity  int
}

type OrderID int

type OrderFlags byte

const (
	// OrderValid indicates that the order is valid.  If this flag is not
	// set, the payment processing for the order is still in progress (or
	// failed), and the order should not be considered "real".
	OrderValid OrderFlags = 1 << iota
)

type OrderSource string

const (
	// OrderFromPublic indicates that this order was placed through Schola's
	// public web site, scholacantorum.org.
	OrderFromPublic OrderSource = "public"

	// OrderFromMembers indicates that this order was placed through
	// Schola's members web site, scholacantorummembers.org.
	OrderFromMembers = "members"

	// OrderFromGala indicates that this order was placed through Schola's
	// gala software, gala.scholacantorummembers.org.
	OrderFromGala = "gala"

	// OrderFromOffice indicates that this order was placed through Schola's
	// order management site, orders.scholacantorum.org.
	OrderFromOffice = "office"

	// OrderInPerson indicates that this order was placed through Schola's
	// in-person sales app.
	OrderInPerson = "inperson"
)

type Order struct {
	ID       OrderID
	Token    string
	Source   OrderSource
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
	CNote    string
	ONote    string
	Coupon   string
	Repeat   time.Time
	Lines    []*OrderLine
	Payments []*Payment
}

type OrderLineID int

type OrderLine struct {
	ID       OrderLineID
	Product  *Product
	Quantity int
	Used     int
	UsedAt   EventID
	Amount   int
	Tickets  []*Ticket
}

type PaymentID int

type PaymentFlags byte

type PaymentType string

const (
	// PaymentCard is a card-not-present Stripe card payment.
	PaymentCard PaymentType = "card"

	// PaymentCardPresent is a card-present Stripe card payment.
	PaymentCardPresent = "card-present"

	// PaymentOther is a non-Stripe payment, described in Method.
	PaymentOther = "other"
)

type Payment struct {
	ID      PaymentID
	Type    PaymentType
	Method  string
	Stripe  string
	Created time.Time
	Flags   PaymentFlags
	Amount  int
}

type Privilege uint8

const (
	// PrivLogin is a privilege held by all valid users.
	PrivLogin Privilege = 1 << iota

	// PrivSetup is the privilege needed to create, modify, or delete
	// products, SKUs, and events.
	PrivSetup

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

type ProductID string

type ProductType string

const (
	// ProdTicket is an event ticket.  It may be valid at one event or at
	// multiple events; it may be valid for a single entry or multiple
	// entries.
	ProdTicket ProductType = "ticket"

	// ProdRecording is a concert recording.  Recordings are available for
	// sale only to performers in that concert.
	ProdRecording = "recording"

	// ProdDonation is a donation.
	ProdDonation = "donation"

	// ProdSheetMusic is a set of sheet music sold to a singer in the
	// chorus.
	ProdSheetMusic = "sheetmusic"

	// ProdAuctionItem is an auction item purchased at the gala or similar.
	ProdAuctionItem = "auctionitem"
)

type Product struct {
	ID          ProductID
	Name        string
	Type        ProductType
	Receipt     string
	TicketCount int
	TicketClass string
	SKUs        []*SKU
	Events      []ProductEvent
}

type ProductEvent struct {
	Priority int
	Event    *Event
}

type Session struct {
	Token      string
	Username   string
	Expires    time.Time
	Member     int
	Privileges Privilege
}

type SKU struct {
	Coupon      string
	SalesStart  time.Time
	SalesEnd    time.Time
	MembersOnly bool
	Quantity    int
	Price       int
}

type TicketID int

type Ticket struct {
	ID    TicketID
	Event *Event
	Used  time.Time
}
