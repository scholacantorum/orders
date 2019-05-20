// Package stripe contains code for accessing the Stripe API.
package stripe

import (
	sapi "github.com/stripe/stripe-go"

	"scholacantorum.org/orders/private"
)

func init() {
	sapi.LogLevel = 1 // log only errors
	sapi.Key = private.StripeSecretKey
}
