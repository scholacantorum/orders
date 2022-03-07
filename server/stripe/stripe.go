// Package stripe contains code for accessing the Stripe API.
package stripe

import (
	"log"
	"strconv"
	"strings"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/paymentintent"
	"github.com/stripe/stripe-go/paymentmethod"
	"github.com/stripe/stripe-go/terminal/connectiontoken"

	"scholacantorum.org/orders/config"
	"scholacantorum.org/orders/model"
)

// FindOrCreateCustomer finds a customer with the name and email in the
// specified order.  If no matching customer was found, it creates one.  It sets
// the customer ID in the order.
func FindOrCreateCustomer(order *model.Order) (err error) {
	var (
		cust   *stripe.Customer
		clistp *stripe.CustomerListParams
		iter   *customer.Iter
	)
	stripe.LogLevel = 1 // log only errors
	stripe.Key = config.Get("stripeSecretKey")
	// Look for an existing customer first.
	clistp = new(stripe.CustomerListParams)
	clistp.Filters.AddFilter("email", "", order.Email)
	iter = customer.List(clistp)
	for iter.Next() {
		c := iter.Customer()
		if c.Description != order.Name || c.Email != order.Email {
			continue
		}
		if c.Metadata["monthly-donation-amount"] != "" {
			continue
		}
		order.Customer = c.ID
		return nil
	}
	// Create a new customer if none was found.
	var cparams = stripe.CustomerParams{Description: &order.Name, Email: &order.Email}
	if cust, err = customer.New(&cparams); err != nil {
		log.Printf("stripe create customer: %s", err)
		return err
	}
	order.Customer = cust.ID
	return nil
}

// ChargeCard charges the user's card specified in the payment.  If the charge
// succeeds, the payment Subtype, Method, and Stripe fields are updated, and the
// function returns (true, card, ""), where card is the Stripe fingerprint for
// the card used.  If the charge is declined, the function returns false, "",
// and the decline message.  If the charge fails due to a Stripe API error, the
// function returns false and two empty strings.
func ChargeCard(order *model.Order, pmt *model.Payment) (success bool, card, cardError string) {
	var (
		iparams *stripe.PaymentIntentParams
		method  *stripe.PaymentMethod
		intent  *stripe.PaymentIntent
		charge  *stripe.Charge
		err     error
	)
	stripe.LogLevel = 1 // log only errors
	stripe.Key = config.Get("stripeSecretKey")
	iparams = &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(int64(pmt.Amount)),
		Confirm:            stripe.Bool(true),
		Currency:           stripe.String(string(stripe.CurrencyUSD)),
		ConfirmationMethod: stripe.String(string(stripe.PaymentIntentConfirmationMethodManual)),
		Params: stripe.Params{
			Metadata: map[string]string{"order-number": strconv.Itoa(int(order.ID))},
		},
	}
	if order.SaveForReuse && order.Customer != "" {
		iparams.Customer = &order.Customer
		iparams.SetupFutureUsage = stripe.String(string(stripe.PaymentIntentSetupFutureUsageOffSession))
	}
	if strings.HasPrefix(pmt.Method, "pm_") {
		iparams.PaymentMethod = &pmt.Method
		pmt.StripePM = pmt.Method
	} else {
		if method, err = paymentmethod.New(&stripe.PaymentMethodParams{
			Type: stripe.String(string(stripe.PaymentMethodTypeCard)),
			Card: &stripe.PaymentMethodCardParams{
				Token: &pmt.Method,
			},
		}); err != nil {
			log.Printf("ERROR: can't create payment method from token: %s", err)
			return false, "", ""
		}
		iparams.PaymentMethod = &method.ID
		pmt.StripePM = method.ID
	}
	intent, err = paymentintent.New(iparams)
	if serr, ok := err.(*stripe.Error); ok && serr.Type == stripe.ErrorTypeCard {
		return false, "", serr.Msg
	}
	if err != nil {
		log.Printf("ERROR: can't create payment intent for order %d: %s", order.ID, err)
		return false, "", ""
	}
	if intent.Status != stripe.PaymentIntentStatusSucceeded {
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
		intent *stripe.PaymentIntent
		chg    *stripe.Charge
	)
	stripe.LogLevel = 1 // log only errors
	stripe.Key = config.Get("stripeSecretKey")
	if intent, err = paymentintent.Capture(pmt.Stripe, &stripe.PaymentIntentCaptureParams{}); err != nil {
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

var brandMap = map[stripe.PaymentMethodCardBrand]string{
	stripe.PaymentMethodCardBrandAmex:       "AmEx ",
	stripe.PaymentMethodCardBrandDiners:     "Diners ",
	stripe.PaymentMethodCardBrandDiscover:   "Discover ",
	stripe.PaymentMethodCardBrandJCB:        "JCB ",
	stripe.PaymentMethodCardBrandMastercard: "MasterCard ",
	stripe.PaymentMethodCardBrandUnionpay:   "UnionPay ",
	stripe.PaymentMethodCardBrandVisa:       "Visa ",
}

// GetConnectionToken returns a Stripe Terminal connecton token, allowing a
// terminal to connect to our Stripe account.  It returns an empty string on
// failure.
func GetConnectionToken() string {
	var (
		token *stripe.TerminalConnectionToken
		err   error
	)
	stripe.LogLevel = 1 // log only errors
	stripe.Key = config.Get("stripeSecretKey")
	if token, err = connectiontoken.New(&stripe.TerminalConnectionTokenParams{}); err != nil {
		log.Printf("ERROR: can't get terminal connection token: %s", err)
		return ""
	}
	return token.Secret
}

// CreatePaymentIntent creates a payment intent for an order, so that it can be
// paid through Stripe Terminal.  It updates the Method field of the payment to
// contain the payment intent secret.  It returns true if successful, false on
// failure.
func CreatePaymentIntent(order *model.Order) bool {
	var (
		intent *stripe.PaymentIntent
		err    error
	)
	stripe.LogLevel = 1 // log only errors
	stripe.Key = config.Get("stripeSecretKey")
	intent, err = paymentintent.New(&stripe.PaymentIntentParams{
		Amount:             stripe.Int64(int64(order.Payments[0].Amount)),
		CaptureMethod:      stripe.String(string(stripe.PaymentIntentCaptureMethodManual)),
		Currency:           stripe.String(string(stripe.CurrencyUSD)),
		PaymentMethodTypes: stripe.StringSlice([]string{"card_present", "card"}),
		Params: stripe.Params{
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
	stripe.LogLevel = 1 // log only errors
	stripe.Key = config.Get("stripeSecretKey")
	_, err = paymentintent.Cancel(id, nil)
	return err
}

// GetCardFingerprint returns the fingerprint of the card used for the specified
// Stripe charge.  It returns an empty string for any error.
func GetCardFingerprint(chargeID string) string {
	var (
		chg *stripe.Charge
		err error
	)
	stripe.LogLevel = 1 // log only errors
	stripe.Key = config.Get("stripeSecretKey")
	if chg, err = charge.Get(chargeID, nil); err != nil {
		log.Printf("ERROR retrieving Stripe charge %s: %s", chargeID, err)
		return ""
	}
	switch chg.PaymentMethodDetails.Type {
	case stripe.ChargePaymentMethodDetailsTypeCard:
		return chg.PaymentMethodDetails.Card.Fingerprint
	case stripe.ChargePaymentMethodDetailsTypeCardPresent:
		return chg.PaymentMethodDetails.CardPresent.Fingerprint
	default:
		return ""
	}
}
