package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"mime"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/model"
)

// CreateSale handles POST /api/sale requests.
func CreateSale(tx *sql.Tx, w http.ResponseWriter, r *http.Request) {
	var body struct {
		Source              string            `json:"source"`
		Customer            model.Customer    `json:"customer"`
		StripePaymentSource string            `json:"stripePaymentSource"`
		Payment             string            `json:"payment"`
		Note                string            `json:"note"`
		Lines               []*model.SaleLine `json:"lines"`
	}
	var (
		session  *model.Session
		bodyType string
		sale     model.Sale
		customer *model.Customer
		sku      *model.SKU
		total    int
		err      error
	)

	// Get current session data, if any.
	if auth.HasSession(r) {
		if session = auth.GetSession(tx, w, r, model.PrivSetup); session == nil {
			return
		}
	}

	// Get the request data.  It may be provided as an HTTP POST form or via
	// JSON.
	bodyType, _, _ = mime.ParseMediaType(r.Header.Get("Content-Type"))
	switch bodyType {
	case "multipart/form-data", "application/x-www-form-urlencoded":
		body.Source = r.FormValue("source")
		body.Customer.Name = r.FormValue("name")
		body.Customer.Email = r.FormValue("email")
		body.Customer.Address = r.FormValue("address")
		body.Customer.City = r.FormValue("city")
		body.Customer.State = r.FormValue("state")
		body.Customer.Zip = r.FormValue("zip")
		body.Customer.Phone = r.FormValue("phone")
		body.StripePaymentSource = r.FormValue("stripePaymentSource")
		body.Payment = r.FormValue("payment")
		body.Note = r.FormValue("note")
		for i := 0; ; i++ {
			k := strconv.Itoa(i)
			if r.FormValue("sku"+k) == "" {
				break
			}
			var sl model.SaleLine
			sl.SKU, _ = strconv.Atoi(r.FormValue("sku" + k))
			sl.Qty, _ = strconv.Atoi(r.FormValue("qty" + k))
			if sl.Amount, err = strconv.Atoi(r.FormValue("amount" + k)); err != nil {
				sl.Amount = -1
			}
			body.Lines = append(body.Lines, &sl)
		}
	case "application/json":
		if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
			BadRequestError(tx, w, err.Error())
			return
		}
	default:
		tx.Rollback()
		http.Error(w, "415 Unsupported Media Type", http.StatusUnsupportedMediaType)
		return
	}

	switch body.Source {
	case "D":
		if session == nil || session.Privileges&model.PrivSell == 0 {
			ForbiddenError(tx, w)
			return
		}
	case "G":
		panic("not implemented")
	case "M":
		panic("not implemented")
	case "O":
		if session == nil || session.Privileges&model.PrivHandleOrders == 0 {
			ForbiddenError(tx, w)
			return
		}
	case "P":
		break
	default:
		BadRequestError(tx, w, "invalid sale source code")
		return
	}

	sale.Source = body.Source
	if sale.Customer = resolveSaleCustomer(tx, &body.Customer, body.StripePaymentSource != ""); sale.Customer == 0 {
		BadRequestError(tx, w, "invalid customer data")
		return
	}
	switch body.Source {
	case "D", "G", "M", "O":
		tx.Rollback()
		http.Error(w, "500 Internal Server Error: source not implemented", http.StatusInternalServerError)
		return
	case "P":
		if body.StripePaymentSource != "" || body.Payment != "" || body.Note != "" ||
			(body.Customer.ID == 0 && (body.Customer.Name == "" || body.Customer.Email == "")) ||
			body.Customer.StripeID != "" || body.Customer.MemberID != 0 {
			BadRequestError(tx, w, "invalid parameters")
			return
		}
	default:
		BadRequestError(tx, w, "invalid source")
		return
	}
	if body.Customer.ID != 0 {
		if customer = model.FetchCustomer(tx, body.Customer.ID); customer == nil {
			BadRequestError(tx, w, "no such customer")
			return
		}
		if (body.Customer.Name != "" && !strings.EqualFold(body.Customer.Name, customer.Name)) ||
			(body.Customer.Email != "" && !strings.EqualFold(body.Customer.Email, customer.Email)) {
			BadRequestError(tx, w, "customer name or email mismatch")
			return
		}
		mergeCustomer(customer, &body.Customer)
	}
	if !validateCustomer(customer) {
		BadRequestError(tx, w, "invalid customer data")
		return
	}
	if len(body.Lines) == 0 {
		BadRequestError(tx, w, "no sale lines")
		return
	}
	for _, l := range body.Lines {
		if l.ID != 0 || l.Sale != 0 || l.Qty <= 0 || l.Amount < 0 {
			BadRequestError(tx, w, "invalid sale line parameters")
			return
		}
		if sku = model.FetchSKU(tx, l.SKU); sku == nil {
			BadRequestError(tx, w, "no such SKU")
			return
		}
		if sku.Price != 0 {
			if l.Amount != 0 && l.Amount != sku.Price*l.Qty {
				BadRequestError(tx, w, "wrong sale line amount")
				return
			}
			l.Amount = sku.Price * l.Qty
		}
		total += l.Amount
	}
	sale.Source = body.Source
	sale.Timestamp = time.Now()
	if body.StripePaymentSource != "" {

	} else {
		sale.Payment = body.Payment
	}
	sale.Note = body.Note
	sale.Lines = body.Lines

	// if event.ID != 0 || event.MembersID < 0 || event.Name == "" || event.Start.IsZero() || event.Capacity < 0 {
	// 	BadRequestError(tx, w, "invalid parameters")
	// 	return
	// }
	// if event.MembersID != 0 && model.FetchEventWithMembersID(tx, event.MembersID) != nil {
	// 	BadRequestError(tx, w, "membersID already in use")
	// 	return
	// }
	customer.Save(tx)
	// event.Save(tx)
	commit(tx)
	log.Printf("%s CREATE SALE %s", session.Username, toJSON(sale))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sale)
}

// resolveSaleCustomer returns the customer ID for a sale being created.  If
// forStripe is true, it ensures that the customer is created in Stripe and has
// the data needed to process credit cards.  If any error or invalid data occur,
// it returns 0.
func resolveSaleCustomer(tx *sql.Tx, in *model.Customer, forStripe bool) model.CustomerID {
	var customer *model.Customer

	if in.ID != 0 {
		if customer = model.FetchCustomer(tx, in.ID); customer == nil {
			return 0
		}
		if (in.Name != "" && !strings.EqualFold(in.Name, customer.Name)) ||
			(in.Email != "" && !strings.EqualFold(in.Email, customer.Email)) {
			return 0
		}
		mergeCustomer(customer, in)
	} else {
		customer = in
	}
	if !validateCustomer(customer) {
		return 0
	}
	if forStripe && customer.StripeID != "" {
		// find or create customer in Stripe
		// TODO don't trust the Stripe ID in the request body.
	}
	customer.Save(tx)
	return customer.ID
}

// mergeCustomer updates c1 with the values of any non-empty fields of c2.
func mergeCustomer(c1, c2 *model.Customer) {
	if c2.StripeID != "" {
		c1.StripeID = c2.StripeID
	}
	if c2.MemberID != 0 {
		c1.MemberID = c2.MemberID
	}
	if c2.Name != "" {
		c1.Name = c2.Name
	}
	if c2.Email != "" {
		c1.Email = c2.Email
	}
	if c2.Address != "" {
		c1.Address = c2.Address
	}
	if c2.City != "" {
		c1.City = c2.City
	}
	if c2.State != "" {
		c1.State = c2.State
	}
	if c2.Zip != "" {
		c1.Zip = c2.Zip
	}
	if c2.Phone != "" {
		c1.Phone = c2.Phone
	}
}

var emailRE = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var zipRE = regexp.MustCompile(`^\d{5}(?:-\d{4})?$`)

// validateCustomer returns whether the customer has valid data.
func validateCustomer(c *model.Customer) bool {
	if c.MemberID < 0 {
		return false
	}
	if c.Email != "" && !emailRE.MatchString(c.Email) {
		return false
	}
	if (c.Address != "" || c.City != "" || c.State != "" || c.Zip != "") &&
		(c.Address == "" || c.City == "" || c.State == "" || c.Zip == "") {
		return false
	}
	if c.Zip != "" && !zipRE.MatchString(c.Zip) {
		return false
	}
	return true
}
