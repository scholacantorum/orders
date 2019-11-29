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

// CreateOrder handles POST /api/pos/order requests.
func CreateOrder(tx db.Tx, w http.ResponseWriter, r *http.Request) {
	var (
		order *model.Order
		err   error
	)
	// Read the order details from the request.
	if order, err = api.ParseCreateOrder(r.Body); err != nil {
		log.Printf("ERROR: can't parse body of POST /api/order request: %s", err)
		api.BadRequestError(tx, w, err.Error())
		return
	}
	CreateOrderParsed(tx, w, r, order)
}

var tokenRE = regexp.MustCompile("^tok_[A-Za-z0-9_]*$")

// CreateOrderParsed creates orders for point of sale clients.
func CreateOrderParsed(tx db.Tx, w http.ResponseWriter, r *http.Request, order *model.Order) {
	var session *model.Session

	// Get current session data, if any.
	if session = auth.GetSession(tx, w, r, model.PrivInPersonSales); session == nil {
		return
	}
	// Validate the order source and permissions.
	if order.Member != 0 {
		log.Printf("ERROR: forbidden order %s", api.EmitOrder(order, true))
		api.ForbiddenError(tx, w)
		return
	}
	if order.Source != model.OrderInPerson {
		log.Printf("ERROR: invalid source %s", api.EmitOrder(order, true))
		api.BadRequestError(tx, w, "invalid source")
		return
	}
	if len(order.Payments) != 0 {
		switch order.Payments[0].Type {
		case model.PaymentCard:
			if !tokenRE.MatchString(order.Payments[0].Method) && !methodRE.MatchString(order.Payments[0].Method) {
				log.Printf("ERROR: invalid payment in order %s", api.EmitOrder(order, true))
				api.BadRequestError(tx, w, "invalid payment")
				return
			}
		case model.PaymentCardPresent, model.PaymentCash, model.PaymentCheck:
			if order.Payments[0].Method != "" {
				log.Printf("ERROR: invalid payment in order %s", api.EmitOrder(order, true))
				api.BadRequestError(tx, w, "invalid payment")
				return
			}
		default:
			log.Printf("ERROR: invalid payment in order %s", api.EmitOrder(order, true))
			api.BadRequestError(tx, w, "invalid payment")
			return
		}
	}
	api.CreateOrderCommon(tx, w, session, order)
}
