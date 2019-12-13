package posapi

import (
	"errors"

	"scholacantorum.org/orders/api"
	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/model"
	"scholacantorum.org/orders/stripe"
)

// SendOrderReceipt processes POST /api/order/${id}/sendReceipt requests, by
// (re-)sending the email receipt for the order.  It takes an optional email=
// query parameter to update the email address on the order.
func SendOrderReceipt(r *api.Request, orderID model.OrderID) error {
	var (
		order  *model.Order
		email  string
		card   string
		sname  string
		semail string
	)
	if r.Privileges&model.PrivInPersonSales == 0 {
		return auth.Forbidden
	}
	// Get the order whose payment we're supposed to capture.
	if order = r.Tx.FetchOrder(orderID); order == nil {
		return api.NotFound
	}
	// Verify that the order is in the desired state.
	if !order.Valid {
		return errors.New("order not complete")
	}
	// Update the email address on the order if requested.
	if email = r.FormValue("email"); email != "" && email != order.Email {
		order.Email = email
		if card = stripe.GetCardFingerprint(order.Payments[0].Stripe); card != "" {
			sname, semail = r.Tx.FetchCard(card)
			if order.Name == "" && email == semail {
				order.Name = sname
			}
			r.Tx.SaveCard(&model.Card{Card: card, Name: order.Name, Email: order.Email})
		}
		r.Tx.SaveOrder(order)
	}
	r.Tx.Commit()
	api.EmitReceipt(order, false)
	return nil
}
