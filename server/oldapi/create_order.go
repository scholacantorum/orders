package oldapi

import (
	"log"
	"net/http"

	"scholacantorum.org/orders/api"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
	"scholacantorum.org/orders/payapi"
	"scholacantorum.org/orders/posapi"
)

// CreateOrder handles POST /api/order requests.  This is an obsolete API; the
// function detects which of the newer APIs was intended and defers the call to
// the handler of that API.
func CreateOrder(tx db.Tx, w http.ResponseWriter, r *http.Request) {
	var (
		order *model.Order
		err   error
	)
	// Read the order details from the request.
	if order, err = api.ParseCreateOrder(r.Body); err != nil {
		log.Printf("ERROR: can't parse body of POST /api/order request: %s", err)
		api.BadRequestError(tx, w, err.Error())
		return
	}
	switch order.Source {
	case model.OrderFromPublic, model.OrderFromMembers:
		payapi.CreateOrderParsed(tx, w, r, order)
	case model.OrderInPerson:
		posapi.CreateOrderParsed(tx, w, r, order)
	default:
		api.BadRequestError(tx, w, "invalid source")
	}
}
