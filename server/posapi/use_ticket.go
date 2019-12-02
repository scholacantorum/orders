package posapi

import (
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/rothskeller/json"

	"scholacantorum.org/orders/api"
	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
)

// UseTicket is the API used by the scanner and at-the-door sales apps to get
// ticket information and mark tickets as being used.  It has two modes,
// distinguished by request method.
//
// In the GET method, it figures out the number of tickets used, available, and
// naturally consumed in each ticket class, and returns that information.
//
// In the POST method, the scan=, class=, and used= parameters dictate the new
// usage count of one or more ticket classes.  The counts may go up or down from
// their current values, but they can't go lower than the usage count at the
// time of the previous GET call.
func UseTicket(tx db.Tx, w http.ResponseWriter, r *http.Request, eventID model.EventID, token string) {
	var (
		session *model.Session
		order   *model.Order
		event   *model.Event
		now     = time.Now()
	)
	// Must have PrivScanTickets to use this API.
	if session = auth.GetSession(tx, w, r, model.PrivScanTickets); session == nil {
		return
	}
	// Make sure the requisite event exists.
	if event = tx.FetchEvent(eventID); event == nil {
		api.NotFoundError(tx, w)
		return
	}
	// Get the requested order.  It could be either an order number or an
	// order token.  If we're in POST mode, it could also be the word
	// "free", meaning that we should create an anonymous order containing
	// only free ticket class usage.
	r.ParseForm()
	if token == "free" && r.Method == http.MethodPost {
		order = &model.Order{
			Source:  model.OrderInPerson,
			Created: now,
			Valid:   true,
			Name:    "Free Entry",
		}
		r.Form["scan"] = []string{api.NewToken()}
	} else if oid, err := strconv.Atoi(token); err == nil {
		order = tx.FetchOrder(model.OrderID(oid))
	} else {
		order = tx.FetchOrderByToken(token)
	}
	if order == nil {
		api.NotFoundError(tx, w)
		return
	}
	// The rest of the processing differs between natural and explicit
	// modes.
	switch r.Method {
	case http.MethodGet:
		useTicketGet(tx, session, w, event, order, now)
	case http.MethodPost:
		useTicketPost(tx, session, w, r, event, order, now)
	}
}

// useTicketGet handles UseTicket GET requests, i.e., the ones made when a
// ticket is first scanned.  It supplies the UI with information about classes
// and ticket counts.
func useTicketGet(
	tx db.Tx, session *model.Session, w http.ResponseWriter, event *model.Event, order *model.Order, now time.Time,
) {
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
			if ol == order.Lines[0] {
				cdata.used++
			}
			cdata.max += ol.Quantity * ol.Product.TicketCount
			for _, t := range ol.Tickets {
				if !t.Used.IsZero() {
					cdata.min++
					cdata.used++
				}
			}
		}
		if cdata.used > cdata.max {
			if fp := free[cname]; fp != nil {
				if ol := addFreeTickets(order, event, fp, cdata.used-cdata.max); ol != nil {
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
	}
	// Clean up and emit results.
	api.Commit(tx)
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
				}
			})
		})
	})
	jw.Close()
}

// useTicketPost handles UseTicket POST requests, i.e., those invoked by the UI
// when the user changes the number of tickets used for a particular class of a
// particular order.
func useTicketPost(
	tx db.Tx, session *model.Session, w http.ResponseWriter, r *http.Request,
	event *model.Event, order *model.Order, now time.Time,
) {
	var (
		linemap map[string][]*model.OrderLine
		free    map[string]*model.Product
		jw      json.Writer
	)
	// Get the order lines for the requested ticket class.
	if linemap = useTicketClassMap(order, event); linemap == nil && len(order.Lines) != 0 {
		api.BadRequestError(tx, w, "not a ticket order")
		return
	}
	if len(r.Form["class"]) != len(r.Form["used"]) {
		api.BadRequestError(tx, w, "different numbers of class and used parameters")
		return
	}
	for cidx, cname := range r.Form["class"] {
		var (
			wanted int
			min    int
			max    int
			used   int
			lines  []*model.OrderLine
			fp     *model.Product
			err    error
		)
		lines = linemap[cname]
		if wanted, err = strconv.Atoi(r.Form["used"][cidx]); err != nil || wanted < 0 {
			api.BadRequestError(tx, w, "invalid used count")
			return
		}
		// Check the desired count against the min and max usage for this class.
		for _, ol := range lines {
			max += ol.Quantity * ol.Product.TicketCount
			min += ol.TicketsUsed()
		}
		used = min
		if wanted < min {
			api.BadRequestError(tx, w, "reducing used count below minimum")
			return
		}
		if wanted > max {
			if free == nil {
				free = getFreeClasses(tx, event)
			}
			if fp = free[cname]; fp == nil {
				api.SendError(tx, w, "Ticket already used")
				return
			}
			if ol := addFreeTickets(order, event, fp, wanted-max); ol != nil {
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
		log.Printf("%s USE TICKETS order:%d event:%s class:%q used:%d want:%d allow:%d-%d",
			session.Username, order.ID, event.ID, cname, used, wanted, min, max)
	}
	// Clean up and return success.
	tx.SaveOrder(order)
	api.Commit(tx)
	jw = json.NewWriter(w)
	jw.Object(func() {
		jw.Prop("id", int(order.ID))
	})
	jw.Close()
}

// useTicketError sends an error for a UseTicket request with an invalid order.
func useTicketError(tx db.Tx, w http.ResponseWriter, order *model.Order, message string) {
	api.Commit(tx)
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
			if sku.Source == model.OrderInPerson && sku.Coupon == "" && sku.InSalesRange(time.Now()) == 0 &&
				sku.Price == 0 {
				fc[p.TicketClass] = p
			}
		}
	}
	return fc
}

// addFreeTickets adds free tickets to the order, of the ticket class given in
// the specified product.  If there is an existing line with tickets of that
// class, it increases the quantity on that line and returns nil.  Otherwise, it
// adds a new line to the order and returns it.  The new tickets are not marked
// used.
func addFreeTickets(order *model.Order, event *model.Event, product *model.Product, count int) (ret *model.OrderLine) {
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
	line = &model.OrderLine{Product: product, Quantity: 0, Price: 0}
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
