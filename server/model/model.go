// Package model contains the data model types and constants for the Schola
// Cantorum ordering system.
package model

import (
	"encoding/json"
	"errors"
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

func (e *Event) MarshalJSON() ([]byte, error) {
	var t = struct {
		ID        EventID `json:"id,omitempty"`
		MembersID int     `json:"membersID,omitempty"`
		Name      string  `json:"name,omitempty"`
		Start     string  `json:"start,omitempty"`
		Capacity  int     `json:"capacity,omitempty"`
	}{e.ID, e.MembersID, e.Name, fmtJSONTime(e.Start), e.Capacity}
	return json.Marshal(&t)
}

func (e *Event) UnmarshalJSON(data []byte) (err error) {
	var t struct {
		ID        EventID `json:"id"`
		MembersID int     `json:"membersID"`
		Name      string  `json:"name"`
		Start     string  `json:"start"`
		Capacity  int     `json:"capacity"`
	}
	if err = json.Unmarshal(data, &t); err != nil {
		return err
	}
	e.ID = t.ID
	e.MembersID = t.MembersID
	e.Name = t.Name
	if e.Start, err = parseJSONTime(t.Start); err != nil {
		return err
	}
	e.Capacity = t.Capacity
	return nil
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

func (o *Order) MarshalJSON() ([]byte, error) {
	var t = struct {
		ID       OrderID      `json:"id,omitempty"`
		Source   OrderSource  `json:"source,omitempty"`
		Name     string       `json:"name,omitempty"`
		Email    string       `json:"email,omitempty"`
		Address  string       `json:"address,omitempty"`
		City     string       `json:"city,omitempty"`
		State    string       `json:"state,omitempty"`
		Zip      string       `json:"zip,omitempty"`
		Phone    string       `json:"phone,omitempty"`
		Customer string       `json:"customer,omitempty"`
		Member   int          `json:"member,omitempty"`
		Created  string       `json:"created,omitempty"`
		CNote    string       `json:"cNote,omitempty"`
		ONote    string       `json:"oNote,omitempty"`
		Coupon   string       `json:"coupon,omitempty"`
		Repeat   string       `json:"repeat,omitempty"`
		Lines    []*OrderLine `json:"lines,omitempty"`
		Payments []*Payment   `json:"payments,omitempty"`
	}{o.ID, o.Source, o.Name, o.Email, o.Address, o.City, o.State, o.Zip,
		o.Phone, o.Customer, o.Member, fmtJSONTime(o.Created), o.CNote,
		o.ONote, o.Coupon, fmtJSONTime(o.Repeat), o.Lines, o.Payments}
	return json.Marshal(&t)
}

func (o *Order) UnmarshalJSON(data []byte) (err error) {
	var t struct {
		ID       OrderID      `json:"id"`
		Source   OrderSource  `json:"source"`
		Name     string       `json:"name"`
		Email    string       `json:"email"`
		Address  string       `json:"address"`
		City     string       `json:"city"`
		State    string       `json:"state"`
		Zip      string       `json:"zip"`
		Phone    string       `json:"phone"`
		Customer string       `json:"customer"`
		Member   int          `json:"member"`
		Created  string       `json:"created"`
		CNote    string       `json:"cNote"`
		ONote    string       `json:"oNote"`
		Coupon   string       `json:"coupon"`
		Repeat   string       `json:"repeat"`
		Lines    []*OrderLine `json:"lines"`
		Payments []*Payment   `json:"payments"`
	}
	if err = json.Unmarshal(data, &t); err != nil {
		return err
	}
	o.ID = t.ID
	o.Token = ""
	o.Source = t.Source
	o.Name = t.Name
	o.Email = t.Email
	o.Address = t.Address
	o.City = t.City
	o.State = t.State
	o.Zip = t.Zip
	o.Phone = t.Phone
	o.Customer = t.Customer
	o.Member = t.Member
	if o.Created, err = parseJSONTime(t.Created); err != nil {
		return err
	}
	o.Flags = 0
	o.CNote = t.CNote
	o.ONote = t.ONote
	o.Coupon = t.Coupon
	if o.Repeat, err = parseJSONTime(t.Repeat); err != nil {
		return err
	}
	o.Lines = t.Lines
	o.Payments = t.Payments
	return nil
}

type OrderLineID int

type OrderLine struct {
	ID       OrderLineID
	Product  *Product
	Quantity int
	Used     int
	UsedAt   EventID
	Price    int
	Tickets  []*Ticket
}

func (ol *OrderLine) MarshalJSON() ([]byte, error) {
	var t = struct {
		ID       OrderLineID `json:"id,omitempty"`
		Product  ProductID   `json:"product,omitempty"`
		Quantity int         `json:"quantity,omitempty"`
		Price    int         `json:"price"`
		Tickets  []*Ticket   `json:"tickets,omitempty"`
	}{ol.ID, ol.Product.ID, ol.Quantity, ol.Price, ol.Tickets}
	return json.Marshal(&t)
}

func (ol *OrderLine) UnmarshalJSON(data []byte) (err error) {
	var t struct {
		ID       OrderLineID `json:"id"`
		Product  ProductID   `json:"product"`
		Quantity int         `json:"quantity"`
		Used     int         `json:"used"`
		UsedAt   EventID     `json:"usedAt"`
		Price    int         `json:"price"`
		Tickets  []*Ticket   `json:"tickets"`
	}
	if err = json.Unmarshal(data, &t); err != nil {
		return err
	}
	ol.ID = t.ID
	ol.Product = &Product{ID: t.Product}
	ol.Quantity = t.Quantity
	ol.Used = t.Used
	ol.UsedAt = t.UsedAt
	ol.Price = t.Price
	ol.Tickets = t.Tickets
	return nil
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

func (p *Payment) MarshalJSON() ([]byte, error) {
	var t = struct {
		ID      PaymentID   `json:"id,omitempty"`
		Type    PaymentType `json:"type,omitempty"`
		Method  string      `json:"method,omitempty"`
		Stripe  string      `json:"stripe,omitempty"`
		Created string      `json:"created,omitempty"`
		Amount  int         `json:"amount"`
	}{p.ID, p.Type, p.Method, p.Stripe, fmtJSONTime(p.Created), p.Amount}
	return json.Marshal(&t)
}

func (p *Payment) UnmarshalJSON(data []byte) (err error) {
	var t struct {
		ID      PaymentID   `json:"id"`
		Type    PaymentType `json:"type"`
		Method  string      `json:"method"`
		Created string      `json:"created"`
		Amount  int         `json:"amount"`
	}
	if err = json.Unmarshal(data, &t); err != nil {
		return err
	}
	p.ID = t.ID
	p.Type = t.Type
	p.Method = t.Method
	p.Stripe = ""
	if p.Created, err = parseJSONTime(t.Created); err != nil {
		return err
	}
	p.Flags = 0
	p.Amount = t.Amount
	return nil
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
	// ProdTicket is an individual event ticket.
	ProdTicket ProductType = "ticket"

	// ProdFlexPass is a ticket that can be used a specified number of times
	// at a specific set of events.  All of its uses can be at one of the
	// events, or one at each of the events, or any mixture.
	ProdFlexPass = "flexpass"

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
	Events      []*Event
}

func (p *Product) MarshalJSON() ([]byte, error) {
	var t = struct {
		ID          ProductID   `json:"id,omitempty"`
		Name        string      `json:"name,omitempty"`
		Type        ProductType `json:"type,omitempty"`
		Receipt     string      `json:"receipt,omitempty"`
		TicketCount int         `json:"ticketCount,omitempty"`
		TicketClass string      `json:"ticketClass,omitempty"`
		SKUs        []*SKU      `json:"skus,omitempty"`
		Events      []EventID   `json:"events,omitempty"`
	}{p.ID, p.Name, p.Type, p.Receipt, p.TicketCount, p.TicketClass, p.SKUs, make([]EventID, len(p.Events))}
	for i, e := range p.Events {
		t.Events[i] = e.ID
	}
	return json.Marshal(&t)
}

func (p *Product) UnmarshalJSON(data []byte) (err error) {
	var t struct {
		ID          ProductID   `json:"id"`
		Name        string      `json:"name"`
		Type        ProductType `json:"type"`
		Receipt     string      `json:"receipt"`
		TicketCount int         `json:"ticketCount"`
		TicketClass string      `json:"ticketClass"`
		SKUs        []*SKU      `json:"skus"`
		Events      []EventID   `json:"events"`
	}
	if err = json.Unmarshal(data, &t); err != nil {
		return err
	}
	p.ID = t.ID
	p.Name = t.Name
	p.Type = t.Type
	p.Receipt = t.Receipt
	p.TicketCount = t.TicketCount
	p.TicketClass = t.TicketClass
	p.SKUs = t.SKUs
	if len(t.Events) == 0 {
		p.Events = nil
	} else {
		p.Events = make([]*Event, len(t.Events))
		for i, eid := range t.Events {
			p.Events[i] = &Event{ID: eid}
		}
	}
	return nil
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
	Price       int
}

func (s *SKU) MarshalJSON() ([]byte, error) {
	var t = struct {
		Coupon      string `json:"coupon,omitempty"`
		SalesStart  string `json:"salesStart,omitempty"`
		SalesEnd    string `json:"salesEnd,omitempty"`
		MembersOnly bool   `json:"membersOnly,omitempty"`
		Price       int    `json:"price"`
	}{s.Coupon, fmtJSONTime(s.SalesStart), fmtJSONTime(s.SalesEnd), s.MembersOnly, s.Price}
	return json.Marshal(&t)
}

func (s *SKU) UnmarshalJSON(data []byte) (err error) {
	var t struct {
		Coupon      string `json:"coupon"`
		SalesStart  string `json:"salesStart"`
		SalesEnd    string `json:"salesEnd"`
		MembersOnly bool   `json:"membersOnly"`
		Price       int    `json:"price"`
	}
	if err = json.Unmarshal(data, &t); err != nil {
		return err
	}
	s.Coupon = t.Coupon
	if s.SalesStart, err = parseJSONTime(t.SalesStart); err != nil {
		return err
	}
	if s.SalesEnd, err = parseJSONTime(t.SalesEnd); err != nil {
		return err
	}
	s.MembersOnly = t.MembersOnly
	s.Price = t.Price
	return nil
}

type TicketID int

type Ticket struct {
	ID    TicketID
	Event *Event
	Used  time.Time
}

func (t *Ticket) MarshalJSON() ([]byte, error) {
	var b struct {
		Event EventID `json:"event,omitempty"`
		Used  string  `json:"used,omitempty"`
	}
	if t.Event != nil {
		b.Event = t.Event.ID
	} else {
		b.Event = ""
	}
	b.Used = fmtJSONTime(t.Used)
	return json.Marshal(&b)
}

func (t *Ticket) UnmarshalJSON(_ []byte) error {
	return errors.New("ticket list not allowed")
}

func fmtJSONTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

func parseJSONTime(s string) (t time.Time, err error) {
	if s == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, s)
}
