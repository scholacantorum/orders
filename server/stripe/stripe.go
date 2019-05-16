// Package stripe contains code for accessing the Stripe API.
package stripe

import (
	sapi "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/product"
	"github.com/stripe/stripe-go/sku"

	"scholacantorum.org/orders/private"
)

func init() {
	sapi.LogLevel = 1 // log only errors
	sapi.Key = private.StripeSecretKey
}

// ValidateProduct verifies that the specified ID is the ID of a usable Stripe
// product.
func ValidateProduct(id string) bool {
	var (
		prod *sapi.Product
		err  error
	)
	if prod, err = product.Get(id, nil); err != nil {
		return false
	}
	if prod.Type != sapi.ProductTypeGood {
		return false
	}
	if !prod.Active {
		return false
	}
	return true
}

// ValidateSKU verifies that the specified SKU ID is the ID of a usable Stripe
// SKU for the product with the specified product ID.  It returns the price of
// the SKU (in cents) and the coupon code of the SKU from Stripe.  If the ID is
// not valid, it returns a negative price.
func ValidateSKU(productID, skuID string) (price int, coupon string) {
	var (
		s   *sapi.SKU
		err error
	)
	if s, err = sku.Get(skuID, nil); err != nil {
		return -1, ""
	}
	if !s.Active {
		return -1, ""
	}
	if s.Attributes != nil {
		coupon = s.Attributes["coupon"]
	}
	return int(s.Price), coupon
}
