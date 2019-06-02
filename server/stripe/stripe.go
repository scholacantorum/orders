// Package stripe contains code for accessing the Stripe API.
package stripe

import (
	"log"
	"strconv"

	sapi "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"

	"scholacantorum.org/orders/config"
	"scholacantorum.org/orders/model"
)

// ChargeCard charges the user's card specified in the payment.  If the charge
// succeeds, the payment Subtype, Method, and Stripe fields are updated, and the
// function returns (true, "").  If the charge is declined, the function returns
// false and the decline message.  If the charge fails due to a Stripe API
// error, the function returns false and an empty string.
func ChargeCard(order *model.Order, pmt *model.Payment) (success bool, cardError string) {
	var (
		intent *sapi.PaymentIntent
		charge *sapi.Charge
		err    error
	)
	sapi.LogLevel = 1 // log only errors
	sapi.Key = config.Get("stripeSecretKey")
	intent, err = paymentintent.New(&sapi.PaymentIntentParams{
		Amount:             sapi.Int64(int64(pmt.Amount)),
		Confirm:            sapi.Bool(true),
		Currency:           sapi.String(string(sapi.CurrencyUSD)),
		ConfirmationMethod: sapi.String(string(sapi.PaymentIntentConfirmationMethodManual)),
		Params: sapi.Params{
			Metadata: map[string]string{"order-number": strconv.Itoa(int(order.ID))},
		},
		PaymentMethod: sapi.String(pmt.Method),
	})
	if serr, ok := err.(*sapi.Error); ok && serr.Type == sapi.ErrorTypeCard {
		return false, serr.Msg
	}
	if err != nil {
		log.Printf("ERROR: can't create payment intent for order %d: %s", order.ID, err)
		return false, ""
	}
	if intent.Status != sapi.PaymentIntentStatusSucceeded {
		log.Printf("ERROR: payment intent for order %d is in invalid status %s", order.ID, intent.Status)
		return false, ""
	}
	charge = intent.Charges.Data[0]
	pmt.Stripe = charge.ID
	pmt.Method = brandMap[charge.PaymentMethodDetails.Card.Brand]
	if pmt.Method == "" {
		pmt.Method = "card "
	}
	pmt.Method += charge.PaymentMethodDetails.Card.Last4
	if charge.PaymentMethodDetails.Card.Wallet != nil {
		pmt.Subtype = string(charge.PaymentMethodDetails.Card.Wallet.Type)
	}
	return true, ""
}

var brandMap = map[sapi.PaymentMethodCardBrand]string{
	sapi.PaymentMethodCardBrandAmex:       "AmEx ",
	sapi.PaymentMethodCardBrandDiners:     "Diners ",
	sapi.PaymentMethodCardBrandDiscover:   "Discover ",
	sapi.PaymentMethodCardBrandJCB:        "JCB ",
	sapi.PaymentMethodCardBrandMastercard: "MasterCard ",
	sapi.PaymentMethodCardBrandUnionpay:   "UnionPay ",
	sapi.PaymentMethodCardBrandVisa:       "Visa ",
}
