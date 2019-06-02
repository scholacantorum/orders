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

    -- Bitmask of order type and status flags.  See model/order.go for values.
    flags integer NOT NULL,

    -- Customer notes about the order.
    cnote text NOT NULL DEFAULT '',

    -- Office notes about the order.
    onote text NOT NULL DEFAULT '',

    -- Coupon code supplied by the customer (empty if none).
    coupon text NOT NULL DEFAULT '',

    -- Date on which this order should stop recurring, in seconds after the Unix
    -- epoch.  If zero, the order does not recur.  Otherwise, the order will
    -- recur on the first of every month until this date is reached.
    repeat integer NOT NULL DEFAULT 0
);
CREATE INDEX order_name_email_index ON orderT (name, email);
CREATE INDEX order_email_index      ON orderT (email);
CREATE INDEX order_repeat_index     ON orderT (repeat) WHERE repeat != 0;

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
    price integer NOT NULL
);
CREATE INDEX order_line_order_index   ON order_line (orderid);
CREATE INDEX order_line_product_index ON order_line (product);

-- The product table describes products that can be ordered.  See also the sku
-- table, which describes the various pricing schemes by which these products
-- can be ordered.  A product represents what the customer gets in return for
-- their order; a SKU represents the pricing scheme by which they ordered it.
CREATE TABLE product (

    -- Unique identifier.  A text identifier is used so that it can be
    -- hard-coded in purchase forms when appropriate.
    id text PRIMARY KEY,

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
-- represents the pricing scheme by which they ordered it.
--
-- Although not expressed in SQL, the code enforces a uniqueness constraint on
-- this table:  for any given combination of product, coupon, and members_only
-- flag, there cannot be overlapping sales_start..sales_end ranges.  When
-- choosing which SKU to use for a given order, the following algorithm is used:
--   - If the caller is not logged in, SKUs with members_only=1 are not
--     considered.
--   - If the caller is logged in, and there are any SKUs with members_only=1,
--     SKUs with members_only=0 are not considered.
--   - If the caller supplied a coupon code, and there are any SKUs with that
--     coupon code, no SKUs without that coupon code are considered.
--   - If the caller does not have PrivHandleOrders privilege, only SKUs whose
--     sales_start..sales_end range contains the current moment are considered.
--   - If the caller does have PrivHandleOrders privilege, preference is given
--     to the SKU whose sales_start..sales_end range contains the current
--     moment, then to the one whose range is before the current moment and
--     closest to it, and finally to the one whose range is after the current
--     moment and closest to it.
CREATE TABLE sku (

    -- Identifier of the product that can be purchased with this SKU.
    product text NOT NULL REFERENCES product ON DELETE CASCADE,

    -- Coupon code for this SKU.  In order to place an order with this SKU, the
    -- customer must specify this coupon code.  Each product should have a SKU
    -- with an empty coupon code, which is used when the customer doesn't
    -- specify a recognized coupon code.
    coupon text NOT NULL DEFAULT '' COLLATE NOCASE,

    -- Start and end times for the time frame during which this SKU can be
    -- purchased by regular customers.  These restrictions do not apply to
    -- orders placed by privileged users.  Empty values indicate no limit.
    sales_start text NOT NULL DEFAULT '',
    sales_end   text NOT NULL DEFAULT '',

    -- Flag indicating that this SKU can only be used by a logged-in Schola
    -- member placing an order through scholacantorummembers.org.
    members_only boolean NOT NULL DEFAULT 0,

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

    -- Type of the payment method, one of "card", "card-present", or "other".
    type text NOT NULL,

    -- Subtype of the payment method: basically, how the card number was
    -- provided.  Free form text, empty when not applicable.
    subtype text NOT NULL,

    -- Text description of the payment method.  For "card" and "card-present"
    -- payments, this is the card type and last 4 digits(e.g. "Visa 1234").  For
    -- "other" payments, this is manual entry.
    method text NOT NULL,

    -- Stripe charge ID or refund ID, if the payment was processed through
    -- Stripe.  (Otherwise empty.)
    stripe text NOT NULL DEFAULT '',

    -- Timestamp of the payment.
    created text NOT NULL,

    -- Bitmask of flags indicating the status of the payment.  See
    -- model/payment.go for values.
    flags integer NOT NULL,

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
