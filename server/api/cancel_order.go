package api

import (
	"log"
	"net/http"

	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
	"scholacantorum.org/orders/stripe"
)

// CancelOrder handles DELETE /api/order/${id} requests, by canceling their
// payment intent if any and deleting the order itself.  Note that the UI
// ignores errors from this call — it's already on an error path — so we also
// log errors.
func CancelOrder(tx db.Tx, w http.ResponseWriter, r *http.Request, orderID model.OrderID) {
	var (
		session *model.Session
		order   *model.Order
		err     error
	)
	// Get current session data, if any.
	if session = auth.GetSession(tx, w, r, model.PrivInPersonSales); session == nil {
		return
	}
	// Get the order we're supposed to cancel.
	if order = tx.FetchOrder(orderID); order == nil {
		NotFoundError(tx, w)
		return
	}
	// Verify that the order is in the desired state.
	if order.Flags&model.OrderValid != 0 || len(order.Payments) != 1 || order.Payments[0].Type != model.PaymentCardPresent ||
		!intentRE.MatchString(order.Payments[0].Stripe) {
		log.Printf("ERROR: cannot cancel order %d as requested because it is not in the proper state", orderID)
		BadRequestError(tx, w, "order not in cancelable state")
		return
	}
	// Cancel the payment intent.
	if err = stripe.CancelPaymentIntent(order.Payments[0].Stripe); err != nil {
		log.Printf("ERROR: cannot cancel payment intent %s for order %d: %s",
			order.Payments[0].Stripe, order.ID, err)
		BadRequestError(tx, w, "order not in cancelable state")
		return
	}
	tx.DeleteOrder(order)
	commit(tx)
	log.Printf("- CANCEL ORDER %s", emitOrder(order, true))
	w.WriteHeader(http.StatusNoContent)
}
