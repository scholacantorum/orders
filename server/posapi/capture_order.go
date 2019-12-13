package posapi

import (
	"errors"
	"net/http"

	"github.com/mailru/easyjson"

	"scholacantorum.org/orders/api"
	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/model"
	"scholacantorum.org/orders/stripe"
)

// CaptureOrderPayment processes POST /api/order/${id}/capturePayment requests,
// by capturing the (already authorized) payment in the order.
func CaptureOrderPayment(r *api.Request, orderID model.OrderID) (err error) {
	var (
		order          *model.Order
		card           string
		tentativeEmail string
	)
	if r.Privileges&model.PrivInPersonSales == 0 {
		return auth.Forbidden
	}
	// Get the order whose payment we're supposed to capture.
	if order = r.Tx.FetchOrder(orderID); order == nil {
		return api.NotFound
	}
	// Verify that the order is in the desired state.
	if order.Valid || len(order.Payments) != 1 || order.Payments[0].Type != model.PaymentCardPresent ||
		!intentRE.MatchString(order.Payments[0].Method) {
		return errors.New("order not in capturable state")
	}
	if card, err = stripe.CapturePayment(order, order.Payments[0]); err != nil {
		return api.HTTPError(http.StatusInternalServerError, "500 Internal Server Error")
	}
	order.Valid = true
	r.Tx.SaveOrder(order)
	_, tentativeEmail = r.Tx.FetchCard(card)
	r.Tx.Commit()
	if order.Email != "" {
		api.EmitReceipt(order, false)
	}
	api.UpdateGoogleSheet(order)
	if order.Email == "" {
		order.Email = tentativeEmail
	}
	easyjson.MarshalToHTTPResponseWriter(order, r)
	return nil
}
