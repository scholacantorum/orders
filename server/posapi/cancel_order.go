package posapi

import (
	"errors"
	"regexp"

	"scholacantorum.org/orders/api"
	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/model"
	"scholacantorum.org/orders/stripe"
)

var intentRE = regexp.MustCompile(`^pi_[A-Za-z0-9_]+$`)

// CancelOrder handles DELETE /api/order/${id} requests, by canceling their
// payment intent if any and deleting the order itself.  Note that the UI
// ignores errors from this call — it's already on an error path — so we also
// log errors.
func CancelOrder(r *api.Request, orderID model.OrderID) (err error) {
	var order *model.Order

	if r.Privileges&model.PrivInPersonSales == 0 {
		return auth.Forbidden
	}
	// Get the order we're supposed to cancel.
	if order = r.Tx.FetchOrder(orderID); order == nil {
		return api.NotFound
	}
	// Verify that the order is in the desired state.
	if order.Valid || len(order.Payments) != 1 || order.Payments[0].Type != model.PaymentCardPresent ||
		!intentRE.MatchString(order.Payments[0].Stripe) {
		return errors.New("order not in cancelable state")
	}
	// Cancel the payment intent.
	if err = stripe.CancelPaymentIntent(order.Payments[0].Stripe); err != nil {
		return errors.New("order cancellation failed")
		return
	}
	r.Tx.DeleteOrder(order)
	r.Tx.Commit()
	return nil
}
