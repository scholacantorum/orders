package api

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rothskeller/json"

	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
)

// CalculateOrder handles GET /api/order requests.  It takes a JSON order in its
// request body, and returns the same order with amounts filled in.  The order
// is not validated except as necessary to compute prices.
func CalculateOrder(tx db.Tx, w http.ResponseWriter, r *http.Request) {
	var (
		session *model.Session
		order   *model.Order
		privs   model.Privilege
		err     error
	)
	// Get current session privileges.
	if auth.HasSession(r) {
		if session = auth.GetSession(tx, w, r, 0); session == nil {
			return
		}
		privs = session.Privileges
	}
	// Read the order details from the request.
	if order, err = parseCalculateOrder(r.Body); err != nil {
		BadRequestError(tx, w, err.Error())
		return
	}
	// Resolve the SKUs and fill in the prices.
	if !resolveSKUs(tx, order, privs, false) {
		BadRequestError(tx, w, "invalid products or prices")
		return
	}
	// Send the result back.
	commit(tx)
	w.Header().Set("Content-Type", "application/json")
	emitCalculatedOrder(w, order)
}

// parseCalculateOrder parses the input order incompletely: just enough to be
// able to calculate the order amounts.
func parseCalculateOrder(r io.Reader) (o *model.Order, err error) {
	var jr = json.NewReader(r)

	o = new(model.Order)
	err = jr.Read(json.ObjectHandler(func(key string) json.Handlers {
		switch key {
		case "coupon":
			return json.StringHandler(func(s string) { o.Coupon = s })
		case "lines":
			return json.ArrayHandler(func() json.Handlers {
				var ol model.OrderLine
				o.Lines = append(o.Lines, &ol)
				return json.ObjectHandler(func(key string) json.Handlers {
					switch key {
					case "product":
						return json.StringHandler(func(s string) {
							ol.Product = &model.Product{ID: model.ProductID(s)}
						})
					case "quantity":
						return json.IntHandler(func(i int) { ol.Quantity = i })
					default:
						return json.IgnoreHandler()
					}

				})
			})
		default:
			return json.IgnoreHandler()
		}
	}))
	return o, err
}

// resolveSKUs walks through each line of the order, finding the listed product
// and computing the amount of the order line, following the SKU rules
// documented in db/schema.sql.  If mustMatch is true, it requires that the line
// have the correct amount.  If mustMatch is false, it overrides the line
// amount, setting it to the correct amount. It returns true if everything
// resolved successfully and false otherwise.
func resolveSKUs(tx db.Tx, order *model.Order, privs model.Privilege, mustMatch bool) bool {
	var couponMatch bool

	for _, line := range order.Lines {
		var sku *model.SKU
		var qty int
		var amount int
		var count int

		if qty = line.Quantity; qty < 1 {
			return false
		}
		if line.Product = tx.FetchProduct(line.Product.ID); line.Product == nil {
			return false
		}
		if line.Product.Type == model.ProdAuctionItem || line.Product.Type == model.ProdDonation {
			if line.Amount < 1 {
				return false
			}
			continue
		}
		for qty > 0 {
			for _, s := range line.Product.SKUs {
				if !matchSKU(order, privs, qty, s) {
					continue
				}
				if s.Coupon != "" {
					couponMatch = true
				}
				if sku == nil || betterSKU(s, sku) {
					sku = s
				}
			}
			if sku == nil {
				return false
			}
			count = line.Quantity / sku.Quantity
			amount += sku.Price * count
			qty -= count * sku.Quantity
		}
		if mustMatch && line.Amount != amount {
			return false
		}
		line.Amount = amount
	}
	if !couponMatch {
		order.Coupon = ""
	}
	return true
}

// matchSKU returns true if all of the criteria for the SKU are met by the order
// being placed.
func matchSKU(order *model.Order, privs model.Privilege, qty int, sku *model.SKU) bool {
	if sku.MembersOnly && privs == 0 {
		return false
	}
	if sku.Coupon != "" && !strings.EqualFold(sku.Coupon, order.Coupon) {
		return false
	}
	if privs&model.PrivHandleOrders != 0 {
		return true
	}
	if !sku.SalesStart.IsZero() && sku.SalesStart.After(time.Now()) {
		return false
	}
	if !sku.SalesEnd.IsZero() && sku.SalesEnd.Before(time.Now()) {
		return false
	}
	if qty < sku.Quantity {
		return false
	}
	return true
}

// betterSKU returns whether sku1 is a higher priority match than sku2.
func betterSKU(sku1, sku2 *model.SKU) bool {
	if sku1.MembersOnly && !sku2.MembersOnly {
		return true
	}
	if !sku1.MembersOnly && sku2.MembersOnly {
		return false
	}
	if sku1.Coupon != "" && sku2.Coupon == "" {
		return true
	}
	if sku1.Coupon == "" && sku2.Coupon != "" {
		return false
	}
	if sku1.Quantity > sku2.Quantity {
		return true
	}
	if sku1.Quantity < sku2.Quantity {
		return false
	}
	if (sku1.SalesStart.IsZero() || !sku1.SalesStart.After(time.Now())) &&
		(sku2.SalesEnd.IsZero() || !sku2.SalesEnd.Before(time.Now())) {
		return true
	}
	return false
}

// emitCalculatedOrder emits the JSON with the calculated order line amounts.
// If the coupon code was recognized, it is emitted as well, to confirm to the
// client that it is valid.
func emitCalculatedOrder(w io.Writer, o *model.Order) {
	var jw = json.NewWriter(w)
	jw.Object(func() {
		if o.Coupon != "" {
			jw.Prop("coupon", o.Coupon)
		}
		jw.Prop("lineAmounts", func() {
			jw.Array(func() {
				for _, ol := range o.Lines {
					jw.Int(ol.Amount)
				}
			})
		})
	})
	jw.Close()
}
