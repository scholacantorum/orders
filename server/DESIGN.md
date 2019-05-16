# Redesigning the Schola Ordering System

Moving into the 2019–20 season, we want to abandon all use of paper ticketing,
and we want to track actual ticket usage.  Our current ordering system,
instituted at the beginning of the 2018–19 season, cannot be readily extended to
meet these new goals.

## Why Track Ticket Usage?

The key issue is the tracking of ticket usage.  There are several reasons to
want it:

* We want to prevent multiple people from entering using the same email order
  confirmation (for an individual ticket).  This is the most frequently cited
  reason, but actually the least persuasive since we have no reason to believe
  such ticket fraud is happening.  (Indeed, per Dan, "it would be a nice
  problem to have.")

* When we go paperless for multi-use tickets like Flex Passes for season
  subscriptions and Summer Sings, we again want to prevent ticket fraud by
  preventing the ticket from being used more times than it's good for.  Again,
  we have no evidence of ticket fraud, but with multiple-use tickets the
  potential loss is higher, and the possibility of *accidental* misuse is also
  higher, so this is a stronger reason.

* We have some events (e.g. this past gala, or Elijah if we stay at one seating)
  that have realistic chances of selling out.  We should be counting tickets
  sold and stopping sales before we overbook.  We have no provision for that
  today.

* We would like to do better analytics of ticket usage, both statistically (e.g.
  which seatings are more popular?) and individually (which concerts did this
  particular donor attend?).  Right now we have good analytics on ticket *sales*
  but poor analytics on ticket *usage*, and no correlation between them.

## Why Redesign?

Our current ordering system has no underlying database that we control.  When
online purchases are made, the purchases are recorded in Stripe and in a Google
Docs spreadsheet.  But neither of these records are suitable for ticket usage
tracking.  Stripe's database is limited to online purchases through Stripe only,
and doesn't have a flexible schema that could handle usage tracking.  The Google
Docs spreadsheet does not support database-style usage; it's effectively a
write-only archive.  To do ticket usage tracking, we will need to track all
ticket sales in a real database, and our current software can't do that.

Of course, before writing custom software, the first thing to do is look for
off-the-shelf alternatives.  The current scholacantorum.org is designed to be
easy for someone besides me to maintain.  A custom ordering system with ticket
usage tracking would necessarily be an order of magnitude more difficult.

Unfortunately, I haven't found any off-the-shelf alternative that meets our
needs.  To use any of the OTS services, like Brown Paper Tickets or Eventbrite,
we would need to change our ticketing policies considerably.  We would need to
abandon the Flex Pass concept; I haven't seen any OTS service that can support
that.  And the OTS services are all expensive, so we would probably have to pass
their ticket service charges onto our patrons.  That would drive more sales
offline, which would greatly degrade our usage data since we can't capture good
analytics for at-the-door sales.

Bottom line, I don't think we have a choice.  If we want paperless tickets and
good usage tracking, it's going to be a custom design, and we'll have to accept
the long-term maintenance cost.

## Back End Technology Choices

The ordering system will need to have a database, and server code to manage it.
Up-front design choices include the authentication and authorization model, the
database type, the web server technology, the hosting domain, and the
implementation language.

> **TL;DR** back end written in Go using SQLite database and CGI protocol,
hosted on a subdomain of scholacantorum.org, delegating authentication and
authorization to the existing Schola members site.

### Authentication and Authorization Model

Some parts of the ordering / ticketing system will need user authentication and
authorization.  For example, cash sales should only be recorded by at-the-door
sales staff, and analytics should only be available to the office staff.  A
proper implementation of authentication requires username and password storage
with encryption, password quality enforcement, profile management for changing
passwords, "forgot my password" recovery functionality, automatic lockout after
successive failed attempts, and "remember me on this computer" functionality.

Schola already has one implementation of these features, in its members web
site.  The ordering / ticketing system will either need to implement them
separately, or it will need to leverage them from the members site.  Leveraging
from the members site would reduce effort and maintenance cost.  It would also
avoid giving the office staff yet another username/password combination to
remember.

Leverage could be achieved by sharing the data and code, or it could be done by
delegation.  For example, the ordering / ticketing system could validate a user
by looking it up in the same database that the members site uses, perhaps even
sharing the same code to do so.  Or, the ordering / ticketing system could
validate a user by sending a request to the members site and asking it to do so.

Sharing the data and code would tie the two systems together unnecessarily; it
would violate basic software architecture principles.  Also, it would dictate
the answers to most of the other technology choices:  SQLite database,
continuously running web server hosted on scholacantorummembers.org, written in
Go.  There is no reason to accept those constraints.  Therefore, I believe the
best option is to leverage the members site's authentication and authorization
models via delegation.

### Database

The nature of the application calls for a traditional SQL-based relational
database, as opposed to a NoSQL alternative like a document store, object store,
or graph database.  Within the execution environment provided by our web host,
the available relational databases are MySQL and SQLite.

MySQL is a full-featured database server.  The database server is separate from
the web server; both are maintained by our web hosting company.  Access to the
database from the back end code is via a remote network connection.  This adds
quite a bit of overhead, although our traffic is low enough that it might not
matter.  It also adds potential failure modes (e.g. the web server is up but it
can't talk to the database server).  One advantage of MySQL is that it can be
accessed remotely by software running someplace other than our web server, if we
so choose.

SQLite is an embedded database: it is a library that becomes part of our back
end code, rather than running as a separate server.  It stores its data in the
file system of the web server, so the data are accessible only by code running
on the web server.  It is lightning fast and very low overhead.  It does not
implement the full SQL standard, but the parts it leaves out are not critical.
Since the data are stored in the file system, a separate mechanism is needed to
ensure the files get backed up.  Its locking mechanisms are coarse grained, and
not really suitable for handling highly concurrent workloads, but our workloads
are not demanding in that regard.

Either database system would serve our needs.  I believe SQLite is the better
choice for us, for two reasons.  First, whenever choosing between third party
software packages, I always lean towards the smallest and least complex one that
will serve the need, and SQLite is far simpler than MySQL.  Second, SQLite is
what the members and gala sites use, so there would be less of a maintenance
burden of needing to understand multiple database systems.

### Web Server Protocol

In this context, the three web server protocol choices are CGI, continuous
HTTPS, or continuous HTTPS with web sockets.

Web sockets allow continuous, asynchronous, two-way communications between
client and server.  They are finicky, and not too many people know how to work
with them, so they add significantly to the maintenance burden.  The gala
software needed that asynchronous communication, and it uses web sockets.
However, I don't antipate that the ordering / ticketing system would need that,
and maintenance issues are more critical for it, so I would argue against using
web sockets for it.

CGI systems start up a new instance of the web server software for each
incoming request.  A continuous HTTPS server is needed when the overhead of
each startup, multiplied by the volume of incoming requests, is too high to
provide adequate user response time.  Schola's members site uses a continuous
HTTPS server because it maintains an in-memory cache of most of its data, and
filling that cache on startup is an expensive operation.  Schola's public site
uses a CGI server because it runs on a web host that does not allow continuous
HTTPS servers.

One can never be certain in advance of actual measurements, but I have no reason
to think that either the startup overhead or the request volume for our ordering
/ ticketing system will be high enough to matter.  I would use a continuous
HTTPS server anyway if it is an option, just on general efficiency principles.
But a CGI server would also meet our needs.

### Hosting Domain

I see no reason to create a new top-level domain, so our choices for hosting are
scholacantorum.org, scholacantorummembers.org, or a subdomain of one of those.
Both of those are hosted by Dreamhost, but in different accounts with different
hosting plans.

scholacantorum.org and its subdomains are hosted in a free hosting plan in
Schola Cantorum's name.  (It's free because Schola is a non-profit.)  The free
plan has one limitation that's relevant in this context:  it does not allow
continuously running server code.  Only CGI servers are supported.  However, as
noted above, that's probably OK for our ordering / ticketing system.

scholacantorummembers.org and its subdomains are hosted in a paid hosting
account under my name.  This is for exactly the reason mentioned above: it needs
a continuously running server for performance, and Schola's free plan didn't
allow it.  (I don't charge Schola for this hosting account; I was paying for it
anyway, to host other non-Schola domains, so adding Schola to it didn't cost me
anything more.)  I hosted the Gala software at gala.scholacantorummembers.org
for the same reason, but stronger: it uses web sockets, so it *had* to have a
continuously running server.

I recommend hosting the ordering / ticketing system on a new subdomain of
scholacantorum.org called orders.scholacantorum.org.  It seems to belong better
on the public side conceptually, since it's mostly a public-facing system rather
than a member-facing system.  And it seems better to have our ordering data in a
hosting account in Schola's name.  The fact that it forces us to use a CGI
server is unfortunate but tolerable.

### Implementation Language

All of our existing back end code, for the public, members, and gala sites, is
in Go.  Go is a compiled, multi-threaded, type-safe language, so it is
dramatically faster and slightly safer than the more commonly used web server
languages like PHP, Python, and Ruby, all of which are interpreted,
single-threaded, and untyped.  Ruby has the additional disadvantage that it is
not readily available in our web hosting environment.

In an effort to minimize the number of languages needed to maintain Schola's
suite of web sites, I believe the back end of the ordering / ticketing system
should also be in Go.

## Front End Technology Choices

Depending on how you count them, we have six front ends to consider:

* Public web site: patrons place online orders there.
* Members web site: members place online orders there (e.g. concert recordings).
* Gala web site: gala staff charges credit cards there.
* Office web site: office staff record offline orders there, and review analytics there.
* Phone app for door sales: staff charge credit cards and record cash receipts there.
* Phone app for ticket scanning: staff record ticket usage there.

These could be merged in some cases (e.g. the office web site could be
implemented as a protected part of the members web site, or the two phone apps
could be implemented as one phone app).

Client code that runs in a browser (e.g., the four web sites above) must be
written in Javascript; that is the only supported language.  The small amounts
of existing Javascript code on the Schola public web site are written natively
without any framework.  The much larger bodies of Javascript code on the members
and gala web sites are written using the Vue framework.  The Vue framework
allows one to write UIs in a declarative fashion rather than a procedural one,
which can greatly simplify coding, at the cost of needing to learn that
framework.  React is a similar (and more commonly used) framework.

For client code on a phone, there are three choices:

* a webapp, written in JavaScript, possibly with the Vue or React framework,
  that's run in a browser on the phone
* a native app, written in the language used by the phone OS (Swift for iOS,
  Java for Android)
* a native app, written in JavaScript with the Vue Native or React Native
  framework, and cross-compiled into the language used by the phone OS

Any of these choices can read QR codes using the phone's camera, and any of them
can process credit cards using a Swipe card reader.  (Note, however, that the
necessary card reader is $60 for the native apps and $300 for the webapp.)
Steve H. has existing code for reading QR codes using the phone's camera, which
he said is written in JavaScript, but I don't know whether that's run in
JavaScript as a webapp or run natively after cross-compilation.

I am inclined to exclude the middle choice, writing in Swift or Java.  I don't
know Swift, and I dislike Java.  More to the point, we should minimize the
number of languages that maintainers of this code need to know, and since they
will already need to know JavaScript, there is little value in adding another
language to the list.  The middle choice also has the disadvantage of needing
two separate implementations if we want to support both iOS and Android.

Between the remaining choices, I lean towards a webapp running in JavaScript.
That is less tightly integrated into the phone UI, and produces a subtly less
desirable user experience, than a native app.  But since the users of the webapp
would be our own staff and ticket takers, I don't think that matters too much.
Meanwhile, the native app brings in a huge number of new technologies, and
greatly increases the complexity of the app and its maintenance burden.  Also,
distribution of a native app requires submitting it to the Apple and/or Google
App stores, which is a tedious and sometimes costly process.

Before considering this choice final, I want to do some prototyping to verify
that I really can read a QR code, with adequate performance, from a webapp on a
phone.  So this should be considered tentative.

## Block Diagram

```x
┌────────┐   ┌─────────┐   ┌───────────┐   ┌──────────┐   ┌─────────────┐
│ sc.org │   │ scm.org │   │ g.scm.org │   │ o.sc.org ├───┤ Card Reader │
└──┬──┬──┘   └──┬───┬──┘   └───┬───┬───┘   └──┬───┬───┘   └──────┬──────┘
   │  └─────────│───┴──────────│───┴──────────│───┴──────────────┤
   └────────────┴───────┬──────┴──────────────┘                  │
┌───────────────────────┴───────────────────┬──────────┐         │
│               Back End Software           │ Database │         │
└───┬───────────────┬────────────────────┬──┴──────────┘         │
┌───┴───┐   ┌───────┴──────┐          ┌──┴───────────────────────┴──────┐
│ Email │   │ Google Sheet │          │             Stripe              │
└───────┘   └──────────────┘          └─────────────────────────────────┘
```

The four client applications that take orders (scholacantorum.org,
scholacantorummembers.org, gala.scholacantorummembers.org, and
orders.scholacantorum.org) each communicate with the back end software and also
with Stripe.  orders.scholacantorum.org, which hosts the phone interface for
at-the-door sales, also communicates with the card reader, which communicates
with Stripe.  The back end software contains the database, and communicates with
Stripe; it culminates with updating the Google sheet and sending emails.

The Stripe communications deserve a bit more detail.  When one of our web sites
shows a credit card entry field, that field is not actually coming from our site
code at all; it is in an IFRAME tag, and is coming from Stripe's own servers.
Anything the user types in it is sent directly to Stripe and is never seen by
our code.  What Stripe gives our site's code is a "token" that represents the
user input without exposing it.  Our site's code passes that to our back end
code, which in turn passes it to Stripe's API.  Stripe then finds the real card
data associated with the token, and charges that card.  A similar process
happens with the phone app and the card reader: the card reader sends the card
data directly to Stripe, and gives the corresponding token to our phone app.

## Database Schema

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
