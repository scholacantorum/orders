package posapi

import (
	"log"
	"net/http"

	"scholacantorum.org/orders/api"
	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
	"scholacantorum.org/orders/stripe"
)

// CaptureOrderPayment processes POST /api/order/${id}/capturePayment requests,
// by capturing the (already authorized) payment in the order.
func CaptureOrderPayment(tx db.Tx, w http.ResponseWriter, r *http.Request, orderID model.OrderID) {
	var (
		session        *model.Session
		order          *model.Order
		card           string
		tentativeEmail string
		err            error
	)
	// Get current session data, if any.
	if session = auth.GetSession(tx, w, r, model.PrivInPersonSales); session == nil {
		return
	}
	// Get the order whose payment we're supposed to capture.
	if order = tx.FetchOrder(orderID); order == nil {
		log.Printf("ERROR: capture of nonexistent order %d", orderID)
		api.NotFoundError(tx, w)
		return
	}
	// Verify that the order is in the desired state.
	if order.Valid || len(order.Payments) != 1 || order.Payments[0].Type != model.PaymentCardPresent ||
		!intentRE.MatchString(order.Payments[0].Method) {
		log.Printf("ERROR: capture of order in wrong state %s", order.ToJSON(true))
		api.BadRequestError(tx, w, "order not in capturable state")
		return
	}
	if card, err = stripe.CapturePayment(order, order.Payments[0]); err != nil {
		api.Commit(tx)
		log.Printf("ERROR: failed to capture payment for order %d: %s", order.ID, err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	order.Valid = true
	tx.SaveOrder(order)
	_, tentativeEmail = tx.FetchCard(card)
	api.Commit(tx)
	log.Printf("- CAPTURE ORDER %s", order.ToJSON(true))
	if order.Email != "" {
		api.EmitReceipt(order, false)
	}
	api.UpdateGoogleSheet(order)
	if order.Email == "" {
		order.Email = tentativeEmail
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(order.ToJSON(false))
}
