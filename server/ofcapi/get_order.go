package ofcapi

import (
	"scholacantorum.org/orders/api"
	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/model"
)

// GetOrder handles GET /ofcapi/order/${id} requests.  Note that this is called
// by the members site to validate recording orders, as well as by the office
// UI.
func GetOrder(r *api.Request, orderID model.OrderID) error {
	var order *model.Order

	// Verify permissions.
	if r.Privileges&model.PrivViewOrders == 0 {
		return auth.Forbidden
	}
	// Get the requested order.
	if order = r.Tx.FetchOrder(orderID); order == nil {
		return api.NotFound
	}
	if r.Privileges&model.PrivManageOrders == 0 && !order.Valid {
		return api.NotFound
	}
	// Send back the order.
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json")
	r.Write(order.ToJSON(false))
	return nil
}
