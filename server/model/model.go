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
	Series    string
	Start     time.Time
	Capacity  int
}

type OrderID int

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
	Valid    bool
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
	CNote    string
	ONote    string
	InAccess bool
	Coupon   string
	Lines    []*OrderLine
	Payments []*Payment
}

type OrderLineID int

type OrderLine struct {
	ID       OrderLineID
	Product  *Product
	Quantity int
	Price    int
	Scan     string
	MinUsed  int
	AutoUse  int
	Tickets  []*Ticket
	Used     int     // not persistent; input only
	UsedAt   EventID // not persistent; input only
	Error    string  // not persistent; output only
}

type PaymentID int

type PaymentFlags byte

const (
	// PaymentInitial marks the initial payment on an order, for ease of
	// queries.
	PaymentInitial PaymentFlags = 1 << iota
)

type PaymentType string

const (
	// PaymentCard is a card-not-present, immediate, single-use Stripe card
	// payment.
	PaymentCard PaymentType = "card"

	// PaymentCardPresent is a card-present Stripe card payment.
	PaymentCardPresent = "card-present"

	// PaymentCash is a cash payment.
	PaymentCash = "cash"

	// PaymentCheck is a check payment, with the check number usually given
	// in Method.
	PaymentCheck = "check"

	// PaymentOther is a nonstandard payment type, described in Method.
	PaymentOther = "other"
)

type Payment struct {
	ID      PaymentID
	Type    PaymentType
	Subtype string
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

	// PrivSetupOrders is the privilege needed to create, modify, or delete
	// products, SKUs, and events.
	PrivSetupOrders

	// PrivViewOrders allows read-only access to all data in the system.
	PrivViewOrders

	// PrivManageOrders allows recording offline orders and making notes on
	// orders.
	PrivManageOrders

	// PrivInPersonSales allows making and recording in-person sales.  Note
	// that this includes refunding those sales within a few minutes of
	// recording them.
	PrivInPersonSales

	// PrivScanTickets allows recording ticket usage.
	PrivScanTickets
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

	// ProdWardrobe is a concert dress or other wardrobe item sold to a
	// singer.
	ProdWardrobe = "wardrobe"

	// ProdOther is an "other" product type, not allowed for new products
	// but used for products in archive orders.
	ProdOther = "other"
)

type Product struct {
	ID          ProductID
	Series      string
	Name        string
	ShortName   string
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
	Name       string
	Email      string
	Address    string
	City       string
	State      string
	Zip        string
}

type SKU struct {
	Source     OrderSource
	Coupon     string
	SalesStart time.Time
	SalesEnd   time.Time
	Flags      SKUFlags
	Price      int
}

// InSalesRange returns -1 if the specified time is before the sales range of
// the SKU, 0 if it is in range, and +1 if it is after the sales range.
func (s *SKU) InSalesRange(t time.Time) int {
	if s.SalesStart.After(t) {
		return -1
	}
	if !s.SalesEnd.IsZero() && s.SalesEnd.Before(t) {
		return +1
	}
	return 0
}

type SKUFlags byte

type TicketID int

type Ticket struct {
	ID    TicketID
	Event *Event
	Used  time.Time
}
