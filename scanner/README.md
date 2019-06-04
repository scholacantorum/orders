# Ticket Scanner Webapp

This directory contains the source code for the ticket scanner application,
which is served from https://orders.scholacantorum.org/scanner (or the same at
orders-test, for sandbox purposes).

## UI Design

Whenever this application is launched or restarted, the first thing it does is
issue an API call to the members web site to confirm that the user is logged in
and get the user's privileges.  If they are not logged in, it redirects them to
the members web site to log in, and the members web site redirects them back.
If they are logged in but lack privileges, it says so and stops (see fatal
errors below).

Once permissions have been confirmed, it asks the server for a list of upcoming
events, and asks its user to choose which event they're taking tickets for.  At
that point we enter the main ticket-taking UI.  There are four main areas of
this UI, vertically stacked:

```x
┌─────────────┐
│             │
│ Camera View │
│             │
├─────────────┤
│ Order Info  │
├─────────────┤
│             │
│ Quantities  │
│             │
├─────────────┤
│  Controls   │
└─────────────┘
```

The Camera View displays the view from the camera, in order to aim at a QR code.
When it sees a new QR code (different from the previous one it saw), it will
change the Order Info and Quantities sections to reflect the new code.

The Order Info section displays the order number from the most recently scanned
QR code, and the customer name if known.  Clicking on it opens a dialog that
shows ticket usage in detail.

The Quantities section displays rows of buttons for selecting the quantity of
tickets used for the event.  There's one logical row for each ticket class that
is present in the customer order.  (It may wrap into multiple physical rows.)
That row contains numbered buttons from one to the number of tickets purchased.
Button colors vary, and always appear left to right in this order:

1. First are zero or more gray buttons (Bootstrap "secondary"), which represent
   tickets used at a different event.
2. Next are zero or more black buttons ("dark"), which represent tickets that
   have been previously used for this event (i.e., in a previous scan of the QR
   code).
3. Finally are zero or more blue buttons ("primary"), which represent tickets
   that have not been previously used for this event.

The black and blue buttons can be either solid or outlined.  (Gray buttons are
always solid.)  Solid means that, when the current scan session is done (i.e.,
when the next QR code is read), the ticket will be marked used.  Outline means
that, when the current scan session is done, the ticket will not be marked used.
Therefore, solid blue is a previously unused ticket that is being marked used;
outline black is a previously used ticket that is being unmarked.  All solid
buttons come to the logical left of all outline buttons.

The initial presentation of each row has every black button solid, and as many
blue buttons solid as there are tickets on the order.  For example, if the order
purchased two tickets, two blue buttons would be solid and the rest outlined.
If the order contained one flex pass, one blue button would be solid.  If there
are not enough blue buttons to meet this rule, the background color of the row
becomes red.

Clicking on a button causes every button to its right to become outlined, and
causes itself and every button to its left to become solid, without changing
color.  There are two exceptions:

* Only the rightmost of the gray buttons does this.  Other gray buttons are
  ignored.  Thus, it is not possible to change the state of tickets used at a
  different event.
* If there are no gray buttons, button 1 is solid, and every button after it is
  outlined, clicking on button 1 will turn it to an outline.  This gives a way
  to select zero without needing a dedicated (and ugly) button for zero.

If the event has any free ticket classes, there will be rows for them, even if
the order did not contain them.  These rows always have outlined blue buttons
on them, even beyond the quantity chosen on the order.  If the ticket taker
changes these to solid, the order quantity will be increased as necessary.

The rows for ticket classes other than general admission have a yellow
background to remind the ticket taker to verify that the people meet the class
restrictions.  The general admission row has a white background.

When an invalid QR code is read (one that is not a Schola order, or is an order
that does not contain tickets to the current event), the Order Info and
Quantities section are both replaced with a large error message on a red
background.

The Controls area contains two controls.  On the left is an entry field for an
order number.  Entering an order number here behaves exactly the same as if the
QR code for that order was scanned.  On the right is a "Free Entry" button.
Clicking that button clears the previously scanned order and displays a
pseudo-order containing only the free ticket classes to the event.  This allows
entry of free guests who did not "purchase" tickets in advance.  If the ticket
taker marks any of these pseudo-tickets used, an anonymous order will be created
to reflect them.

If a network error occurs, the entire page will switch to a red background with
an error message an OK button.  If some other fatal error occurs, the button
will say Reload instead.  If the camera is not accessible, the Camera View will
contain an error message, but the rest of the UI will be accessible.
