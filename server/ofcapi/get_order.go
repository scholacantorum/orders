package ofcapi

import (
	"net/http"

	"scholacantorum.org/orders/api"
	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
)

// GetOrder handles GET /ofcapi/order/${id} requests.
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
		api.NotFoundError(tx, w)
		return
	}
	if session.Privileges&model.PrivManageOrders == 0 && !order.Valid {
		api.NotFoundError(tx, w)
		return
	}
	// Send back the order.
	api.Commit(tx)
	w.Header().Set("Content-Type", "application/json")
	w.Write(api.EmitOrder(order, false))
}
