package payapi

import (
	"log"
	"net/http"
	"regexp"

	"scholacantorum.org/orders/api"
	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/config"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
)

var methodRE = regexp.MustCompile(`^pm_[A-Za-z0-9_]+$`)

// CreateOrder handles POST /payapi/order requests.
//
// Parameters:
//     auth:  members site authorization token [only if source="members"]
//     source:  order source [must be "public" or "members"]
//     name:  customer name [required]
//     email:  customer email [required]
//     address:  customer street address
//     city:  customer address city
//     state:  customer address state code
//     zip:  customer address zip code
//     phone:  customer phone number
//     customer:  customer Stripe ID
//     member:  customer ID on members web site [iff source="members"]
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
//     payment#.type:  type of payment # [must be "card"]
//     payment#.subtype:  subtype of payment #
//     payment#.method:  method of payment # [must be Stripe payment method ID]
//     payment#.amount:  amount of payment #
// Emits an HTTP error status for invalid data or internal error.
// Emits JSON {"error": "..."} for card declined or other card problem.
// Emits JSON order for success.
func CreateOrder(tx db.Tx, w http.ResponseWriter, r *http.Request) {
	var (
		session *model.Session
		order   *model.Order
	)
	// Read the order details from the request.
	if order = api.GetOrderFromRequest(w, r); order == nil {
		tx.Rollback()
		return
	}
	// Validate the order source and permissions.
	switch order.Source {
	case model.OrderFromPublic:
		if order.Member != 0 {
			log.Printf("ERROR: member set in public order %s", order.ToJSON(true))
			api.BadRequestError(tx, w, "invalid member ID")
			return
		}
	case model.OrderFromMembers:
		// Get member authorization.
		if session = auth.GetSessionMembersAuth(tx, w, r, r.FormValue("auth")); session == nil {
			return
		}
		// If a member ID is specified, it must match that of the
		// session.
		if order.Member != 0 && order.Member != session.Member {
			log.Printf("ERROR: invalid session %s", order.ToJSON(true))
			api.ForbiddenError(tx, w)
			return
		}
		order.Member = session.Member
	case model.OrderFromGala:
		if r.FormValue("auth") != config.Get("galaAPIKey") {
			log.Printf("ERROR: invalid auth key for gala order")
			api.ForbiddenError(tx, w)
			return
		}
	default:
		log.Printf("ERROR: invalid source %s", order.ToJSON(true))
		api.BadRequestError(tx, w, "invalid source")
		return
	}
	if order.Name == "" || order.Email == "" {
		log.Printf("ERROR: invalid customer data in order %s", order.ToJSON(true))
		api.BadRequestError(tx, w, "invalid customer data")
		return
	}
	if len(order.Payments) != 0 {
		if order.Payments[0].Type != model.PaymentCard || !methodRE.MatchString(order.Payments[0].Method) {
			log.Printf("ERROR: invalid payment in order %s", order.ToJSON(true))
			api.BadRequestError(tx, w, "invalid payment")
			return
		}
	}

	// COMMON CODE
	api.CreateOrderCommon(tx, w, session, order)
}
