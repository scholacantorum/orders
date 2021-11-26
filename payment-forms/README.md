# Schola Cantorum Ticket Order Form

On the Schola Cantorum public web site, any page that sells event tickets has a
`<script>` element at the point where the ticket purchase button should appear.
That script element invokes this Vue component, and specifies the products to be
sold through the ticket purchase button.

```html
<script src="https://orders.scholacantorum.org/payment-forms?product=ticket-2019-07-08"></script>
```

In othe

This form is invoked by the Schola Cantorum public web site, through a
`<script>` tag.

## How to Test Locally

1. Install the desired database at ~/src/orders/server/data/orders.db.
2. Start the orders server: cd ~/src/orders/server; go run ./cmd/serve.
3. Start the payment forms server: cd ~/src/orders/payment-forms; yarn serve.
4. Set the local web server configuration to use the local orders server:
   edit ~/src/schola6p/public-site-framework/config.yaml and set
   params.ordersURL to http://localhost:8200/.
5. Change the orders assets in
   ~/src/schola6p/public-site-framework/data/development/resources/buy-tickets.json
   to reference buy-tickets.js (without a fingerprint) only. No chunk-vendors
   and no css.
6. Start the local web server: cd ~/src/schola6p/public-site-framework;
   hugo serve.
