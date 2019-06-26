package api

import (
	"net/http"

	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
)

// GetOrder retrieves the details of the specified order.
func GetOrder(tx db.Tx, w http.ResponseWriter, r *http.Request, orderID model.OrderID) {
	var (
		session *model.Session
		order   *model.Order
	)
	// Verify permissions.
	if session = auth.GetSession(tx, w, r, model.PrivViewOrders); session == nil {
		return
	}
	// Get the requested order.
	if order = tx.FetchOrder(orderID); order == nil {
		NotFoundError(tx, w)
		return
	}
	if session.Privileges&model.PrivManageOrders == 0 && order.Flags&model.OrderValid == 0 {
		NotFoundError(tx, w)
		return
	}
	// Send back the order.
	commit(tx)
	w.Header().Set("Content-Type", "application/json")
	w.Write(emitOrder(order, false))
}
