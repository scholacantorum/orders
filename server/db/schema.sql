-- Database schema for orders.scholacantorum.org.

-- The customer table lists customers who have placed orders from us (including
-- our own singers).  When a new order is placed, an existing customer entry
-- will be reused if (a) either the names are the same or one of the names is
-- empty; (b) either the emails are the same or one of the emails is empty; and
-- (c) the memberIDs are the same.  Otherwise, a new customer will be created
-- (both in this database and in Stripe).
CREATE TABLE customer (

    -- Unique identifier of the customer.
    id integer PRIMARY KEY,

    -- Stripe ID for the customer.  This will be populated only if the customer
    -- paid for an order through Stripe.
    stripeID text UNIQUE,

    -- Member ID for the customer.  This will be non-zero only if the customer
    -- is a member of Schola, in which case it will map to the user.id column in
    -- the scholacantorummembers.org database.
    memberID integer NOT NULL DEFAULT 0,

    -- Name of the customer.  This may be empty for in-person sales where no
    -- name is given.
    name text NOT NULL DEFAULT '' COLLATE NOCASE,

    -- Email address of the customer.  This may be empty.
    email text NOT NULL DEFAULT '' COLLATE NOCASE,

    -- Other contact information for the customer.  These columns may be empty.
    address text NOT NULL DEFAULT '',
    city    text NOT NULL DEFAULT '',
    state   text NOT NULL DEFAULT '',
    zip     text NOT NULL DEFAULT '',
    phone   text NOT NULL DEFAULT ''
);
CREATE INDEX customer_name_email_index ON customer (name, email, memberID);
CREATE INDEX customer_email_index      ON customer (email, memberID);

-- The event table lists all events to which we sell tickets.
CREATE TABLE event (

    -- Unique identifier of the event.
    id integer PRIMARY KEY,

    -- Identifier of the event in the event table of the
    -- scholacantorummembers.org database.
    membersID integer UNIQUE,

    -- Name of the event (as it should be shown to a customer on a receipt).  Do
    -- not include the date in the name.
    name text NOT NULL,

    -- Start time of the event (in seconds since the Unix epoch).
    start integer NOT NULL,

    -- Seating capacity of the event.  Zero means unlimited.
    capacity integer NOT NULL DEFAULT 0
);

-- The product table lists all products that have been, are, or will be for
-- sale.  (See also the "sku" table.)
CREATE TABLE product (

    -- Unique identifier of the product.
    id integer PRIMARY KEY,

    -- Stripe ID for the product.
    stripeID text NOT NULL UNIQUE,

    -- Name of the product (as it should appear on customer receipts).
    name text NOT NULL,

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

-- The product_event table specifies which products grant admission to which
-- events.  A non-ticket product will have no entries in this table.  An
-- individual event ticket product will have a single entry in this table.  A
-- Flex Pass product will have multiple entries in this table, one for each
-- event where the Flex Pass can be used.
CREATE TABLE product_event (

    -- Identifier of the product.
    product integer NOT NULL REFERENCES product ON DELETE CASCADE,

    -- Identifier of the event.
    event integer NOT NULL REFERENCES event ON DELETE CASCADE,
    PRIMARY KEY (product, event)
);
CREATE INDEX product_event_event_index ON product_event (event);

-- The sale table lists all sales.
CREATE TABLE sale (

    -- Unique identifier of the sale.  This is customer visible, as the "Schola
    -- order number".
    id integer PRIMARY KEY,

    -- Stripe ID of the corresponding Stripe order, if the sale was paid through
    -- Stripe.
    stripeID text UNIQUE,

    -- Identifier of the sale customer.
    customer integer NOT NULL REFERENCES customer,

    -- Source of the sale.  Possible values:
    --   - 'P' scholacantorum.org public site
    --   - 'M' scholacantorummembers.org site
    --   - 'G' gala.scholacantorummembers.org site
    --   - 'O' offline order processed by office staff
    --   - 'D' at-the-door or other in-person sale
    source text NOT NULL CHECK (source IN ('D','G','M','O','P')),

    -- Timestamp of the sale (in seconds since the Unix epoch).
    timestamp integer NOT NULL,

    -- Payment method for the sale.  If the sale was paid through Stripe, this
    -- will be a description of the card used (e.g. "Visa 5023").  Otherwise,
    -- this is a free-form text field with an office-supplied description.
    payment text NOT NULL CHECK (payment != ''),

    -- Office-supplied notes about the sale.
    note text NOT NULL DEFAULT ''
);
CREATE INDEX sale_customer_index ON sale (customer);

-- The sale_line table lists the specific line items in each sale.
CREATE TABLE sale_line (

    -- Unique identifier of the sale line.
    id integer PRIMARY KEY,

    -- Identifier of the sale to which this sale line belongs.
    sale integer NOT NULL REFERENCES sale ON DELETE CASCADE,

    -- Identifier of the SKU sold on this sale line.
    sku integer NOT NULL REFERENCES sku,

    -- The number of units of the underlying product sold on this sale line.
    qty integer NOT NULL CHECK (qty > 0),

    -- The total amount charged for this sale line, in cents.  This will usually
    -- be qty * sku.price, but it may be different if the SKU has a variable
    -- price (e.g. a donation) or if the SKU's price has changed since the
    -- sale was processed.
    amount  integer NOT NULL CHECK (amount >= 0)
);
CREATE INDEX sale_line_sale_index ON sale_line (sale);
CREATE INDEX sale_line_sku_index  ON sale_line (sku);

-- The session table lists all active user sessions.  User authentication and
-- authorization are delegated to scholacantorummembers.org.
CREATE TABLE session (

    -- Unique identifier of a login session (a random alphanumeric string).
    -- This identifier is included in all HTTPS requests to the server.
    token text PRIMARY KEY,

    -- The username of the session user.  It maps to the user.username column in
    -- the scholacantorummembers.org database.
    username text NOT NULL,

    -- The expiration date and time of the session (in seconds since the Unix
    -- epoch).  After this time, the user must log in again.
    expires integer NOT NULL,

    -- The privileges granted to this user and therefore to this session (a
    -- bitmask).  These are derived from the user.roles column in the
    -- scholacantorummembers.org database.
    privileges integer NOT NULL
);

-- The sku table lists all of the SKUs that have been, are, or will be for sale.
-- Each product in the product table has one or more SKUs in the sku table,
-- representing different price points or purchase methods for the product.
-- Generally speaking, the product is what the customer receives in return for
-- their money; the SKU represents the sales arrangements for the purchase.
--
-- Although not expressed in SQL, the code enforces a uniqueness constraint on
-- this table:  for any given combination of product, coupon, and members_only
-- flag, there cannot be overlapping sales_start..sales_end ranges.  Thus, when
-- an order for a product is placed, there are at most four matching SKUs:
--   - matching coupon code and members_only flag set
--   - empty coupon code    and members_only flag set
--   - matching coupon code and members_only flag clear
--   - empty coupon code    and members_only flag clear
-- The first one of those that applies to the purchase is the one used.
CREATE TABLE sku (

    -- Unique identifier of the SKU.
    id integer PRIMARY KEY,

    -- Stripe ID for the SKU.
    stripeID text NOT NULL UNIQUE,

    -- Identifier of the product that can be purchased with this SKU.
    product integer NOT NULL REFERENCES product,

    -- Coupon code for this SKU.  In order to place an order with this SKU, the
    -- customer must specify this coupon code.  Each product should have a SKU
    -- with an empty coupon code, which is used when the customer doesn't
    -- specify a recognized coupon code.
    coupon text NOT NULL DEFAULT '' COLLATE NOCASE,

    -- Start and end times for the time frame during which this SKU can be
    -- purchased by regular customers (as seconds since the Unix epoch).  These
    -- restrictions do not apply to orders placed by privileged users.  Zero
    -- values indicate no limit.
    sales_start integer NOT NULL DEFAULT 0,
    sales_end   integer NOT NULL DEFAULT 0,

    -- Flag indicating that this SKU can only be used by a logged-in Schola
    -- member placing an order through scholacantorummembers.org.
    members_only boolean NOT NULL DEFAULT 0,

    -- Price to purchase the product using this SKU, in cents.  A zero value may
    -- mean that the product is free when purchased with this SKU, or it may
    -- mean that the product's price is variable and will be specified at the
    -- time of order (e.g. a donation).
    price integer NOT NULL CHECK (price >= 0)
);
CREATE INDEX sku_product_index ON sku (product);

-- The ticket table lists all tickets that have been sold.  Each ticket is a
-- single admission to a single event.
CREATE TABLE ticket (

    -- Unique identifier of the ticket.
    id integer PRIMARY KEY,

    -- Token representing the set of interchangeable tickets purchased with this
    -- one (i.e., all tickets purchased with the same sale_line have the same
    -- token).  This is a random alphanumeric string.
    token text NOT NULL,

    -- Identifier of the sale line on which this ticket was sold.
    sale_line integer NOT NULL REFERENCES sale_line,

    -- Identifier of the event at which this ticket was used, or must be used.
    -- If the ticket is valid at multiple future events, this is NULL.
    event integer REFERENCES event,

    -- Timestamp when this ticket was used (in seconds past the Unix epoch), or
    -- zero if it has not yet been used.
    used integer NOT NULL DEFAULT 0 CHECK (used=0 OR event IS NOT NULL)
);
CREATE INDEX ticket_event_index     ON ticket (event);
CREATE INDEX ticket_sale_line_index ON ticket (sale_line);
CREATE INDEX ticket_token_index     ON ticket (token, used);
