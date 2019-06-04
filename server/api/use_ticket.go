package api

import (
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/rothskeller/json"
	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
)

// UseTicket is the API used by the scanner app to get ticket information and
// mark tickets as being used.  When called without scan=, class=, and used=
// query parameters, it uses tickets from every class in their natural quantity.
// When called with scan=, class=, and used= parameters, it uses or unuses
// tickets of the specified class in order to achieve the desired used count.
// In either case it returns an array of class structures describing the counts.
func UseTicket(tx db.Tx, w http.ResponseWriter, r *http.Request, eventID model.EventID, token string) {
	var (
		session *model.Session
		order   *model.Order
		event   *model.Event
		scan    string
		now     = time.Now()
	)
	// Must have PrivSell to use this API.
	if session = auth.GetSession(tx, w, r, model.PrivSell); session == nil {
		return
	}
	// Make sure the requisite event exists.
	if event = tx.FetchEvent(eventID); event == nil {
		NotFoundError(tx, w)
		return
	}
	// Find out whether we're in natural or explicit mode.
	scan = r.FormValue("scan")
	// Get the requested order.  It could be either an order number or an
	// order token.  If we're in explicit mode, it could also be the word
	// "free", meaning that we should create an anonymous order containing
	// only free ticket class usage.
	if token == "free" && scan == "free" {
		order = &model.Order{
			Source:  model.OrderInPerson,
			Created: now,
			Flags:   model.OrderValid,
			Name:    "Free Entry",
		}
		r.Form["scan"] = []string{newToken()}
	} else if oid, err := strconv.Atoi(token); err == nil {
		order = tx.FetchOrder(model.OrderID(oid))
	} else {
		order = tx.FetchOrderByToken(token)
	}
	if order == nil {
		NotFoundError(tx, w)
		return
	}
	// The rest of the processing differs between natural and explicit
	// modes.
	if scan == "" {
		useTicketNatural(tx, session, w, event, order, now)
	} else {
		useTicketExplicit(tx, session, w, r, event, order, now)
	}
}

// useTicketNatural handles natural UseTicket requests, i.e., the ones made when
// a ticket is first scanned.  These consume the "natural" number of tickets out
// of each ticket class on the order, and then supply the UI with information
// about classes and ticket counts.
func useTicketNatural(tx db.Tx, session *model.Session, w http.ResponseWriter, event *model.Event, order *model.Order, now time.Time) {
	// Information collected and returned about each ticket class.
	type class struct {
		// class name
		name string
		// lowest allowed value for the usage count; this is the number
		// of tickets marked used in previous scan sessions
		min int
		// highest allowed value for the usage count; this is the total
		// number of tickets.  Set to 1000 if usage is unlimited (free
		// class).
		max int
		// usage count after this call
		used int
		// true if there were not enough tickets left to use the natural
		// number of tickets for this class
		overflow bool
	}
	var (
		lines   map[string][]*model.OrderLine
		free    map[string]*model.Product
		jw      json.Writer
		classes []*class
		scan    = newToken()
	)
	// Get the order lines for each ticket class.
	if lines = useTicketClassMap(order, event); lines == nil {
		useTicketError(tx, w, order, "Not a ticket order")
		return
	} else if len(lines) == 0 {
		useTicketError(tx, w, order, "Wrong event")
		return
	}
	// Get the free classes and make sure the lines map contains each of
	// them.
	free = getFreeClasses(tx, event)
	for fc := range free {
		if _, ok := lines[fc]; !ok {
			lines[fc] = nil
		}
	}
	// Handle each class.
	for cname, clines := range lines {
		var cdata = class{name: cname}
		classes = append(classes, &cdata)
		for _, ol := range clines {
			ol.Scan = scan
			ol.MinUsed = 0
			cdata.used += ol.AutoUse
			cdata.max += ol.Quantity * ol.Product.TicketCount
			for _, t := range ol.Tickets {
				if !t.Used.IsZero() {
					cdata.min++
					cdata.used++
					ol.MinUsed++
				}
			}
		}
		if cdata.used > cdata.max {
			if fp := free[cname]; fp != nil {
				if ol := addFreeTickets(order, event, fp, scan, cdata.used-cdata.max); ol != nil {
					lines[cname] = append(lines[cname], ol)
				}
				cdata.max = 1000
			} else {
				cdata.used = cdata.max
				cdata.overflow = true
			}
		} else if free[cname] != nil {
			cdata.max = 1000
		}
		if cdata.used > cdata.min {
			consumeTickets(clines, event, now, cdata.used-cdata.min)
		}
	}
	// Clean up and emit results.
	tx.SaveOrder(order)
	commit(tx)
	sort.Slice(classes, func(i, j int) bool {
		return classes[i].name < classes[j].name
	})
	w.Header().Set("Content-Type", "application/json")
	jw = json.NewWriter(w)
	jw.Object(func() {
		jw.Prop("id", int(order.ID))
		if order.Name != "" {
			jw.Prop("name", order.Name)
		}
		jw.Prop("scan", scan)
		jw.Prop("classes", func() {
			jw.Array(func() {
				for _, class := range classes {
					jw.Object(func() {
						jw.Prop("name", class.name)
						jw.Prop("min", class.min)
						jw.Prop("max", class.max)
						jw.Prop("used", class.used)
						if class.overflow {
							jw.Prop("overflow", true)
						}
					})
					log.Printf("%s USE TICKETS order:%d event:%s %v",
						session.Username, order.ID, event.ID, class)
				}
			})
		})
	})
	jw.Close()
}

// useTicketExplicit handles explicit UseTicket requests, i.e., those invoked by
// the UI when the user changes the number of tickets used for a particular
// class of a particular order.
func useTicketExplicit(
	tx db.Tx, session *model.Session, w http.ResponseWriter, r *http.Request,
	event *model.Event, order *model.Order, now time.Time,
) {
	var (
		wanted int
		min    int
		max    int
		used   int
		lines  []*model.OrderLine
		free   map[string]*model.Product
		fp     *model.Product
		jw     json.Writer
		err    error
		scan   = r.FormValue("scan")
		cname  = r.FormValue("class")
	)
	if wanted, err = strconv.Atoi(r.FormValue("used")); err != nil || wanted < 0 {
		BadRequestError(tx, w, "invalid used count")
		return
	}
	// Get the order lines for the requested ticket class.
	if linemap := useTicketClassMap(order, event); linemap == nil && len(order.Lines) != 0 {
		BadRequestError(tx, w, "not a ticket order")
		return
	} else {
		lines = linemap[cname]
	}
	// Check the desired count against the min and max usage for this class.
	for _, ol := range lines {
		if ol.Scan != scan {
			BadRequestError(tx, w, "wrong scan session")
			return
		}
		max += ol.Quantity * ol.Product.TicketCount
		min += ol.MinUsed
		for _, t := range ol.Tickets {
			if !t.Used.IsZero() {
				used++
			}
		}
	}
	if wanted < min {
		BadRequestError(tx, w, "reducing used count below minimum")
		return
	}
	if wanted > max {
		free = getFreeClasses(tx, event)
		if fp = free[cname]; fp == nil {
			sendError(tx, w, "Ticket already used")
			return
		}
		if ol := addFreeTickets(order, event, fp, scan, wanted-max); ol != nil {
			lines = append(lines, ol)
		}
		max = 1000
	}
	// Adjust the usage as requested.
	if wanted > used {
		consumeTickets(lines, event, now, wanted-used)
	}
	if wanted < used {
		unconsumeTickets(lines, event, used-wanted)
	}
	// Clean up and return success.
	tx.SaveOrder(order)
	commit(tx)
	log.Printf("%s USE TICKETS order:%d event:%s class:%q used:%d want:%d allow:%d-%d",
		session.Username, order.ID, event.ID, cname, used, wanted, min, max)
	jw = json.NewWriter(w)
	jw.Object(func() {
		jw.Prop("id", int(order.ID))
		jw.Prop("scan", scan)
	})
	jw.Close()
}

// useTicketError sends an error for a UseTicket request with an invalid order.
func useTicketError(tx db.Tx, w http.ResponseWriter, order *model.Order, message string) {
	commit(tx)
	w.Header().Set("Content-Type", "application/json")
	jw := json.NewWriter(w)
	jw.Object(func() {
		jw.Prop("id", int(order.ID))
		if order.Name != "" {
			jw.Prop("name", order.Name)
		}
		jw.Prop("error", message)
	})
	jw.Close()
}

// useTicketClassMap returns a map from ticket class name to the list of order
// lines on the specified order containing tickets to the specified event with
// the named class.  Each list is in priority order.  If the map is empty, the
// order contains ticket lines but none that are applicable to the specified
// event.  If the returned map is nil, the order does not contain ticket lines.
func useTicketClassMap(order *model.Order, event *model.Event) (cm map[string][]*model.OrderLine) {
	var (
		ticket bool
		prios  = map[model.OrderLineID]int{}
	)
	// Build the map.
	cm = make(map[string][]*model.OrderLine)
	for _, ol := range order.Lines {
		for _, pe := range ol.Product.Events {
			ticket = true
			if pe.Event.ID == event.ID {
				cm[ol.Product.TicketClass] = append(cm[ol.Product.TicketClass], ol)
				prios[ol.ID] = pe.Priority
			}
		}
	}
	if !ticket {
		return nil
	}
	// Prioritize the lists.
	for _, cl := range cm {
		sort.Slice(cl, func(i, j int) bool {
			return prios[cl[i].ID] < prios[cl[j].ID]
		})
	}
	return cm
}

// getFreeClasses returns a map from class name to product for each ticket class
// that is available free to the specified event.
func getFreeClasses(tx db.Tx, event *model.Event) (fc map[string]*model.Product) {
	fc = make(map[string]*model.Product)
	for _, p := range tx.FetchProductsByEvent(event) {
		for _, sku := range p.SKUs {
			if sku.Price == 0 && sku.Coupon == "" && !sku.MembersOnly {
				fc[p.TicketClass] = p
			}
			// Note deliberately ignoring SalesStart..SalesEnd
			// range.  Free student tickets are usually set up with
			// a range that never matches so they don't appear as
			// explicitly orderable.
		}
	}
	return fc
}

// addFreeTickets adds free tickets to the order, of the ticket class given in
// the specified product.  If there is an existing line with tickets of that
// class, it increases the quantity on that line and returns nil.  Otherwise, it
// adds a new line to the order and returns it.  The new tickets are not marked
// used.
func addFreeTickets(order *model.Order, event *model.Event, product *model.Product, scan string, count int) (ret *model.OrderLine) {
	var line *model.OrderLine

	for _, ol := range order.Lines {
		if ol.Product.TicketClass != product.TicketClass {
			continue
		}
		for _, pe := range ol.Product.Events {
			if pe.Event.ID == event.ID {
				line = ol
				goto found
			}
		}
	}
	line = &model.OrderLine{Product: product, Scan: scan, AutoUse: 0, Quantity: 0, MinUsed: 0, Price: 0}
	order.Lines = append(order.Lines, line)
	ret = line
found:
	line.Quantity += count
	for i := 0; i < count; i++ {
		line.Tickets = append(line.Tickets, &model.Ticket{})
	}
	return ret
}

// consumeTickets consumes count tickets of the specified ticket class to the
// specified event, using the tickets on the supplied set of lines.
func consumeTickets(lines []*model.OrderLine, event *model.Event, now time.Time, count int) {
	for _, ol := range lines {
		for _, t := range ol.Tickets {
			if t.Used.IsZero() {
				t.Used = now
				t.Event = event
				count--
				if count == 0 {
					return
				}
			}
		}
	}
	panic("ran out of unused tickets")
}

// unconsumeTickets unconsumes count tickets of the specified ticket class to
// the specified event, using the tickets on the supplied set of lines.
func unconsumeTickets(lines []*model.OrderLine, event *model.Event, count int) {
	for i := len(lines) - 1; i >= 0; i-- {
		ol := lines[i]
		for j := len(ol.Tickets) - 1; j >= 0; j-- {
			t := ol.Tickets[j]
			if !t.Used.IsZero() {
				t.Used = time.Time{}
				count--
				if count == 0 {
					return
				}
			}
		}
	}
	panic("ran out of used tickets")
}
