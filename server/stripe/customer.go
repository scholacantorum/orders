// Package stripe contains code for accessing the Stripe API.
package stripe

import (
	"log"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/setupintent"

	"scholacantorum.org/orders/config"
	"scholacantorum.org/orders/model"
)

// CreateCustomer creates a customer with the specified name and email, and the
// specified card source as its default for payments.  It returns the customer
// ID, payment method ID for future payments, and card description if
// successful, and a problem string if there's an issue with the card.
func CreateCustomer(name, email, card string) (custid, pmtmeth, desc, problem string) {
	var (
		cust *stripe.Customer
		si   *stripe.SetupIntent
		err  error
	)
	stripe.LogLevel = 1 // log only errors
	stripe.Key = config.Get("stripeSecretKey")
	var cparams = stripe.CustomerParams{
		Description: &name,
		Email:       &email,
	}
	if cust, err = customer.New(&cparams); err != nil {
		log.Printf("stripe create customer: %s", err)
		return
	}
	var siparams = stripe.SetupIntentParams{
		Confirm:       stripe.Bool(true),
		Customer:      &cust.ID,
		PaymentMethod: &card,
		Usage:         stripe.String(string(stripe.SetupIntentUsageOffSession)),
	}
	siparams.AddExpand("payment_method")
	if si, err = setupintent.New(&siparams); err != nil {
		if serr, ok := err.(*stripe.Error); ok && serr.Type == stripe.ErrorTypeCard {
			return "", "", "", serr.Msg
		}
		log.Printf("stripe create setup intent: %s", err)
		return
	}
	if si.Status != stripe.SetupIntentStatusSucceeded {
		return "", "", "", "card not authorized for delayed charges"
	}
	desc = brandMap[si.PaymentMethod.Card.Brand]
	if desc == "" {
		desc = "card "
	}
	desc += si.PaymentMethod.Card.Last4
	return cust.ID, si.PaymentMethod.ID, desc, ""
}

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

// UpdateCustomer updates the name, email, and optionally payment card of an
// existing customer.
func UpdateCustomer(id, name, email, card string) (desc string, err error) {
	stripe.LogLevel = 1 // log only errors
	stripe.Key = config.Get("stripeSecretKey")
	cparams := stripe.CustomerParams{Description: &name, Email: &email}
	if _, err = customer.Update(id, &cparams); err != nil {
		log.Printf("stripe update customer: %s", err)
		return "", err
	}
	if card == "" {
		return "", nil
	}
	siparams := stripe.SetupIntentParams{
		Confirm:       stripe.Bool(true),
		Customer:      &id,
		PaymentMethod: &card,
		Usage:         stripe.String(string(stripe.SetupIntentUsageOffSession)),
	}
	siparams.AddExpand("payment_method")
	var si *stripe.SetupIntent
	if si, err = setupintent.New(&siparams); err != nil {
		log.Printf("stripe create setup intent: %s", err)
		if ce, ok := err.(*stripe.Error); ok && ce.Type == stripe.ErrorTypeCard {
			return "", CardError(ce.Msg)
		}
		return "", err
	}
	if si.Status != stripe.SetupIntentStatusSucceeded {
		log.Printf("card not authorized for delayed charges")
		return "", CardError("card not authorized for delayed charges")
	}
	desc = brandMap[si.PaymentMethod.Card.Brand]
	if desc == "" {
		desc = "card "
	}
	desc += si.PaymentMethod.Card.Last4
	return desc, nil
}
