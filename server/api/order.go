package api

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"scholacantorum.org/orders/config"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
	"scholacantorum.org/orders/stripe"
)

var emailRE = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var stateRE = regexp.MustCompile(`^[A-Z][A-Z]$`)
var zipRE = regexp.MustCompile(`^\d{5}(?:-\d{4})?$`)
var customerRE = regexp.MustCompile(`^cus_[A-Za-z0-9]+$`)

// GetOrderFromRequest reads the order details from the request body and returns
// them.  If it returns nil, the order details were invalid; an appropriate
// error response has been issued, and the error has been logged.
func GetOrderFromRequest(w http.ResponseWriter, r *http.Request) (o *model.Order) {
	var err error

	o = new(model.Order)
	o.Source = model.OrderSource(r.FormValue("source"))
	switch o.Source {
	case "":
		o.Source = model.OrderFromPublic
	case model.OrderFromPublic, model.OrderFromMembers, model.OrderFromGala, model.OrderFromOffice, model.OrderInPerson:
		// no-op
	default:
		log.Printf("ERROR: invalid source %q", o.Source)
		http.Error(w, `400 Bad Request: invalid "source"`, http.StatusBadRequest)
		goto ERROR
	}
	o.Name = strings.TrimSpace(r.FormValue("name"))
	o.Email = strings.TrimSpace(r.FormValue("email"))
	if o.Email != "" && !emailRE.MatchString(o.Email) {
		log.Printf("ERROR: invalid email %q", o.Email)
		http.Error(w, `400 Bad Request: invalid "email"`, http.StatusBadRequest)
		goto ERROR
	}
	o.Address = strings.TrimSpace(r.FormValue("address"))
	o.City = strings.TrimSpace(r.FormValue("city"))
	o.State = strings.ToUpper(strings.TrimSpace(r.FormValue("state")))
	if o.State != "" && !stateRE.MatchString(o.State) {
		log.Printf("ERROR: invalid state %q", o.State)
		http.Error(w, `400 Bad Request: invalid "state"`, http.StatusBadRequest)
		goto ERROR
	}
	o.Zip = strings.TrimSpace(r.FormValue("zip"))
	if o.Zip != "" && !zipRE.MatchString(o.Zip) {
		log.Printf("ERROR: invalid zip %q", o.Zip)
		http.Error(w, `400 Bad Request: invalid "zip"`, http.StatusBadRequest)
		goto ERROR
	}
	if (o.Address != "" || o.City != "" || o.State != "" || o.Zip != "") &&
		(o.Address == "" || o.City == "" || o.State == "" || o.Zip == "") {
		log.Printf("ERROR: have address %v city %v state %v zip %v",
			o.Address != "", o.City != "", o.State != "", o.Zip != "")
		http.Error(w, `400 Bad Request: specify all or none of "address"+"city"+"state"+"zip"`, http.StatusBadRequest)
		goto ERROR
	}
	o.Phone = strings.TrimSpace(r.FormValue("phone"))
	o.Customer = strings.TrimSpace(r.FormValue("customer"))
	if o.Customer != "" && !customerRE.MatchString(o.Customer) {
		log.Printf("ERROR: invalid customer %q", o.Customer)
		http.Error(w, `400 Bad Request: invalid "customer"`, http.StatusBadRequest)
		goto ERROR
	}
	if mstr := r.FormValue("member"); mstr != "" {
		if o.Member, err = strconv.Atoi(mstr); err != nil || o.Member < 1 {
			log.Printf("ERROR: invalid member %q", mstr)
			http.Error(w, `400 Bad Request: invalid "member"`, http.StatusBadRequest)
			goto ERROR
		}
	}
	o.CNote = strings.TrimSpace(r.FormValue("cNote"))
	o.ONote = strings.TrimSpace(r.FormValue("oNote"))
	if iastr := r.FormValue("inAccess"); iastr != "" {
		if o.InAccess, err = strconv.ParseBool(iastr); err != nil {
			log.Printf("ERROR: invalid inAccess %q", iastr)
			http.Error(w, `400 Bad Request: invalid "inAccess"`, http.StatusBadRequest)
			goto ERROR
		}
	}
	o.Coupon = strings.ToUpper(strings.TrimSpace(r.FormValue("coupon")))
	for idx := 1; true; idx++ {
		var (
			ol     model.OrderLine
			prefix = fmt.Sprintf("line%d.", idx)
		)
		if pname := r.FormValue(prefix + "product"); pname != "" {
			ol.Product = &model.Product{ID: model.ProductID(pname)}
		} else {
			break
		}
		if ol.Quantity, err = strconv.Atoi(r.FormValue(prefix + "quantity")); err != nil || ol.Quantity < 1 {
			log.Printf("ERROR: invalid quantity %q", r.FormValue(prefix+"quantity"))
			http.Error(w, `400 Bad Request: invalid "quantity"`, http.StatusBadRequest)
			goto ERROR
		}
		if ol.Price, err = strconv.Atoi(r.FormValue(prefix + "price")); err != nil || ol.Quantity < 0 {
			log.Printf("ERROR: invalid price %q", r.FormValue(prefix+"price"))
			http.Error(w, `400 Bad Request: invalid "price"`, http.StatusBadRequest)
			goto ERROR
		}
		ol.GuestName = r.FormValue(prefix + "guestName")
		ol.GuestEmail = r.FormValue(prefix + "guestEmail")
		ol.Option = r.FormValue(prefix + "option")
		if uval := r.FormValue(prefix + "used"); uval != "" {
			if ol.Used, err = strconv.Atoi(uval); err != nil || ol.Used < 0 {
				log.Printf("ERROR: invalid used amount %q", uval)
				http.Error(w, `400 Bad Request: invalid "used"`, http.StatusBadRequest)
				goto ERROR
			}
			if ol.UsedAt = model.EventID(r.FormValue(prefix + "usedAt")); ol.Used > 0 && ol.UsedAt == "" {
				log.Printf("ERROR: missing usedAt")
				http.Error(w, `400 Bad Request: "usedAt" is required when "used" is nonzero`, http.StatusBadRequest)
				goto ERROR
			}
		}
		o.Lines = append(o.Lines, &ol)
	}
	for idx := 1; true; idx++ {
		var (
			p      model.Payment
			prefix = fmt.Sprintf("payment%d.", idx)
		)
		if p.Type = model.PaymentType(r.FormValue(prefix + "type")); p.Type == "" {
			break
		}
		switch p.Type {
		case model.PaymentCard, model.PaymentCardPresent, model.PaymentCash, model.PaymentCheck, model.PaymentOther:
			// no-op
		default:
			log.Printf("ERROR: invalid payment type %q", p.Type)
			http.Error(w, `400 Bad Request: invalid "type"`, http.StatusBadRequest)
			goto ERROR
		}
		p.Subtype = strings.TrimSpace(r.FormValue(prefix + "subtype"))
		p.Method = strings.TrimSpace(r.FormValue(prefix + "method"))
		if p.Amount, err = strconv.Atoi(r.FormValue(prefix + "amount")); err != nil {
			log.Printf("ERROR: invalid payment amount %q", r.FormValue(prefix+"amount"))
			http.Error(w, `400 Bad Request: invalid "amount"`, http.StatusBadRequest)
			goto ERROR
		}
		o.Payments = append(o.Payments, &p)
	}
	return o

ERROR:
	log.Printf("    in %s %s %+v", r.Method, r.RequestURI, r.Form)
	return nil
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
		log.Printf("ERROR: invalid products or prices in order %s", order.ToJSON(true))
		BadRequestError(tx, w, "invalid products or prices")
		return
	}
	// Validate the customer data.
	if !validateCustomer(tx, order, session) {
		log.Printf("ERROR: invalid customer data in order %s", order.ToJSON(true))
		BadRequestError(tx, w, "invalid customer data")
		return
	}
	// Make sure the rest of the order details are OK.
	if !validateOrderDetails(tx, order, privs) {
		log.Printf("ERROR: invalid parameters in order %s", order.ToJSON(true))
		BadRequestError(tx, w, "invalid parameters")
		return
	}
	// Calculate the order total and verify the payment.
	if !validatePayment(order) {
		log.Printf("ERROR: invalid payment in order %s", order.ToJSON(true))
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
		receipt = true
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
			receipt = true
		case model.PaymentCard:
			success, card, message = stripe.ChargeCard(order, order.Payments[0])
			tx = db.Begin()
			if !success {
				tx.DeleteOrder(order)
				Commit(tx)
				if message == "" {
					message = "We're sorry, but our payment processor isn't working right now.  Please try again later, or contact our office at (650) 254-1700."
				}
				log.Printf("ERROR: payment rejected (%q) in order %s", message, order.ToJSON(true))
				SendError(tx, w, message)
				return
			}
			order.Valid = true
			tx.SaveOrder(order)
			tx.SaveCard(card, order.Name, order.Email)
			receipt = true
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
				log.Printf("ERROR: can't create payment intent for order %s", order.ToJSON(true))
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
		log.Printf("%s %s ORDER %s", session.Username, logverb, order.ToJSON(true))
	} else {
		log.Printf("- %s ORDER %s", logverb, order.ToJSON(true))
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(order.ToJSON(false))
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

// validateCustomer returns whether the customer data in the order are valid.
func validateCustomer(tx db.Tx, order *model.Order, session *model.Session) bool {

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

		// Additional constraints by product type:
		switch line.Product.Type {
		case model.ProdAuctionItem:
			// Auction items aren't supported yet.
			return false // TODO
		case model.ProdDonation, model.ProdRecording, model.ProdSheetMusic, model.ProdRegistration:
			// Donations, concert recordings, sheet music, and event
			// registrations must have a quantity of 1.
			if line.Quantity != 1 || line.Used != 0 || line.UsedAt != "" {
				return false
			}
		case model.ProdWardrobe:
			if line.Quantity < 1 || line.Used != 0 || line.UsedAt != "" {
				return false
			}
		case model.ProdTicket:
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
	// And it should have the correct amount.
	var pmt = order.Payments[0]
	if pmt.Amount != total {
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
