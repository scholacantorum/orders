package posapi

import (
	"log"
	"net/http"
	"regexp"

	"scholacantorum.org/orders/api"
	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
)

var methodRE = regexp.MustCompile(`^pm_[A-Za-z0-9_]+$`)
var tokenRE = regexp.MustCompile("^tok_[A-Za-z0-9_]*$")

// CreateOrder handles POST /posapi/order requests.
//
// Headers:
//     Auth:  session token
// Parameters:
//     source:  order source [must be "inperson"]
//     name:  customer name
//     email:  customer email
//     address:  customer street address
//     city:  customer address city
//     state:  customer address state code
//     zip:  customer address zip code
//     phone:  customer phone number
//     customer:  customer Stripe ID
//     member:  customer ID on members web site [must be 0 or absent]
//     cNote:  order note from customer
//     oNote:  order note from office
//     inAccess:  flag whether order is in office Access database
//     coupon:  coupon code used for order
//     [line# begins at 1]
//     line#.product:  product ID for line #
//     line#.quantity:  quantity for line #
//     line#.price:  price (in cents) per unit for line #
//     line#.guestName:  name of guest for line #
//     line#.guestEmail:  email address of guest for line #
//     line#.option:  product option for line #
//     line#.used:  number of tickets used for line #
//     line#.usedAt:  event ID of event at which tickets were used for line #
//     [payment# begins at 1]
//     payment#.type:  type of payment # [must be "card", "card-present", "cash", or "check"]
//     payment#.subtype:  subtype of payment #
//     payment#.method:  method of payment # [if type is "card", must be Stripe payment method ID; otherwise must be "" or absent]
//     payment#.amount:  amount of payment #
// Emits an HTTP error status for invalid data or internal error.
// Emits JSON {"error": "..."} for card declined or other card problem.
// Emits JSON order for success.
func CreateOrder(tx db.Tx, w http.ResponseWriter, r *http.Request) {
	var (
		session *model.Session
		order   *model.Order
	)
	// Verify permissions.
	if session = auth.GetSession(tx, w, r, model.PrivInPersonSales); session == nil {
		return
	}
	// Read the order details from the request.
	if order = api.GetOrderFromRequest(w, r); order == nil {
		tx.Rollback()
		return
	}
	// Validate the order source and permissions.
	if order.Member != 0 {
		log.Printf("ERROR: forbidden order %s", order.ToJSON(true))
		api.ForbiddenError(tx, w)
		return
	}
	if order.Source != model.OrderInPerson {
		log.Printf("ERROR: invalid source %s", order.ToJSON(true))
		api.BadRequestError(tx, w, "invalid source")
		return
	}
	if len(order.Payments) != 0 {
		switch order.Payments[0].Type {
		case model.PaymentCard:
			if !tokenRE.MatchString(order.Payments[0].Method) && !methodRE.MatchString(order.Payments[0].Method) {
				log.Printf("ERROR: invalid payment in order %s", order.ToJSON(true))
				api.BadRequestError(tx, w, "invalid payment")
				return
			}
		case model.PaymentCardPresent, model.PaymentCash, model.PaymentCheck:
			if order.Payments[0].Method != "" {
				log.Printf("ERROR: invalid payment in order %s", order.ToJSON(true))
				api.BadRequestError(tx, w, "invalid payment")
				return
			}
		default:
			log.Printf("ERROR: invalid payment in order %s", order.ToJSON(true))
			api.BadRequestError(tx, w, "invalid payment")
			return
		}
	}
	api.CreateOrderCommon(tx, w, session, order)
}
