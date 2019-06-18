// Package stripe contains code for accessing the Stripe API.
package stripe

import (
	"log"
	"strconv"
	"strings"

	sapi "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/paymentintent"
	"github.com/stripe/stripe-go/paymentmethod"
	"github.com/stripe/stripe-go/terminal/connectiontoken"

	"scholacantorum.org/orders/config"
	"scholacantorum.org/orders/model"
)

// ChargeCard charges the user's card specified in the payment.  If the charge
// succeeds, the payment Subtype, Method, and Stripe fields are updated, and the
// function returns (true, card, ""), where card is the Stripe fingerprint for
// the card used.  If the charge is declined, the function returns false, "",
// and the decline message.  If the charge fails due to a Stripe API error, the
// function returns false and two empty strings.
func ChargeCard(order *model.Order, pmt *model.Payment) (success bool, card, cardError string) {
	var (
		iparams *sapi.PaymentIntentParams
		method  *sapi.PaymentMethod
		intent  *sapi.PaymentIntent
		charge  *sapi.Charge
		err     error
	)
	sapi.LogLevel = 1 // log only errors
	sapi.Key = config.Get("stripeSecretKey")
	iparams = &sapi.PaymentIntentParams{
		Amount:             sapi.Int64(int64(pmt.Amount)),
		Confirm:            sapi.Bool(true),
		Currency:           sapi.String(string(sapi.CurrencyUSD)),
		ConfirmationMethod: sapi.String(string(sapi.PaymentIntentConfirmationMethodManual)),
		Params: sapi.Params{
			Metadata: map[string]string{"order-number": strconv.Itoa(int(order.ID))},
		},
	}
	if strings.HasPrefix(pmt.Method, "pm_") {
		iparams.PaymentMethod = &pmt.Method
	} else {
		if method, err = paymentmethod.New(&sapi.PaymentMethodParams{
			Type: sapi.String(string(sapi.PaymentMethodTypeCard)),
			Card: &sapi.PaymentMethodCardParams{
				Token: &pmt.Method,
			},
		}); err != nil {
			log.Printf("ERROR: can't create payment method from token: %s", err)
			return false, "", ""
		}
		iparams.PaymentMethod = &method.ID
	}
	intent, err = paymentintent.New(iparams)
	if serr, ok := err.(*sapi.Error); ok && serr.Type == sapi.ErrorTypeCard {
		return false, "", serr.Msg
	}
	if err != nil {
		log.Printf("ERROR: can't create payment intent for order %d: %s", order.ID, err)
		return false, "", ""
	}
	if intent.Status != sapi.PaymentIntentStatusSucceeded {
		log.Printf("ERROR: payment intent for order %d is in invalid status %s", order.ID, intent.Status)
		return false, "", ""
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
	return true, charge.PaymentMethodDetails.Card.Fingerprint, ""
}

// CapturePayment captures a card-present payment that has already been
// authorized by the Stripe Terminal SDK.  If the capture succeeds, the payment
// Subtype, Method, and Stripe fields are updated.
func CapturePayment(order *model.Order, pmt *model.Payment) (card string, err error) {
	var (
		intent *sapi.PaymentIntent
		chg    *sapi.Charge
	)
	sapi.LogLevel = 1 // log only errors
	sapi.Key = config.Get("stripeSecretKey")
	if intent, err = paymentintent.Capture(pmt.Stripe, &sapi.PaymentIntentCaptureParams{}); err != nil {
		return "", err
	}
	chg = intent.Charges.Data[0]
	pmt.Stripe = chg.ID
	pmt.Method = brandMap[chg.PaymentMethodDetails.CardPresent.Brand]
	if pmt.Method == "" {
		pmt.Method = "card "
	}
	pmt.Method += chg.PaymentMethodDetails.CardPresent.Last4
	pmt.Subtype = chg.PaymentMethodDetails.CardPresent.ReadMethod
	card = chg.PaymentMethodDetails.CardPresent.Fingerprint
	return card, nil
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

// GetConnectionToken returns a Stripe Terminal connecton token, allowing a
// terminal to connect to our Stripe account.  It returns an empty string on
// failure.
func GetConnectionToken() string {
	var (
		token *sapi.TerminalConnectionToken
		err   error
	)
	sapi.LogLevel = 1 // log only errors
	sapi.Key = config.Get("stripeSecretKey")
	if token, err = connectiontoken.New(&sapi.TerminalConnectionTokenParams{}); err != nil {
		log.Printf("ERROR: can't get terminal connection token: %s", err)
		return ""
	}
	return token.Secret
}

// CreatePaymentIntent creates a payment intent for an order, so that it can be
// paid through Stripe Terminal.  It updates the Method field of the payment to
// contain the payment intent ID.  It returns true if successful, false on
// failure.
func CreatePaymentIntent(order *model.Order) bool {
	var (
		intent *sapi.PaymentIntent
		err    error
	)
	sapi.LogLevel = 1 // log only errors
	sapi.Key = config.Get("stripeSecretKey")
	intent, err = paymentintent.New(&sapi.PaymentIntentParams{
		Amount:             sapi.Int64(int64(order.Payments[0].Amount)),
		CaptureMethod:      sapi.String(string(sapi.PaymentIntentCaptureMethodManual)),
		Currency:           sapi.String(string(sapi.CurrencyUSD)),
		PaymentMethodTypes: sapi.StringSlice([]string{"card_present", "card"}),
		Params: sapi.Params{
			Metadata: map[string]string{"order-number": strconv.Itoa(int(order.ID))},
		},
	})
	if err != nil {
		log.Printf("ERROR: can't create payment intent for order %d: %s", order.ID, err)
		return false
	}
	order.Payments[0].Method = intent.ClientSecret
	order.Payments[0].Stripe = intent.ID
	return true
}

// CancelPaymentIntent cancels the payment intent for an order.  It is used when
// cleaning up after an order processing failure.
func CancelPaymentIntent(id string) error {
	var (
		err error
	)
	sapi.LogLevel = 1 // log only errors
	sapi.Key = config.Get("stripeSecretKey")
	_, err = paymentintent.Cancel(id, nil)
	return err
}

// GetCardFingerprint returns the fingerprint of the card used for the specified
// Stripe charge.  It returns an empty string for any error.
func GetCardFingerprint(chargeID string) string {
	var (
		chg *sapi.Charge
		err error
	)
	sapi.LogLevel = 1 // log only errors
	sapi.Key = config.Get("stripeSecretKey")
	if chg, err = charge.Get(chargeID, nil); err != nil {
		log.Printf("ERROR retrieving Stripe charge %s: %s", chargeID, err)
		return ""
	}
	switch chg.PaymentMethodDetails.Type {
	case sapi.ChargePaymentMethodDetailsTypeCard:
		return chg.PaymentMethodDetails.Card.Fingerprint
	case sapi.ChargePaymentMethodDetailsTypeCardPresent:
		return chg.PaymentMethodDetails.CardPresent.Fingerprint
	default:
		return ""
	}
}
