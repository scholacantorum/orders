package api

import (
	"log"
	"net/http"

	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
	"scholacantorum.org/orders/stripe"
)

// SendOrderReceipt processes POST /api/order/${id}/sendReceipt requests, by
// (re-)sending the email receipt for the order.  It takes an optional email=
// query parameter to update the email address on the order.
func SendOrderReceipt(tx db.Tx, w http.ResponseWriter, r *http.Request, orderID model.OrderID) {
	var (
		session *model.Session
		order   *model.Order
		email   string
		card    string
		sname   string
		semail  string
	)
	// Get current session data, if any.
	if session = auth.GetSession(tx, w, r, model.PrivInPersonSales); session == nil {
		return
	}
	// Get the order whose payment we're supposed to capture.
	if order = tx.FetchOrder(orderID); order == nil {
		NotFoundError(tx, w)
		return
	}
	// Verify that the order is in the desired state.
	if order.Flags&model.OrderValid == 0 {
		BadRequestError(tx, w, "order not complete")
		return
	}
	// Update the email address on the order if requested.
	if email = r.FormValue("email"); email != "" && email != order.Email {
		order.Email = email
		if card = stripe.GetCardFingerprint(order.Payments[0].Stripe); card != "" {
			sname, semail = tx.FetchCard(card)
			if order.Name == "" && email == semail {
				order.Name = sname
			}
			tx.SaveCard(card, order.Name, order.Email)
		}
		tx.SaveOrder(order)
	}
	commit(tx)
	log.Printf("- RESEND RECEIPT for order %d to %s", orderID, order.Email)
	w.WriteHeader(http.StatusNoContent)
	EmitReceipt(order, false)
}
