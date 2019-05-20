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
Hugo and it's not trivial to introduce CGI scripts to it.  I'm not hosting it on
scholacantorummembers.org (or a subdomain of that), because that domain is part
of my personal web hosting account at Dreamhost, and it doesn't seem appropriate
to put Schola's order processing there.

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
              │        order          │      * = zero or more
              └───┬────────────────┬──┘      1 = exactly one
                 +│               *│         + = one or more
           ┌──────┴─────┐     ┌────┴────┐
           │ order_line │     │ payment │
           └──┬───┬───┬─┘     └────┬────┘
   *┌─────────┘   │1  └──────┐*   +│
┌───┴────┐  ┌─────┴────┐  ┌──┴─────┴─────┐
│ ticket │  │ product  │  │ payment_line │
└───┬────┘  └─┬─────┬──┘  └──────────────┘
   ?│         │    +│
┌───┴────┐    │  ┌──┴──┐
│ event  │    │  │ sku │
└───┬────┘    │  └─────┘
   +│        *│
┌───┴─────────┴─┐              ┌─────────┐
│ product_event │              │ session │
└───────────────┘              └─────────┘
```

There are database tables for order management and for ticket usage tracking.
(There are also tables for auditing and logging, but those are not detailed
here.)

### Ordering

The ordering part of the system has a pretty straightforward database schema,
which centers around customers, products, SKUs, and orders.

```sql
CREATE TABLE customer (
    id       integer PRIMARY KEY,
    stripeID text    UNIQUE,
    memberID integer,
    name     text    NOT NULL,
    email    text,
    address  text,
    city     text,
    state    text,
    zip      text
);
```

Customers are basically a collection of contact information and identifiers.
`stripeID` is the ID of the customer in Stripe's database; we continue to store
credit card and similar information in Stripe, of course.  `memberID` is set
only for customers who are Schola singers; it gives their ID number in the
members site database.

```sql
CREATE TABLE product (
    id          integer PRIMARY KEY,
    stripeID    text    NOT NULL UNIQUE,
    name        text    NOT NULL,
    ticketCount integer NOT NULL,
    ticketClass text
);
CREATE TABLE sku (
    id          integer  PRIMARY KEY,
    stripeID    text     NOT NULL UNIQUE,
    product     integer  NOT NULL REFERENCES product,
    coupon      text,
    salesStart  datetime,
    salesEnd    datetime,
    membersOnly boolean  NOT NULL DEFAULT false,
    price       integer  CHECK (price IS NULL OR price >= 0),
    UNIQUE (product, coupon, membersOnly)
);
```

Products and SKUs are closely related; each product is a collection of one or
more SKUs.  Products represent what the customer gets (e.g., adult admission to
a particular concert) while SKUs represent how they got it (full price, coupon
code, early bird discount, etc.).

In the `product` table, `ticketCount` is the number of (virtual) event tickets
the customer gets as a result of buying one unit of this product.  For a
non-ticket product, this would be zero.  For a ticket to a specific event, this
would be one.  For a Flex Pass, this would be the number of events covered by
the pass (e.g. 4 for a season subscription, 6 for summer sings, etc.).
`ticketClass`, if specified, is the restriction on who is allowed to use the
ticket (e.g. "Senior" or "Student").  It is not specified for unrestricted-use
tickets or for non-ticket products.

In the `sku` table, the `price` is the price of the item (as an integer in
cents, to avoid rounding problems).  If it is `NULL`, that means the SKU's price
is variable (e.g. a donation or an auction bid).  The `coupon` is the coupon
code that the customer has to enter in order to get that price.  Every product
should have a sku with a `NULL` coupon to set the price when no valid coupon
code is entered.  The `salesStart` and `salesEnd` columns, when not `NULL`,
specify the period of time during when regular customers can place orders.
(Office staff can place orders for anything at any time.)  The `membersOnly`
flag indicates a product that can only be ordered through the members web site
(e.g. concert recordings).

It should be noted that Schola has experimented with a wide variety of sales
incentives over the years: quantity discounts, percentage discounts, early bird
discounts, etc.  This schema only accounts for the sales models currently in
use.  Additional ones can be added on demand.

Note that `product.stripeID` cannot be changed once an order is placed for that
product, and `sku.stripeID`, `sku.product`, `sku.coupon`, and `sku.membersOnly`
cannot be changed once an order is placed for that SKU.  Ideally we would treat
`sku.price` the same way, but that's actually set in Stripe's dashboard, over
which we have no control.

```sql
CREATE TABLE order (
    id       integer PRIMARY KEY,
    stripeID text     UNIQUE,
    customer integer  NOT NULL REFERENCES customer,
    source   text     NOT NULL CHECK (source IN ('door', 'web', 'other')),
    tstamp   datetime NOT NULL,
    payment  text     NOT NULL,
    note     text
);
CREATE TABLE order_line (
    id       integer PRIMARY KEY,
    orderID  integer NOT NULL REFERENCES orders,
    sku      integer NOT NULL REFERENCES sku,
    qty      integer NOT NULL CHECK (qty > 0),
    amount   integer NOT NULL CHECK (amount >= 0)
);
```

In the `order` table, the `stripeID` is Stripe's ID for the order, which is set
only if the order was paid through Stripe.The `payment` is a text description of
how the order was paid.  For Stripe purchases, it will be automatically
populated with the card type and last four digits.  The `note` is for office
use.

Orders comprise one or more order lines, stored in the `order_line` table under
the same `orderID` and with unique `lineID` values.  Each line represents a
purchase of a `qty` of a `sku`.  The `amount` is the total number of cents
charged for the line.  Normally this will be `qty * sku.price`, but it may be
different for SKUs whose price varies.

This schema makes no provision for canceled or refunded orders.  The extra
complexity needed to handle them properly is enormous, and they happen too
rarely to make that worthwhile.

### Ticket Usage Tracking

The ticket usage tracking part of the schema is more novel.  It involves events,
tickets, and admissions.

```sql
CREATE TABLE event (
    id        integer  PRIMARY KEY,
    membersID integer  UNIQUE,
    name      text     NOT NULL,
    start     datetime NOT NULL,
    capacity  integer
);
CREATE TABLE product_event (
    product integer NOT NULL REFERENCES product,
    event   integer NOT NULL REFERENCES event,
    PRIMARY KEY (product, event)
);
```

The `event` table lists events for which tickets are sold and ticket usage is
tracked.  Note that a concert with two seatings gets two rows in this table.
The `membersID` is the ID of the event in the members site database, which may
be useful for cross references.  The `capacity`, if not `NULL`, is the capacity
of the house; ticket sales will be stopped when the maximum is reached.  (But
see note below regarding Flex Passes.)

The `product_event` join table indicates which products correspond to tickets
for which events.  A product may have zero rows in this table (it's something
other than an event ticket); one row in this table (it's a ticket to a specific
event); or multiple rows in this table (it's a Flex Pass that can be used for
multiple events).

```sql
CREATE TABLE ticket (
    id         integer PRIMARY KEY,
    token      text    NOT NULL,
    order_line integer NOT NULL REFERENCES order_line,
    event      integer REFERENCES event,
    used       datetime
);
```

The `ticket` table has an entry for every purchased ticket, i.e., for every
entry of a single person to a single event.  Each ticket has a unique `id`.
Each ticket also has a random alphanumeric `token`, which is what gets encoded
into the bar code for the ticket.  For each `order_line` row, there will be
`order_line.qty * product.ticketCount` rows in the `ticket` table; those rows
will have unique `id` values but will all have the same random `token`.  For
example, if someone purchases two season subscriptions to a season with four
concert sets, there will be eight rows added to the `ticket` table, each with a
unique `id` and each with the same random `token`.  (Random alphanumeric tokens
are used in bar codes instead of integers to prevent them from being spoofed.)

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
analytics.

### Sales Configuration APIs

```x
GET    /api/product               Get list of all products
GET    /api/product/$id           Get details of a product
POST   /api/product               Create a product
PUT    /api/product/$id           Modify a product
DELETE /api/product/$id           Delete a product
GET    /api/product/$id/sku       Get list of all SKUs for a product
GET    /api/product/$id/sku/$id   Get details of one SKU for a product
POST   /api/product/$id/sku       Create a SKU for a product
PUT    /api/product/$id/sku/$id   Modify a SKU for a product
DELETE /api/product/$id/sku/$id   Delete a SKU for a product
GET    /api/event                 Get list of all events
GET    /api/event/$id             Get details of an event
POST   /api/event                 Create an event
PUT    /api/event/$id             Modify an event
DELETE /api/event/$id             Delete an event
```

These APIs are used to maintain the sales configuration data, i.e., everything
in the `product`, `sku`, `event`, and `product_event` tables described above.
(`product_event` information is embedded in the `/api/product` APIs.)  The
`POST`, `PUT`, and `DELETE` calls require authorization, using an `Auth:` header
on the HTTPS request.  The `GET` calls return limited information if they are
called without authorization.  In particular:

* Product GETs omit the `stripeID`.
* SKU GETs return information only for the `NULL` coupon and the coupon named in
  the `coupon=` parameter if any.  They omit the `stripeID`.  They do not return
  information for SKUs whose `salesStart` is in the future or whose `salesEnd`
  is in the past.
* Event GETs omit the `membersID` and `maximum`.  The event list GET returns
  information only for future events.

Consistency checks are performed on all changes.  In particular:

* Nothing can be deleted if there are references to it.
* `product.stripeID` cannot be changed once an order is placed for that product.
* `sku.stripeID`, `sku.product`, `sku.coupon`, and `sku.membersOnly` cannot be
  changed once an order is placed for that SKU.
* All `stripeIDs` are verified against Stripe's database.
* `sku.price` is populated from Stripe's database, not from the APIs.
* `event.capacity` cannot be reduced below the number of tickets already sold.

### Order Processing APIs

```x
GET    /api/order       Get list of orders
GET    /api/order/$id   Get details of one order
POST   /api/order       Create an order
PUT    /api/order/$id   Modify an existing order
DELETE /api/order/$id   Delete an existing order
```

These APIs are used to place and maintain orders (and order lines).  Except for
`POST`, they all require authorization, using an `Auth:` header on the HTTPS
request.

The `POST` call creates a new order, and if credit card information is included
in the request, it charges the credit card.  The new order is automatically
removed if the credit card charge fails.  `POST` requires authorization if:

* A payment method other than credit card is provided.
* A source other than `web` is provided.
* A SKU is provided whose `salesStart` is in the future, whose `salesEnd` is in
  the past, or whose `membersOnly` flag is set.
* A ticket is purchased that would exceed the maximum for the event.

Orders normally require a customer name and email address.  The exception is
orders that have a source `door` and contain only individual tickets for the
event at which they are purchased.  When a Flex Pass is purchased at the door,
one ticket from it is automatically marked used for the event at which it was
purchased.  When an individual ticket for an event is purchased at the door of
that event, it is automatically marked used.

> Procedural question:  the person who buys a ticket at the door (individual or
Flex pass) isn't going to have anything for the ticket scanner to scan.  How
do they get past the ticket scanner?  I'm guessing the door sales person gives
them some sort of chit (e.g. raffle ticket stub) to use in place of a bar code,
which is why I'm asserting that the ticket they purchase is automatically marked
used.  That's not quite "paperless", though.  Other ideas?

The `PUT` call allows editing an existing order.  Only the `note` and `payment`
fields can be changed, and the `payment` field can be changed only for a
non-Stripe order.

The `DELETE` call allows deleting an existing order.  If the order was charged
to a credit card, this will reverse the charge.  This is not a general-purpose
refund mechanism; it is intended only for use in an "undo" feature, to reverse
an order that was incorrectly placed moments before.  To avoid abuse, this call
will only work on orders placed within the previous five minutes.

### Ticket Collection APIs

```x
POST   /api/ticket/$token   Mark a ticket as used
DELETE /api/ticket/$token   Mark a ticket as unused
```

These calls are used by the ticket scanner.  Both of them require authorization,
using an `Auth:` header on the HTTPS request.  The `POST` request searches for
an unused ticket, valid for the current event, with the specified token.

The `DELETE` request is intended for use in an "undo" feature, such as if the
ticket scanner accidentally scans someone's pass the wrong number of times.  It
will look for a ticket with the specified token that was marked used within the
last 5 minutes, and it will make it unused again.

Both of these calls, whether successful or not, return a list of all of the
tickets with the same token, so that the ticket taker can answer questions about
how many are left unused, when the other ones were used, etc.

### Ticket Usage APIs

```x
GET /api/usage.csv
```

This call, which requires authorization, will return a CSV file listing all
tickets along with their ordering and usage information.  This can be imported
into Excel and used for investigations, charts, etc.

Once we have a better idea of what analytics we want on a recurring basis, we
can add additional APIs that return targeted data for them.  However, trying to
predict that is a fool's game.
