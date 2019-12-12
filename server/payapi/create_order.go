package payapi

import (
	"errors"
	"regexp"

	"scholacantorum.org/orders/api"
	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/model"
)

var methodRE = regexp.MustCompile(`^pm_[A-Za-z0-9_]+$`)

// CreateOrder handles POST /api/pay/order requests.
func CreateOrder(r *api.Request) (err error) {
	var order *model.Order

	// Read the order details from the request.
	if order, err = api.GetOrderFromRequest(r); err != nil {
		return err
	}
	// Validate the order source and permissions.
	switch order.Source {
	case model.OrderFromPublic:
		if order.Member != 0 {
			return errors.New("invalid member ID")
		}
	case model.OrderFromMembers:
		// Get member authorization.
		if err = auth.GetSessionMembersAuth(r, r.FormValue("auth")); err != nil {
			return err
		}
		// If a member ID is specified, it must match that of the
		// session.
		if order.Member != 0 && order.Member != r.Session.Member {
			return auth.Forbidden
		}
		order.Member = r.Session.Member
	default:
		return errors.New("invalid source")
	}
	if order.Name == "" || order.Email == "" {
		return errors.New("missing name or email")
	}
	if len(order.Payments) != 0 {
		if order.Payments[0].Type != model.PaymentCard || !methodRE.MatchString(order.Payments[0].Method) {
			return errors.New("invalid payment type or method")
		}
	}

	// COMMON CODE
	return api.CreateOrderCommon(r, order)
}
