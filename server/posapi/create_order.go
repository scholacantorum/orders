package posapi

import (
	"errors"
	"regexp"

	"scholacantorum.org/orders/api"
	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/model"
)

var methodRE = regexp.MustCompile(`^pm_[A-Za-z0-9_]+$`)
var tokenRE = regexp.MustCompile("^tok_[A-Za-z0-9_]*$")

// CreateOrder handles POST /api/pos/order requests.
func CreateOrder(r *api.Request) (err error) {
	var order *model.Order

	// Verify permissions.
	if r.Privileges&model.PrivInPersonSales == 0 {
		return auth.Forbidden
	}
	// Read the order details from the request.
	if order, err = api.GetOrderFromRequest(r); err != nil {
		return err
	}
	// Validate the order source and permissions.
	if order.Member != 0 {
		return errors.New("invalid member ID")
	}
	if order.Source != model.OrderInPerson {
		return errors.New("invalid source")
	}
	if len(order.Payments) != 0 {
		switch order.Payments[0].Type {
		case model.PaymentCard:
			if !tokenRE.MatchString(order.Payments[0].Method) && !methodRE.MatchString(order.Payments[0].Method) {
				return errors.New("invalid payment method")
			}
		case model.PaymentCardPresent, model.PaymentCash, model.PaymentCheck:
			if order.Payments[0].Method != "" {
				return errors.New("invalid payment method")
			}
		default:
			return errors.New("invalid payment type")
		}
	}
	return api.CreateOrderCommon(r, order)
}
