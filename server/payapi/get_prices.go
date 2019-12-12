package payapi

//go:generate easyjson -lower_camel_case -omit_empty get_prices.go

import (
	"errors"
	"fmt"
	"time"

	"github.com/mailru/easyjson/jwriter"

	"scholacantorum.org/orders/api"
	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/config"
	"scholacantorum.org/orders/model"
)

type getPricesData struct {
	id      model.ProductID
	name    string
	message string
	price   int
}

// GetPrices returns the prices and availability of one or more products.  It is
// used to drive payment forms.
func GetPrices(r *api.Request) (err error) {
	var (
		source      model.OrderSource
		productIDs  []string
		coupon      string
		couponMatch bool
		pdata       []*getPricesData
		product     *model.Product
		masterSKU   *model.SKU
		jw          *jwriter.Writer
	)
	// Get the request source and authorization.
	switch source = model.OrderSource(r.FormValue("source")); source {
	case "":
		source = model.OrderFromPublic
	case model.OrderFromPublic:
		// no-op
	case model.OrderFromMembers:
		if err = auth.GetSessionMembersAuth(r, r.FormValue("auth")); err != nil {
			return err
		}
	default:
		return errors.New("invalid source")
	}
	// Read the request details from the request.
	productIDs = r.Form["p"]
	if coupon = r.FormValue("coupon"); coupon == "" {
		couponMatch = true
	}
	// Look up the prices for each product.
	for _, pid := range productIDs {
		var (
			sku *model.SKU
			pd  getPricesData
		)
		// Get the product.  Skip nonexistent ones.
		if product = r.Tx.FetchProduct(model.ProductID(pid)); product == nil {
			continue
		}
		pd.id = product.ID
		pd.name = product.ShortName
		// Find the best SKU for this product.
		for _, s := range product.SKUs {
			if !api.MatchingSKU(s, coupon, source, true) {
				continue
			}
			if s.Coupon == coupon {
				couponMatch = true
			}
			sku = api.BetterSKU(s, sku)
		}
		if sku == nil {
			// No applicable SKUs; ignore product.
			continue
		}
		// Generate the product data to return.
		if !api.ProductHasCapacity(r, product) {
			pd.message = "This event is sold out."
		} else if pd.message = noSalesMessage(sku); pd.message == "" {
			pd.price = sku.Price
		}
		pdata = append(pdata, &pd)
		// The masterSKU determines the message to be shown in lieu of
		// the purchase button if none of the products are available for
		// sale.
		masterSKU = api.BetterSKU(sku, masterSKU)
	}
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jw = new(jwriter.Writer)
	if len(pdata) == 0 {
		// No products available for sale, or even with messages to
		// display, so we return a null.
		jw.RawString("null")
	} else if message := noSalesMessage(masterSKU); message != "" {
		// No products available for sale, but we do have a message to
		// display.
		jw.String(message)
	} else {
		// Return the product data.
		emitGetPrices(jw, r.Session, couponMatch, pdata)
	}
	jw.DumpTo(r)
	return nil
}

// noSalesMessage returns the string describing why the SKU isn't available for
// sale, or an empty string if it is available.
func noSalesMessage(sku *model.SKU) string {
	if sku.InSalesRange(time.Now()) < 0 {
		return fmt.Sprintf("Sales start on %s.", sku.SalesStart.Format("January\u00A02"))
	}
	return ""
}

//easyjson:json
type getPricesResponse struct {
	User            *getPricesResponseSession
	StripePublicKey string
	Coupon          bool
	Products        []*getPricesData
}

//easyjson:json
type getPricesResponseSession struct {
	ID       int
	Username string
	Name     string
	Email    string
	Address  string
	City     string
	State    string
	Zip      string
}

// emitGetPrices writes the JSON response.
func emitGetPrices(jw *jwriter.Writer, session *model.Session, couponMatch bool, pdata []*getPricesData) {
	var resp = getPricesResponse{
		Coupon:   couponMatch,
		Products: pdata,
	}
	if session != nil && session.Name != "" {
		resp.User = &getPricesResponseSession{
			ID:       session.Member,
			Username: session.Username,
			Name:     session.Name,
			Email:    session.Email,
			Address:  session.Address,
			City:     session.City,
			State:    session.State,
			Zip:      session.Zip,
		}
		resp.StripePublicKey = config.Get("stripePublicKey")
	}
	encodeGetPricesResponse(jw, resp)
}
