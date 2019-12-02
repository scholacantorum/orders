-- Database schema for orders.scholacantorum.org.

-- The orderT table tracks all Schola Cantorum orders.  (The "T" suffix is used
-- to work around the fact that "order" is a reserved word in SQL.)  There is
-- one row in this table for each order that has been placed (successfully or
-- unsuccessfully) or is in the process of being placed.
CREATE TABLE orderT (

    -- Unique identifier of the order.  This is public, referred to as the
    -- "Schola Order Number" in customer communications.
    id integer PRIMARY KEY, -- autoincrement

    -- Another unique identifier of the order, but this one's opaque.  This is
    -- the one used in the generated QR codes.
    token string UNIQUE,

    -- Flag indicating that the order is valid, i.e., it was finalized.  An
    -- order that doesn't have this flag set either is still in process or was
    -- canceled before it was finalized.
    valid boolean NOT NULL DEFAULT 0,

    -- Source of the order, one of "public", "members", "gala", "office", or
    -- "inperson".
    source text NOT NULL,

    -- Information about the customer who placed the order.  Any of these fields
    -- may be empty.  However, if any of address, city, state, and zip are
    -- non-empty, they must all be non-empty.
    name     text    NOT NULL DEFAULT '',
    email    text    NOT NULL DEFAULT '',
    address  text    NOT NULL DEFAULT '',
    city     text    NOT NULL DEFAULT '',
    state    text    NOT NULL DEFAULT '', -- two-letter code, upper case
    zip      text    NOT NULL DEFAULT '', -- 5-digit or 9-digit code
    phone    text    NOT NULL DEFAULT '', -- ###-###-####(x#*)?
    customer text    NOT NULL DEFAULT '', -- Stripe customer ID if any
    member   integer NOT NULL DEFAULT 0,  -- Schola member ID if any

    -- Creation time of the order.
    created text NOT NULL,

    -- Customer notes about the order.
    cnote text NOT NULL DEFAULT '',

    -- Office notes about the order.
    onote text NOT NULL DEFAULT '',

    -- Flag indicating that the order and customer details have been posted
    -- into the office Access database.
    in_access boolean NOT NULL DEFAULT 0,

    -- Coupon code supplied by the customer (empty if none).
    coupon text NOT NULL DEFAULT ''
);
CREATE INDEX order_name_email_index ON orderT (name, email);
CREATE INDEX order_email_index      ON orderT (email);

-- The order_line table tracks lines of Schola Cantorum orders.  Every order has
-- at least one line.
CREATE TABLE order_line (

    -- Unique identifier of the line
    id integer PRIMARY KEY,

    -- Identifier of the order containing the line.
    orderid integer NOT NULL REFERENCES orderT ON DELETE CASCADE,

    -- The product ordered on this line.
    product text NOT NULL REFERENCES product,

    -- The quantity of the product ordered on this line.  For non-quantifiable
    -- products (e.g. donation), this is 1.
    quantity integer NOT NULL,

    -- The price per unit for this line, in cents.  The total amount for this
    -- line will always be price * quantity.
    price integer NOT NULL,

    -- Scan instance key.  A new key is assigned each time the order token is
    -- scanned (or the order ID is entered into the scanner tool).  Incremental
    -- changes are allowed only by providing the key.
    scan text NOT NULL DEFAULT '',

    -- The minimum number of tickets used for this line.  This is set to the
    -- current usage count each time the order token is scanned, so incremental
    -- changes cannot go below this.  Meaningless for non-ticket products.
    min_used integer NOT NULL DEFAULT 0,

    -- The number of tickets to be used automatically each time the order token
    -- is scanned.  Meaningless for non-ticket products.
    auto_use integer NOT NULL
);
CREATE INDEX order_line_order_index   ON order_line (orderid);
CREATE INDEX order_line_product_index ON order_line (product);

-- The product table describes products that can be ordered.  See also the sku
-- table, which describes the various pricing schemes by which these products
-- can be ordered.  A product represents what the customer gets in return for
-- their order; a SKU represents the pricing scheme by which they ordered it.
CREATE TABLE product (

    -- Unique identifier.  A text identifier is used so that it can be
    -- hard-coded in purchase forms when appropriate.  Care should be taken to
    -- ensure that IDs sort in a reasonable order within their product type.
    id text PRIMARY KEY,

    -- Series with which the product is associated.  This is usually a "20YY-YY"
    -- string for a regular season or a "20YY Summer" string for a summer sing
    -- series, or "" for products not associated with a series.
    series text NOT NULL DEFAULT '',

    -- Full name of the product, as it should appear where there is no context.
    name text NOT NULL,

    -- Short name of the product, as it should appear on an order form that
    -- gives context.
    shortname text NOT NULL,

    -- Product type.  This selects various type-specific code including ticket
    -- tracking, receipt generation, etc.  See model/product.go for values.
    type text NOT NULL,

    -- Text associated with the product on the receipt.  Only used by some
    -- product types.
    receipt text NOT NULL DEFAULT '',

    -- Number of tickets that should be issued for each unit of this product.
    -- For non-ticket products, this will be zero.  For individual event
    -- tickets, this will be 1.  For Flex Passes, this will be the number of
    -- event admissions included in the pass.
    ticket_count integer NOT NULL DEFAULT 0,

    -- Ticket class (i.e., the sort of person who is allowed to use this
    -- ticket: "Senior", "Student", etc.).  Empty string for non-ticket products
    -- or unrestricted-use tickets.
    ticket_class text NOT NULL DEFAULT ''
);

-- The sku table lists all of the SKUs that have been, are, or will be for sale.
-- Each product in the product table has one or more SKUs in the sku table,
-- representing different price points or purchase methods for the product.  A
-- product represents what the customer gets in return for their order; a SKU
-- represents the method and pricing scheme by which they ordered it.
--
-- Although not expressed in SQL, the code enforces a uniqueness constraint on
-- this table:  for any given combination of product, source, coupon, and flags,
-- there cannot be overlapping sales_start..sales_end ranges.  When choosing
-- which SKU to use for a given order, the following algorithm is used:
--   - If the order source is not "office", SKUs with a different source are not
--     considered.
--   - If the caller supplied a coupon code, and there are any SKUs with that
--     coupon code, no SKUs without that coupon code are considered.
--   - If the order source is not "office", only SKUs whose
--     sales_start..sales_end range contains the current moment are considered.
--   - If the order source is "office", preference is given to the SKU whose
--     sales_start..sales_end range contains the current moment, then to the one
--     whose range is before the current moment and closest to it, and finally
--     to the one whose range is after the current moment and closest to it.
CREATE TABLE sku (

    -- Identifier of the product that can be purchased with this SKU.
    product text NOT NULL REFERENCES product ON DELETE CASCADE,

    -- Source to which this SKU applies.  Orders can be placed with this SKU
    -- only from this source (or from the office, which can use any SKU).
    source text NOT NULL,

    -- Coupon code for this SKU.  In order to place an order with this SKU, the
    -- customer must specify this coupon code.  Each product should have a SKU
    -- with an empty coupon code, which is used when the customer doesn't
    -- specify a recognized coupon code.  Non-empty values are relevant only for
    -- source "public".
    coupon text NOT NULL DEFAULT '' COLLATE NOCASE,

    -- Start and end times for the time frame during which this SKU can be
    -- purchased by regular customers.  These restrictions do not apply to
    -- orders placed by privileged users.  Empty values indicate no limit.
    sales_start text NOT NULL DEFAULT '',
    sales_end   text NOT NULL DEFAULT '',

    -- Price to purchase the product using this SKU, in cents.  Depending on the
    -- product type, a zero value may mean that the product is free when
    -- purchased with this SKU, or it may mean that the product's price is
    -- variable and will be specified at the time of order (e.g. a donation).
    price integer NOT NULL
);
CREATE INDEX sku_product_index ON sku (product);

-- The event table lists all events to which we sell tickets.
CREATE TABLE event (

    -- Unique identifier of the event.  This is usually the YYYY-MM-DD form of
    -- the event date.  When there are multiple events on the same date, it may
    -- have an additional suffix for uniqueness.
    id text PRIMARY KEY,

    -- Identifier of the event in the event table of the
    -- scholacantorummembers.org database.
    members_id integer UNIQUE,

    -- Name of the event (as it should be shown to a customer on a receipt).  Do
    -- not include the date in the name.
    name text NOT NULL,

    -- Name of the series to which the event belongs.  This is generally
    -- "20XX-YY" or "20YY Summer".
    series text NOT NULL,

    -- Start time of the event.
    start text NOT NULL,

    -- Seating capacity of the event.  Zero means unlimited.
    capacity integer NOT NULL DEFAULT 0
);

-- The product_event table specifies which products grant admission to which
-- events.  A non-ticket product will have no entries in this table.  An
-- individual event ticket product will have a single entry in this table.  A
-- Flex Pass product will have multiple entries in this table, one for each
-- event where the Flex Pass can be used.
CREATE TABLE product_event (

    -- Identifier of the product.
    product text NOT NULL REFERENCES product ON DELETE CASCADE,

    -- Identifier of the event.
    event text NOT NULL REFERENCES event ON DELETE CASCADE,

    -- Priority of this product-event match.  When using tickets to an event,
    -- the tickets with the lower priority number are used first.  Priority zero
    -- is special in that the tickets generated from this product are considered
    -- "pre-allocated" to the event with priority zero, if they have one, for
    -- the purposes of capping sales to enforce event capacity.
    --
    -- Typical usage:
    --    0 for tickets targeted at the event
    --   10 for tickets targeted at the other seating of the same concert
    --   20 for flex passes including the event
    --   30 for tickets targeted at a different concert in the season
    priority integer NOT NULL,

    PRIMARY KEY (product, event)
);
CREATE INDEX product_event_event_index ON product_event (event);

-- The ticket table lists all tickets that have been sold.  Each ticket is a
-- single admission to a single event.
CREATE TABLE ticket (

    -- Unique identifier of the ticket.
    id integer PRIMARY KEY,

    -- Identifier of the order line on which this ticket was sold.
    order_line integer NOT NULL REFERENCES order_line ON DELETE CASCADE,

    -- Identifier of the event at which this ticket was used, or must be used.
    -- If the ticket is valid at multiple future events, this is NULL.
    event integer REFERENCES event,

    -- Timestamp when this ticket was used, or empty if it has not yet been
    -- used.
    used text NOT NULL DEFAULT ''
);
CREATE INDEX ticket_event_index      ON ticket (event);
CREATE INDEX ticket_order_line_index ON ticket (order_line);

-- The payment table tracks payments for Schola Cantorum orders.  Note that
-- this includes refunds, which are treated as negative payments.
CREATE TABLE payment (

    -- Unique identifier of the payment.
    id integer PRIMARY KEY,

    -- Identifier of the order to which this payment is associated.
    orderid integer NOT NULL REFERENCES orderT,

    -- Type of the payment method, one of "card", "card-present", "cash",
    -- "check", or "other".
    type text NOT NULL,

    -- Subtype of the payment method: basically, how the card number was
    -- provided.  Free form text.  Only used for types "card" or "card-present".
    subtype text NOT NULL,

    -- Text description of the payment method.  For "card" and "card-present"
    -- payments, this is the card type and last 4 digits (e.g. "Visa 1234").
    -- For other types, this is manual entry.  It is generally empty for "cash",
    -- the check number for "check", and free-form text for "other".
    method text NOT NULL,

    -- Stripe charge ID or refund ID, if the payment was processed through
    -- Stripe.  (Otherwise empty.)
    stripe text NOT NULL DEFAULT '',

    -- Timestamp of the payment.
    created text NOT NULL,

    -- Amount of the payment, in cents.  Negative amounts indicate refunds.
    amount integer NOT NULL
);
CREATE INDEX payment_order_index ON payment (orderid);

-- The session table lists all active user sessions.  User authentication and
-- authorization are delegated to scholacantorummembers.org.
CREATE TABLE session (

    -- Unique identifier of a login session (a random alphanumeric string).
    -- This identifier is included in all HTTPS requests to the server.
    token text PRIMARY KEY,

    -- The username of the session user.  It maps to the user.username column in
    -- the scholacantorummembers.org database.
    username text NOT NULL,

    -- The expiration date and time of the session.  After this time, the user
    -- must log in again.
    expires text NOT NULL,

    -- The user ID of the session user, in the user table of the
    -- scholacantorummembers.org table.
    member integer NOT NULL,

    -- The privileges granted to this user and therefore to this session (a
    -- bitmask).  These are derived from the user.roles column in the
    -- scholacantorummembers.org database.
    privileges integer NOT NULL
);

-- The card_email table maps payment card fingerprints (from Stripe) to email
-- addresses.  Any time a card is used with a known email address, the address
-- for that card is added or updated here.  When an order is placed with no
-- address (e.g. at the door), this table is consulted to "guess" the address
-- for the receipt.
CREATE TABLE card_email (

    -- Card fingerprint is lookup key.
    card text PRIMARY KEY,

    -- Name associated with the card.  May be empty if we've never had a name
    -- for this card (always used at the door).
    name text NOT NULL DEFAULT '',

    -- Email address associated with the card.
    email text NOT NULL DEFAULT ''
);
