// Package stripe contains code for accessing the Stripe API.
package stripe

import (
	"log"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"

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

// UpdateCustomer updates the name and email of an existing customer.
func UpdateCustomer(id, name, email string) (success bool) {
	stripe.LogLevel = 1 // log only errors
	stripe.Key = config.Get("stripeSecretKey")
	var cparams = stripe.CustomerParams{Description: &name, Email: &email}
	if _, err := customer.Update(id, &cparams); err != nil {
		log.Printf("stripe update customer: %s", err)
		return false
	}
	return true
}
