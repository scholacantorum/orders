package payapi

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

// CreateOrder handles POST /api/pay/order requests.
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
