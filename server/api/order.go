package api

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"github.com/rothskeller/json"
	"scholacantorum.org/orders/config"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
	"scholacantorum.org/orders/stripe"
)

// ParseCreateOrder reads the order details from the request body.
func ParseCreateOrder(r io.Reader) (o *model.Order, err error) {
	var (
		jr = json.NewReader(r)
	)
	o = new(model.Order)
	err = jr.Read(json.ObjectHandler(func(key string) json.Handlers {
		switch key {
		case "source":
			return json.StringHandler(func(s string) { o.Source = model.OrderSource(s) })
		case "name":
			return json.StringHandler(func(s string) { o.Name = s })
		case "email":
			return json.StringHandler(func(s string) { o.Email = s })
		case "address":
			return json.StringHandler(func(s string) { o.Address = s })
		case "city":
			return json.StringHandler(func(s string) { o.City = s })
		case "state":
			return json.StringHandler(func(s string) { o.State = s })
		case "zip":
			return json.StringHandler(func(s string) { o.Zip = s })
		case "phone":
			return json.StringHandler(func(s string) { o.Phone = s })
		case "customer":
			return json.StringHandler(func(s string) { o.Customer = s })
		case "member":
			return json.IntHandler(func(i int) { o.Member = i })
		case "cNote":
			return json.StringHandler(func(s string) { o.CNote = s })
		case "oNote":
			return json.StringHandler(func(s string) { o.ONote = s })
		case "coupon":
			return json.StringHandler(func(s string) { o.Coupon = s })
		case "lines":
			return json.ArrayHandler(func() json.Handlers {
				var ol model.OrderLine
				o.Lines = append(o.Lines, &ol)
				return parseCreateOrderLine(&ol)
			})
		case "payments":
			return json.ArrayHandler(func() json.Handlers {
				var p model.Payment
				o.Payments = append(o.Payments, &p)
				return parseCreateOrderPayment(&p)
			})
		default:
			return json.RejectHandler()
		}
	}))
	return o, err
}
func parseCreateOrderLine(ol *model.OrderLine) json.Handlers {
	return json.ObjectHandler(func(key string) json.Handlers {
		switch key {
		case "product":
			return json.StringHandler(func(s string) { ol.Product = &model.Product{ID: model.ProductID(s)} })
		case "quantity":
			return json.IntHandler(func(i int) { ol.Quantity = i })
		case "used":
			return json.IntHandler(func(i int) { ol.Used = i })
		case "usedAt":
			return json.StringHandler(func(s string) { ol.UsedAt = model.EventID(s) })
		case "price":
			return json.IntHandler(func(i int) { ol.Price = i })
		default:
			return json.RejectHandler()
		}
	})
}
func parseCreateOrderPayment(p *model.Payment) json.Handlers {
	return json.ObjectHandler(func(key string) json.Handlers {
		switch key {
		case "type":
			return json.StringHandler(func(s string) { p.Type = model.PaymentType(s) })
		case "subtype":
			return json.StringHandler(func(s string) { p.Subtype = s })
		case "method":
			return json.StringHandler(func(s string) { p.Method = s })
		case "amount":
			return json.IntHandler(func(i int) { p.Amount = i })
		default:
			return json.RejectHandler()
		}
	})
}

// CreateOrderCommon is the common part of creating an order, shared by the
// various APIs.  Each of them makes authorization and validity checks specific
// to it, and then calls this function to perform the common checks and create
// the order.
func CreateOrderCommon(tx db.Tx, w http.ResponseWriter, session *model.Session, order *model.Order) {
	var (
		privs   model.Privilege
		success bool
		card    string
		message string
		receipt bool
		logverb = "PLACE"
	)
	if session != nil {
		privs = session.Privileges
	}
	// Resolve the products and SKUs and validate the prices.
	if !resolveSKUs(tx, order) {
		log.Printf("ERROR: invalid products or prices in order %s", EmitOrder(order, true))
		BadRequestError(tx, w, "invalid products or prices")
		return
	}
	// Validate the customer data.
	if !validateCustomer(tx, order, session) {
		log.Printf("ERROR: invalid customer data in order %s", EmitOrder(order, true))
		BadRequestError(tx, w, "invalid customer data")
		return
	}
	// Make sure the rest of the order details are OK.
	if !validateOrderDetails(tx, order, privs) {
		log.Printf("ERROR: invalid parameters in order %s", EmitOrder(order, true))
		BadRequestError(tx, w, "invalid parameters")
		return
	}
	// Calculate the order total and verify the payment.
	if !validatePayment(order) {
		log.Printf("ERROR: invalid payment in order %s", EmitOrder(order, true))
		BadRequestError(tx, w, "invalid payment")
		return
	}
	// Assign a token to the order (after the new transaction is opened, to
	// ensure uniqueness).
	order.Token = newOrderToken(tx)
	// Generate tickets if needed.
	generateTickets(tx, order)
	// If we don't have to charge a card through Stripe, the order is now
	// complete.
	if len(order.Payments) == 0 {
		order.Valid = true
	} else {
		switch order.Payments[0].Type {
		case model.PaymentCard, model.PaymentCardPresent:
			break
		default:
			order.Valid = true
		}
	}
	// Save the order to the database.
	tx.SaveOrder(order)
	Commit(tx)
	// If we do have to charge a card through Stripe, do it now.
	if len(order.Payments) == 1 {
		switch order.Payments[0].Type {
		case model.PaymentCash, model.PaymentCheck, model.PaymentOther:
			receipt = order.Email != ""
		case model.PaymentCard:
			success, card, message = stripe.ChargeCard(order, order.Payments[0])
			tx = db.Begin()
			if !success {
				tx.DeleteOrder(order)
				Commit(tx)
				if message == "" {
					message = "We're sorry, but our payment processor isn't working right now.  Please try again later, or contact our office at (650) 254-1700."
				}
				log.Printf("ERROR: payment rejected (%q) in order %s", message, EmitOrder(order, true))
				SendError(tx, w, message)
				return
			}
			order.Valid = true
			tx.SaveOrder(order)
			tx.SaveCard(card, order.Name, order.Email)
			receipt = order.Email != ""
			order.Name, order.Email = tx.FetchCard(card)
			Commit(tx)
		case model.PaymentCardPresent:
			// For card present transactions, we have to create the
			// order and notify Stripe before processing the card.
			// Do that now, and return the (uncompleted) order with
			// the payment intent in it.
			success = stripe.CreatePaymentIntent(order)
			tx = db.Begin()
			if !success {
				tx.DeleteOrder(order)
				Commit(tx)
				SendError(tx, w, "We're sorry, but our payment processor isn't working right now.  Please try again later, or contact our office at (650) 254-1700.")
				log.Printf("ERROR: can't create payment intent for order %s", EmitOrder(order, true))
				return
			}
			tx.SaveOrder(order)
			Commit(tx)
			logverb = "CREATE"
			receipt = false
		}
	}
	// Log and return the completed order.
	if session != nil {
		log.Printf("%s %s ORDER %s", session.Username, logverb, EmitOrder(order, true))
	} else {
		log.Printf("- %s ORDER %s", logverb, EmitOrder(order, true))
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(EmitOrder(order, false))
	if receipt && order.Email != "" {
		EmitReceipt(order, false)
	}
	if order.Valid {
		UpdateGoogleSheet(order)
	}
}

// resolveSKUs walks through each line of the order, finding the listed product
// and verifying the amount of the order line, following the SKU rules
// documented in db/schema.sql. It returns true if everything resolved
// successfully and false otherwise. Note that if a coupon is specified in the
// order but not used by any SKU, it is removed; this keeps the order reporting
// system clean of invalid coupon codes.
func resolveSKUs(tx db.Tx, order *model.Order) bool {
	var couponMatch bool

	for _, line := range order.Lines {
		var sku *model.SKU

		if line.Product = tx.FetchProduct(line.Product.ID); line.Product == nil {
			return false
		}
		if line.Product.Type == model.ProdAuctionItem || line.Product.Type == model.ProdDonation {
			if line.Price < 1 {
				return false
			}
			continue
		}
		for _, s := range line.Product.SKUs {
			if !MatchingSKU(s, order.Coupon, order.Source, false) {
				continue
			}
			if s.Coupon != "" {
				couponMatch = true
			}
			sku = BetterSKU(sku, s)
		}
		if sku == nil {
			continue
		}
		if line.Price != sku.Price {
			return false
		}
	}
	if !couponMatch {
		order.Coupon = ""
	}
	return true
}

var emailRE = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var stateRE = regexp.MustCompile(`^[A-Z][A-Z]$`)
var zipRE = regexp.MustCompile(`^\d{5}(?:-\d{4})?$`)
var customerRE = regexp.MustCompile(`^cus_[A-Za-z0-9]+$`)

// validateCustomer returns whether the customer data in the order are valid.
func validateCustomer(tx db.Tx, order *model.Order, session *model.Session) bool {

	// Email must be a valid email address according to the regular
	// expression above (which is the same RE used by <input type="email">
	// in HTML5).
	if order.Email != "" && !emailRE.MatchString(order.Email) {
		return false
	}

	// If any of the address fields is set, they must all be set, and they
	// need to match the appropriate regular expressions.
	if (order.Address != "" || order.City != "" || order.State != "" || order.Zip != "") &&
		(order.Address == "" || order.City == "" || order.State == "" || order.Zip == "") {
		return false
	}
	if order.State != "" && !stateRE.MatchString(order.State) {
		return false
	}
	if order.Zip != "" && !zipRE.MatchString(order.Zip) {
		return false
	}
	if order.Customer != "" && !customerRE.MatchString(order.Customer) {
		return false
	}

	// A Stripe customer ID is allowed only for gala sales.  TODO
	if order.Customer != "" {
		return false
	}

	// A member ID is required for sheet music or concert recording sales.
	for _, line := range order.Lines {
		if line.Product.Type == model.ProdRecording || line.Product.Type == model.ProdSheetMusic {
			if order.Member == 0 {
				return false
			}
		}
	}
	return true
}

// validateOrderDetails returns whether the order details are valid.  Note that
// this does not check authorization.  It also doesn't check anything specific
// to the product type or the order type.
func validateOrderDetails(tx db.Tx, order *model.Order, privs model.Privilege) bool {

	// New orders should not have an ID or a created timestamp, and must
	// have at least one line.
	if order.ID != 0 || order.Token != "" || !order.Created.IsZero() || len(order.Lines) == 0 {
		return false
	}
	order.Created = time.Now()

	// Office notes are allowed only by office staff.
	if order.ONote != "" && privs&model.PrivManageOrders == 0 {
		return false
	}

	// Remove any lines with zero quantity.  Make sure there's at least one
	// line left.
	var j = 0
	for i := range order.Lines {
		if order.Lines[i].Quantity < 0 {
			return false
		}
		if order.Lines[i].Quantity != 0 {
			order.Lines[j] = order.Lines[i]
			j++
		}
	}
	order.Lines = order.Lines[:j]
	if len(order.Lines) == 0 {
		return false
	}

	// Check the validity of each order line.
	for _, line := range order.Lines {

		// Lines must not have IDs or tickets.
		if line.ID != 0 || len(line.Tickets) != 0 {
			return false
		}

		// Additional constraints by product type:
		switch line.Product.Type {
		case model.ProdAuctionItem:
			// Auction items aren't supported yet.
			return false // TODO
		case model.ProdDonation, model.ProdRecording, model.ProdSheetMusic:
			// Donations, concert recordings, and sheet music must
			// have a quantity of 1.
			if line.Quantity != 1 || line.Used != 0 || line.UsedAt != "" {
				return false
			}
		case model.ProdWardrobe:
			if line.Quantity < 1 || line.Used != 0 || line.UsedAt != "" {
				return false
			}
		case model.ProdTicket:
			line.AutoUse = line.Quantity
			if line.Used < 0 || line.Used > line.Quantity*line.Product.TicketCount {
				return false
			}
			if line.Used != 0 {
				var found bool
				for _, e := range line.Product.Events {
					if e.Event.ID == line.UsedAt {
						found = true
						break
					}
				}
				if !found {
					return false
				}
			}
		}
	}
	return true
}

// validatePayment returns whether the order payment is valid for the order type
// and has the correct amount.
func validatePayment(order *model.Order) bool {
	var total int

	// Calculate the order total.
	for _, ol := range order.Lines {
		total += ol.Price * ol.Quantity
	}
	// If this is a free order, it's OK if there is no payment.
	if total == 0 && len(order.Payments) == 0 {
		return true
	}
	// Otherwise, there should be exactly one payment.
	if len(order.Payments) != 1 {
		return false
	}
	// And it should not have an ID, a Stripe ID, or a created timestamp,
	// and it should have the correct amount.
	var pmt = order.Payments[0]
	if pmt.ID != 0 || pmt.Stripe != "" || !pmt.Created.IsZero() || pmt.Amount != total {
		return false
	}
	// If this is a free order and has a payment, its type must be "cash".
	// We remove it; no point in storing a zero payment.
	if pmt.Amount == 0 {
		if pmt.Type != model.PaymentCash {
			return false
		}
		order.Payments = order.Payments[:0]
		return true
	}
	pmt.Created = order.Created
	return true
}

// generateTickets creates tickets as needed for the new order.
func generateTickets(tx db.Tx, order *model.Order) {
	var (
		event *model.Event
		found bool
	)
	// Walk through all of the order lines.
	for _, ol := range order.Lines {
		// We only care about lines on which tickets are sold.
		if ol.Product.TicketCount == 0 {
			continue
		}
		// Figure out whether this ticket is allocated to a particular
		// event.
		for _, pe := range ol.Product.Events {
			// If the event has priority zero, the ticket is
			// dedicated to that event by definition.
			if pe.Priority == 0 {
				event = pe.Event
				break
			}
			// Otherwise, look to see if it is the only future event
			// at which the ticket is valid.  "Future" is taken with
			// one hour slop to allow for at-the-door sales after
			// curtain.
			if pe.Event.Start.After(time.Now().Add(-time.Hour)) {
				if found {
					event = nil // multiple matches
				} else {
					found = true
					event = pe.Event
				}
			}
		}
		// Create the ticket objects.
		for i := 0; i < ol.Product.TicketCount*ol.Quantity; i++ {
			var tick = model.Ticket{Event: event}
			if ol.Used > 0 {
				tick.Event = &model.Event{ID: ol.UsedAt}
				tick.Used = order.Created
				ol.Used--
			}
			ol.Tickets = append(ol.Tickets, &tick)
		}
	}
}

// newOrderToken generates a token for a new order, retrying until it has one
// that's unique.
func newOrderToken(tx db.Tx) (token string) {
	for token == "" || tx.FetchOrderByToken(token) != nil {
		token = NewToken()
	}
	return token
}

// EmitOrder generates a JSON representation of the order.  If log is true, it
// includes internal details.
func EmitOrder(o *model.Order, log bool) []byte {
	var (
		buf bytes.Buffer
		jw  = json.NewWriter(&buf)
	)
	jw.Object(func() {
		jw.Prop("id", int(o.ID))
		if log {
			jw.Prop("token", o.Token)
		}
		jw.Prop("source", string(o.Source))
		if o.Name != "" {
			jw.Prop("name", o.Name)
		}
		if o.Email != "" {
			jw.Prop("email", o.Email)
		}
		if o.Address != "" {
			jw.Prop("address", o.Address)
			jw.Prop("city", o.City)
			jw.Prop("state", o.State)
			jw.Prop("zip", o.Zip)
		}
		if o.Phone != "" {
			jw.Prop("phone", o.Phone)
		}
		if o.Customer != "" {
			jw.Prop("customer", o.Customer)
		}
		if o.Member != 0 {
			jw.Prop("member", o.Member)
		}
		jw.Prop("created", o.Created.Format(time.RFC3339))
		if log {
			jw.Prop("valid", o.Valid)
			jw.Prop("inAccess", o.InAccess)
		}
		if o.CNote != "" {
			jw.Prop("cNote", o.CNote)
		}
		if o.ONote != "" {
			jw.Prop("oNote", o.ONote)
		}
		if o.Coupon != "" {
			jw.Prop("coupon", o.Coupon)
		}
		jw.Prop("lines", func() {
			jw.Array(func() {
				for _, ol := range o.Lines {
					jw.Object(func() {
						if log {
							jw.Prop("id", int(ol.ID))
						}
						jw.Prop("product", string(ol.Product.ID))
						jw.Prop("quantity", ol.Quantity)
						jw.Prop("price", ol.Price)
						if len(ol.Tickets) != 0 {
							jw.Prop("tickets", func() {
								jw.Array(func() {
									for _, t := range ol.Tickets {
										jw.Object(func() {
											if log {
												jw.Prop("id", int(t.ID))
											}
											if t.Event != nil {
												jw.Prop("event", string(t.Event.ID))
											}
											if !t.Used.IsZero() {
												jw.Prop("used", t.Used.Format(time.RFC3339))
											}
										})
									}
								})
							})
						}
					})
				}
			})
		})
		jw.Prop("payments", func() {
			jw.Array(func() {
				for _, p := range o.Payments {
					jw.Object(func() {
						if log {
							jw.Prop("id", int(p.ID))
						}
						jw.Prop("type", string(p.Type))
						if p.Subtype != "" {
							jw.Prop("subtype", p.Subtype)
						}
						jw.Prop("method", p.Method)
						if p.Stripe != "" && log {
							jw.Prop("stripe", p.Stripe)
						}
						jw.Prop("created", p.Created.Format(time.RFC3339))
						jw.Prop("amount", p.Amount)
					})
				}
			})
		})
	})
	jw.Close()
	return buf.Bytes()
}

// UpdateGoogleSheet starts a subprocess that updates the Google orders
// spreadsheet with information about the specified order.
func UpdateGoogleSheet(order *model.Order) {
	var (
		cmd *exec.Cmd
		err error
	)
	cmd = exec.Command(config.Get("bin")+"/update-orders-sheet", strconv.Itoa(int(order.ID)))
	if err = cmd.Start(); err != nil {
		log.Printf("ERROR: can't update orders sheet for order %d: %s", order.ID, err)
		return
	}
	// Note that we are intentionally not waiting for the subprocess to
	// finish.  This CGI script will exit immediately, so that the user gets
	// a fast response to their order.  The subprocess will continue as an
	// orphan, and its zombie will be reaped by the init daemon.
}
