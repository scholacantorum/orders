# Redesigning the Schola Ordering System

Moving into the 2019–20 season, we want to make three major changes to our
ordering system:

* We want to abandon all use of paper ticketing.

  This will save us hundreds of dollars each year in ticket printing and mailing
  costs, and bring us in line with other arts organizations that have done the
  same.

* We want to track actual ticket usage.

  This will reduce the (admittedly already low) chance of ticket fraud, enable
  us to prevent overselling a limited event like the gala, and most importantly,
  give us better analytics of ticket usage, both statistically (e.g. which
  seatings are more popular?) and individually (which concerts did this donor
  attend?).

* We want to use Stripe rather than Square for at-the-door credit card sales.

  This allows us to do all of our credit card processing through a single
  processor, which simplifies bookkeeping and management.  We also are less
  comfortable with Square after it lost all of our at-the-door sales for the
  November 3, 2018 concert.

Our current ordering system, instituted at the beginning of the 2018–19 season,
cannot be readily extended to meet these goals.  This is primarily because it
has no underlying database that we control.  It records all purchases in Stripe
and in a Google Docs spreadsheet, but neither of these records are suitable for
ticket usage tracking.  Stripe's database is also limited to purchases made
through Stripe; it cannot be used to track cash or check sales.

Another problem is that Stripe introduced a new API to handle card-present sales
(e.g. at-the-door sales).  The new API can be used for all types of sales, and
for simplicity we probably should.  But the new API does not support tracking of
products, SKUs, and order details; it only deals in charges.  So order and
inventory tracking will need to be added to our code.

It should be noted that I looked for, but not found, off-the-shelf software or
commercial services that would do what we want.  We will either need to change
our practices (e.g. by eliminating Flex Passes and adding per-purchase service
fees), or we will need to implement our own system.  I choose the latter.

## Types of Sales

Our ordering system has to deal with a number of different sales paths and
product ordering styles.  These include:

* Public web site:
  * One-off donations:  card not present, immediate charge, variable amount.
  * Recurring donations:  card not present, delayed charge, variable amount.
  * Individual event ticket:  card not present, immediate charge, fixed amount,
    single use, multiple quantity.
  * Flex Pass ticket (to season or summer sings):  card not present,
    immediate charge, fixed amount, multiple use, multiple quantity.
* Members web site:
  * Concert recording:  card not present, immediate charge, fixed amount,
    perpetual use.
  * Sheet music:  card not present, immediate charge, fixed amount.
* Gala software:
  * Registration:  card not present *or* card present, immediate charge but save
    card for delayed charges, fixed amount, multiple quantity.
  * Fund-a-Need:  card present, card not present/immediate charge, or card not
    present/delayed charge; fixed amount, multipile quantity.
  * Auction item: card present, card not present/immediate charge, or card not
    present/delayed charge; variable amount.
* Order management web site:
  * All of the above (except recurring donation), paid via cash or check or
    other.
* Door sales app:
  * Individual event ticket: cash, check, or card present, immediate charge,
    fixed amount, single use, multiple quantity.
  * Flex Pass ticket: cash, check, or card present, immediate charge, fixed
    amount, multiple use, multiple quantity.

Therefore, it has to handle:

* Payment by cash, check, card, or other (office note).
* For cards, payment with card present, card not present, or saved card.
* For cards, save for later use or not.
* Fixed and variable amount products.
* Single use and multiple use products.
* Quantity-always-1 and variable-quantity products.

## Back End Technology Choices

**Authentication and Authorization Model**  
We need username and password storage with encryption, password quality
enforcement, profile management for changing passwords, "forgot my password"
recovery functionality, automatic lockout after successive failed attempts, etc.
The members web site already has all of this, and I see no value in
reimplementing it.  On the other hand, I don't want to couple the two sites
together very closely by having them share a database or share code.  So my plan
is that the ordering system will delegate authentication and authorization to
the members site via APIs.

**Database**  
All of the existing Schola sites use SQLite 3, so that's what I plan to use for
the ordering system.  MySQL is another alternative, but it offers nothing
additional that we need and nothing that would justify the maintenance cost of
using two different database systems.

**Hosting Domain**  
I plan to host the ordering system on orders.scholacantorum.org, which is part
of Schola's free web hosting account at Dreamhost.  It doesn't allow constantly
running servers, so the ordering system will be invoked as CGI scripts, but our
system will be very low usage so that shouldn't be a problem.  I'm not hosting
it on the main scholacantorum.org because that is a static site generated by
Hugo and it's not trivial to introduce CGI scripts to it.

**Implementation Language**  
All of the existing Schola sites have their back end code written in Go, so
that's what I plan to use for the ordering system.  PHP and Python would be
viable alternatives, but they don't provide any additional functionality that we
need, and there's nothing to justify the maintenance cost of using multiple
back-end languages.

**Stripe API Usage**  
We have to use Stripe's new `PaymentIntent` and `PaymentMethod` objects to
process card-present transactions.  Since those work for all types of
transactions, I plan to use them for all transactions.  The alternative would be
to continue to use the older `Source` objects for card-not-present transactions,
but that just makes the code more complex.  Note that we will be abandoning use
of Stripe's tracking of `Order`, `Product` and `SKU` objects, which
`PaymentIntent` and `PaymentMethod` don't support.  But we're going to need to
track our orders, products, and SKUs on our own server anyway, to encompass
non-credit-card orders, so that's not a bit loss.

**Stripe Customer Tracking**  
Transactions with delayed charges, such as recurring donations or gala auction
purchases, must be associated with a Stripe `Customer` object, where the payment
method is stored until we decide to charge it.  My plan is to use Stripe
`Customer` objects only for those transactions, and do all other customer
tracking in our own database.  We're going to need to track customers in our own
database anyway, so there's little value in doing it redundantly in Stripe,
except for the few cases where it's required.

## Front End Technology Choices

For the ordering system, we will need five new front ends.

**Order Management Web Site**  
Used by the office staff to record offline orders and review analytics.  This
will be a single-page application using the Nuxt framework.

**Order Dialog**  
Customer-facing dialog form for entering and confirming order details; the
various Schola web sites will bring up this dialog in an `<iframe>` to handle
order placement.  This will be a Vue-based webapp.

**Ticket Summary**  
Customers who scan the QR code on their ticket will be shown this site, which
will give details of the ticket and the order through which it was purchased.
This will be a static page served by the back end.

**Door Sales App**  
Used by the front-of-house staff at a concert to record at-the-door ticket
sales, both cash and credit card.  In order to handle card-present transactions,
this will be implemented as a React Native application and run on an iOS device.
(Android support is coming but not available yet.)

**Ticket Taking App**  
Used by the front-of-house staff at a concert to scan the QR codes on tickets,
validate them, and mark them as used.  This needs to run on a mobile phone, to
use the camera.  It will probably be implemented as a Vue-based webapp.  But if
it turns out to be convenient to include it in the same React Native application
used for door sales, I will do that instead.

## Block Diagram

```x
┌───────────────┐   ┌──────────────────┐   ┌─────────────────┐   ┌──────────────────┐
│ Gala Software │   │ Order Management │   │ Public Web Site │   │ Members Web Site │
└──────┬─┬──────┘   └───────┬──┬───────┘   └────────┬────────┘   └─────────┬────────┘
       │ └──────────────────│──┴──────────────────┬─┴──────────────────────┘
       └────────┐ ┌─────────┘                     │
┌────────────┐  │ │  ┌────────────────┐   ┌───────┴──────┐   ┌────────────────┐   ┌─────────────┐
│ Ticket App │  │ │  │ Ticket Summary │   │ Order Dialog │   │ Door Sales App ├───┤ Card Reader │
└──────┬─────┘  │ │  └────────┬───────┘   └─────┬──┬─────┘   └──────┬──┬──────┘   └──────┬──────┘
       │        │ │           │                 │  └────────────────│──┴──┬──────────────┘
       └────────┴─┴───────┬───┴─────────────────┴───────────────────┘     │
┌─────────────────────────┴─────────────────────────┬──────────┐          │
│                 Back End Software                 │ Database │          │
└───┬───────────────┬────────────────────────────┬──┴──────────┘          │
┌───┴───┐   ┌───────┴──────┐              ┌──────┴────────────────────────┴─────────────────────┐
│ Email │   │ Google Sheet │              │                        Stripe                       │
└───────┘   └──────────────┘              └─────────────────────────────────────────────────────┘
```

The Stripe communications deserve a bit more detail.  When the order dialog
shows a credit card entry field, that field is not actually coming from our site
code at all; it is in an IFRAME tag, and is coming from Stripe's own servers.
Anything the user types in it is sent directly to Stripe and is never seen by
our code.  What Stripe gives our site's code is a "token" that represents the
user input without exposing it.  The order dialog passes that to the back end,
which in turn passes it to Stripe's API.  Stripe then finds the real card data
associated with the token, and charges that card.  A similar process happens
with the door sales app and the card reader: the card reader sends the card data
directly to Stripe, and gives the corresponding token to our phone app.

## Data Model and Database Schema

```x
              ┌───────────────────────┐      ? = zero or one
              │        orderT         │      * = zero or more
              └───┬────────────────┬──┘      1 = exactly one
                 +│               *│         + = one or more
           ┌──────┴─────┐     ┌────┴────┐
           │ order_line │     │ payment │
           └──┬───┬─────┘     └─────────┘
   *┌─────────┘   │1  
┌───┴────┐  ┌─────┴────┐
│ ticket │  │ product  │
└───┬────┘  └─┬─────┬──┘
   ?│         │    +│
┌───┴────┐    │  ┌──┴──┐   ┌────────────┐
│ event  │    │  │ sku │   │ card_email │
└───┬────┘    │  └─────┘   └────────────┘
   +│        *│
┌───┴─────────┴─┐             ┌─────────┐
│ product_event │             │ session │
└───────────────┘             └─────────┘
```

There are database tables for order management and for ticket usage tracking.

### Ordering

The ordering part of the system has a pretty straightforward database schema,
which centers around products, SKUs, orders, and payments.

```sql
CREATE TABLE product (
    id           text    PRIMARY KEY,
    series       text,
    name         text    NOT NULL,
    shortname    text    NOT NULL,
    type         text    NOT NULL,
    receipt      text,
    ticket_count integer,
    ticket_class text
    options      text
);
CREATE TABLE sku (
    product     text     NOT NULL REFERENCES product,
    source      text     NOT NULL CHECK (source IN ('public', 'members', 'gala', 'inperson', 'office')),
    coupon      text,
    sales_start datetime,
    sales_end   datetime,
    price       integer
);
```

Products and SKUs are closely related; each product is a collection of one or
more SKUs.  Products represent what the customer gets (e.g., adult admission to
a particular concert) while SKUs represent how they got it (full price, coupon
code, early bird discount, etc.).

In the `product` table, products are arranged by `series` (e.g., season) and
`type` (`ticket`, `donation`, etc.).  Each product has two names:  a full `name`
that stands on its own, and a `shortname` used in contexts where the full name
is redundant.  If `receipt` is non-NULL, it is the template for a receipt email
that will be sent to the purchaser.  `ticket_count` is the number of (virtual)
event tickets the customer gets as a result of buying one unit of this product.
For a non-ticket product, this would be zero.  For a ticket to a specific event,
this would be one.  For a Flex Pass, this would be the number of events covered
by the pass (e.g. 4 for a season subscription, 6 for summer sings, etc.).
`ticket_class`, if specified, is the restriction on who is allowed to use the
ticket (e.g. "Senior" or "Student").  It is not specified for unrestricted-use
tickets or for non-ticket products.  Some products have options to choose from
(e.g., dinner entree or T-shirt size), in which case those are listed in the
`options`.

In the `sku` table, the `price` is the price of the item (as an integer in
cents, to avoid rounding problems).  If it is `NULL`, that means the SKU's price
is variable (e.g. a donation or an auction bid).  The `source` is the source
through which the order must be placed in order to get that price.  The `coupon`
is the coupon code that the customer has to enter in order to get that price.
Every product should have a sku with a `NULL` coupon to set the price when no
valid coupon code is entered.  The `sales_start` and `sales_end` columns, when
not `NULL`, specify the period of time during when regular customers can place
orders.  (Office staff can place orders for anything at any time.)

It should be noted that Schola has experimented with a wide variety of sales
incentives over the years: quantity discounts, percentage discounts, early bird
discounts, etc.  This schema only accounts for the sales models currently in
use.  Additional ones can be added on demand.

```sql
CREATE TABLE orderT (
    id        integer  PRIMARY KEY,
    token     text     UNIQUE,
    valid     boolean  NOT NULL DEFAULT false,
    source    text     NOT NULL CHECK (source IN ('public', 'members', 'gala', 'inperson', 'office')),
    name      text,
    email     text,
    address   text,
    city      text,
    state     text,
    zip       text,
    phone     text,
    customer  text,
    member    integer,
    created   datetime NOT NULL,
    cnote     text,
    onote     text,
    in_access boolean,
    coupon    text,
);
CREATE TABLE order_line (
    id          integer PRIMARY KEY,
    orderid     integer NOT NULL REFERENCES orders,
    product     text    NOT NULL,
    quantity    integer NOT NULL,
    price       integer NOT NULL,
    guest_name  text,
    guest_email text,
    option      text
);
```

In the `orderT` table — which has a "T" in the name because "order" is a
reserved word in SQL — both the `id` and the `token` are unique identifiers of
the order.  The `id` is a monotonic integer used for database access; the
`token` is an opaque string embedded in ticket QR codes.  The valid flag
indicates whether the order was actually placed and paid; it is false when an
order is in progress.  `source` is the source from which the order was placed.
`name`, `email`, `address`, `city`, `state`, `zip`, and `phone` are information
about the customer placing the order.  `customer` is the customer's Stripe ID,
if any.  `member` is the ID of the customer on the members site, if the customer
is a member.  `created` is the time the order was created.  `cnote` and `onote`
are notes on the order from the customer and the office, respectively.
`in_access` indicates whether information about the order has been transferred
into the office Access database.  `coupon` is the coupon code used to place the
order, if any.

Orders comprise one or more order lines, stored in the `order_line` table under
the same `orderid` and with unique `id` values.  Each line represents a purchase
of a `quantity` of a `product`, paying `price` (in cents) for each unit.  Thus,
the total amount for the line is always `quantity * price`.  For some products,
a guest name, guest email address, and/or options are recorded.

```sql
CREATE TABLE payment (
    id integer PRIMARY KEY,
    orderid integer NOT NULL REFERENCES orderT,
    type    text    NOT NULL CHECK (type IN ('card', 'card-present', 'check', 'cash', 'other')),
    subtype text,
    method  text,
    stripe  text,
    created datetime NOT NULL,
    initial boolean  NOT NULL,
    amount  integer  NOT NULL
)
```

In the `payment` table, there is one row for each payment or refund on an order.
The payment is identified by `id` and the order is identified by `orderid`.
`type` indicates the basic type of payment.  `subtype` is primarily used with
`card` and `card-present` to indicate how the card number was collected.
`method` is a textual description of the payment method, used in receipts; for
example, for `card` payments it includes the card type and last four digits.  If
the payment is processed through Stripe, `stripe` contains the charge ID or
refund ID.  `created` is the time when the payment or refund occurred.
`initial` is true for the first payment on an order, and false for all others.
`amount` is the amount of the payment (in cents); negative numbers are refunds.

```sql
CREATE TABLE card_email (
    card  text PRIMARY KEY,
    name  text,
    email text,
)
```

In the `card_email` table, Stripe card fingerprints (`card`) are mapped to the
`name` and/or `email` of the person who most recently placed an order using that
card.  This allows us to guess their name and email when they present that card
for subsequent orders.

### Ticket Usage Tracking

The ticket usage tracking part of the schema is more novel.  It involves events,
tickets, and admissions.

```sql
CREATE TABLE event (
    id         integer  PRIMARY KEY,
    members_id integer  UNIQUE,
    name       text     NOT NULL,
    series     text     NOT NULL,
    start      datetime NOT NULL,
    capacity   integer
);
CREATE TABLE product_event (
    product  integer NOT NULL REFERENCES product,
    event    integer NOT NULL REFERENCES event,
    priority integer NOT NULL,
    PRIMARY KEY (product, event)
);
```

The `event` table lists events for which tickets are sold and ticket usage is
tracked.  Note that a concert with two seatings gets two rows in this table.
The `members_id` is the ID of the event in the members site database, which may
be useful for cross references.  The `capacity`, if not `NULL`, is the capacity
of the house; ticket sales will be stopped when the maximum is reached.  (But
see note below regarding Flex Passes.)

The `product_event` join table indicates which products correspond to tickets
for which events.  A product may have zero rows in this table (it's something
other than an event ticket); one row in this table (it's a ticket to a specific
event); or multiple rows in this table (it's a Flex Pass that can be used for
multiple events).  The `priority` field helps determine which ticket on an order
should be consumed when presented at an event.  Typical usage is:

* 0 for individual tickets targeted at the event itself
* 10 for individual tickets targeted at the other seating of the same concert
* 20 for flex passes that include the event
* 30 for individual tickets for a different concert in the same season

```sql
CREATE TABLE ticket (
    id         integer  PRIMARY KEY,
    order_line integer  NOT NULL REFERENCES order_line,
    event      integer  REFERENCES event,
    used       datetime
);
```

The `ticket` table has an entry for every purchased ticket, i.e., for every
entry of a single person to a single event.  Each ticket has a unique `id`.  For
each `order_line` row, there will be `order_line.qty * product.ticket_count`
rows in the `ticket` table.  For example, if someone purchases two season
subscriptions to a season with four concert sets, there will be eight rows added
to the `ticket` table, each with a unique `id`.

The `event` column specifies the event for which the ticket was or will be used.
For Flex Pass tickets, this column remains `NULL` until the ticket is used or
until there is only one event remaining that it can be used for.  The `used`
column remains `NULL` until the ticket is used, and then records the exact time
of usage.

Note that the `capacity` for an event is measured against the number of entries
in the `ticket` table whose `event` columns are set to that event.  In other
words, outstanding Flex Pass tickets are *not counted* against the capacity,
unless there are no other events left at which they could be used.  It is still
possible to oversell the house.  For this reason, the `capacity` is really only
effective for events that are not covered by a Flex Pass.

## Back End APIs

The ordering / ticketing system needs APIs to support sales configuration, order
processing, ticket validation and recording of usage, and ticket usage
analytics.  The APIs are broken up into different endpoints based on their
usage.

### Office APIs

These APIs are used by the Schola Office webapp.

```x
POST /ofcapi/event       Create an event
POST /ofcapi/login       Authenticate
GET  /ofcapi/order/$id   Get details of an order
POST /ofcapi/product     Create a product
GET  /ofcapi/report      Run a report
```

The `login` API is used to log into the office webapp.  (Authentication is
delegated to the members web site.)  The `report` API is used to tabulate order
information, and the `order/$id` API gets information about a specific order.
The other two APIs, while implemented, are not yet used.

### Payment APIs

These APIs are used by the payment forms displayed on the various Schola web
sites for buying tickets, donating, registering for the gala, purchasing concert
recordings, etc.

```x
POST /payapi/order    Create an order
GET  /payapi/prices   Get pricing for product(s)
```

### Point-of-Sale APIs

These APIs are used by the webapp and iOS app used at the front of the house at
a concert for in-person sales and ticket scanning.

```x
GET    /posapi/event                      List all events
GET    /posapi/event/$id/orders           List orders for tickets to an event
GET    /posapi/event/$id/prices           Get pricing for tickets to an event
GET    /posapi/event/$id/ticket/$token    Use a ticket at an event
POST   /posapi/event/$id/ticket/$token    Use a ticket at an event
POST   /posapi/login                      Authenticate
POST   /posapi/order                      Create an order
DELETE /posapi/order/$id                  Delete an order
POST   /posapi/order/$id/capturePayment   Capture the payment for an order
POST   /posapi/order/$id/sendReceipt      Send a receipt for an order
GET    /posapi/stripe/connectTerminal     Get a Stripe terminal connection token
```

### Ticket-Taking GUI

This entrypoint actually serves a web page rather than an API.  It is the web
page people see if they scan the QR code on their ticket.  It shows the usage of
the ticket.

```x
GET /ticket/$token   Show information about a ticket
```
