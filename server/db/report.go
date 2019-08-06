package db

import (
	"database/sql"
	"strings"
	"time"

	"scholacantorum.org/orders/model"
)

// Maximum number of matched records.
const maxReportSize = 500

// These constants are used as a bitmask parameter to the lineMatches function,
// telling it which of the report criteria to enforce (either all of them or all
// but one of them).
const (
	critOrderSource byte = 1 << iota
	critCustomer
	critOrderCreated
	critOrderCoupon
	critProduct
	critPaymentType
	critTicketClass
	critUsedAtEvent
	critAll = critOrderSource | critCustomer | critOrderCreated | critOrderCoupon |
		critProduct | critPaymentType | critTicketClass | critUsedAtEvent
)

// reportProduct contains the information we cache about each product in the DB.
type reportProduct struct {
	ptype  model.ProductType
	series string
	name   string
	tclass string
	tcount int
	count  int
}

// reportOrder contains the information we cache about the order we're currently
// handling.  It persists through all of the lines on the same order and is then
// replaced for the next order.
type reportOrder struct {
	id          model.OrderID
	source      model.OrderSource
	name        string
	email       string
	created     time.Time
	flags       model.OrderFlags
	coupon      string
	paymentType string
	counted     bool
}

// reportLine contains the information we cache about each order line while
// processing it.
type reportLine struct {
	pid   model.ProductID
	prod  *reportProduct
	qty   int
	price int
	order *reportOrder

	// This field maps event ID to the number of tickets used at that event,
	// with an entry for "" counting unused tickets.  It's nil for non-
	// ticket lines or historical ticket lines where usage wasn't tracked.
	tusage map[model.EventID]int
}

// paymentTypeMap maps from the type and subtype fields of the payment table
// (concatenated with a comma) to the PaymentType string we used for reporting.
// Note that this is not a one-to-one mapping; some distinctions in the subtypes
// are not useful in reports.  Entries without a comma in the key are used as
// prefixes for unrecognized subtypes.
var paymentTypeMap = map[string]string{
	",":                            "Free",
	"card":                         "Card,",
	"card,":                        "Card,Typed",
	"card,API":                     "Wallet,Browser",
	"card,API basic-card":          "Wallet,Browser",
	"card,apple_pay":               "Wallet,Apple Pay",
	"card,google_pay":              "Wallet,Google Pay",
	"card,manual":                  "Card,Typed",
	"card-present":                 "Card Present,",
	"card-present,":                "Card Present,Unknown",
	"card-present,contact_emv":     "Card,Inserted",
	"card-present,contactless_emv": "Card,Tapped",
	"card-present,contactless_magstripe_mode": "Card,Tapped",
	"card-present,magnetic_stripe_fallback":   "Card,Swiped",
	"card-present,magnetic_stripe_track2":     "Card,Swiped",
	"other":                                   "Other ",
	"other,":                                  "Other",
	"other,cash":                              "Cash",
	"other,check":                             "Check",
	"saved":                                   "Card, Reuse Saved ",
	"saved,manual":                            "Card,Reuse Saved",
}

// RunReport executes the report defined by the supplied report definition.
func (tx Tx) RunReport(def *model.ReportDefinition) *model.ReportResults {
	var (
		rows            *sql.Rows
		ticketUsageStmt *sql.Stmt
		orderStmt       *sql.Stmt
		order           *reportOrder
		err             error

		// If we don't have any criteria, we return only statistics and
		// no rows.
		hasAnyCriteria = !def.CreatedAfter.IsZero() || !def.CreatedBefore.IsZero() || def.Customer != "" ||
			len(def.OrderCoupons) != 0 || len(def.OrderSources) != 0 || len(def.PaymentTypes) != 0 ||
			len(def.Products) != 0 || len(def.TicketClasses) != 0 || len(def.UsedAtEvents) != 0

		// Initialize the result.
		result = model.ReportResults{
			OrderSources:  make(map[model.OrderSource]int),
			OrderCoupons:  make(model.StringCounts),
			PaymentTypes:  make(model.StringCounts),
			TicketClasses: make(model.StringCounts),
		}

		// Initialize the products and events caches.
		products     = make(map[model.ProductID]*reportProduct)
		usedAtEvents = make(map[model.EventID]*model.ReportEventCount)
	)
	if hasAnyCriteria {
		result.Lines = make([]*model.ReportLine, 0)
	}

	// Because each report needs to return counts of criteria permutations
	// as well as matching records, every report run will inevitably scan
	// the entire database. Given that, it's simpler (and possibly even more
	// efficient) that we generate the report using a single linear scan
	// rather than a bunch of complicated, targeted queries.

	// To begin, we do want to cache the product information since we will
	// be needing random repeated access to it.
	rows, err = tx.tx.Query(
		`SELECT p.id, p.type, p.series, p.name, p.ticket_class, p.ticket_count FROM product p`)
	panicOnError(err)
	for rows.Next() {
		var (
			pid  model.ProductID
			prod reportProduct
		)
		panicOnError(rows.Scan(&pid, &prod.ptype, &prod.series, &prod.name, &prod.tclass, &prod.tcount))
		if prod.tcount == 0 {
			prod.tcount = 1 // so that we can unconditionally multiply and divide by it
		}
		products[pid] = &prod
	}
	panicOnError(rows.Err())

	// Also cache the event information.
	rows, err = tx.tx.Query(`SELECT id, name, series, start FROM event`)
	panicOnError(err)
	for rows.Next() {
		var event model.ReportEventCount
		panicOnError(rows.Scan(&event.ID, &event.Name, &event.Series, (*Time)(&event.Start)))
		usedAtEvents[event.ID] = &event
	}
	panicOnError(rows.Err())
	usedAtEvents[""] = &model.ReportEventCount{Name: "(unused)"}

	// We're also going to want a couple of prepared statements for queries
	// we run often.
	ticketUsageStmt, err = tx.tx.Prepare(`
SELECT COUNT(*), CASE WHEN used!='' THEN event ELSE '' END AS used_event FROM ticket WHERE order_line=? GROUP BY used_event`)
	panicOnError(err)
	defer ticketUsageStmt.Close()
	orderStmt, err = tx.tx.Prepare(`
SELECT o.source, o.name, o.email, o.created, o.flags, o.coupon, p.type, p.subtype
FROM ordert o LEFT JOIN payment p ON p.orderid=o.id AND p.flags&1 WHERE o.id=?`)
	panicOnError(err)
	defer orderStmt.Close()

	// Now, read every order line in the database.  Sort by order ID so that
	// all of the lines for an order are read together.
	rows, err = tx.tx.Query(`SELECT id, orderid, product, quantity, price FROM order_line ORDER BY orderid, id`)
	panicOnError(err)
	for rows.Next() {
		var (
			olid model.OrderLineID
			ol   reportLine
			oid  model.OrderID
		)

		// Get the order line data, and the corresponding order,
		// product, and ticket usage data.
		panicOnError(rows.Scan(&olid, &oid, &ol.pid, &ol.qty, &ol.price))
		if order == nil || order.id != oid {
			var ptype, psubtype sql.NullString
			order = &reportOrder{id: oid}
			panicOnError(orderStmt.QueryRow(oid).Scan(&order.source, &order.name, &order.email,
				(*Time)(&order.created), &order.flags, &order.coupon, &ptype, &psubtype))
			if mapped := paymentTypeMap[ptype.String+","+psubtype.String]; mapped != "" {
				order.paymentType = mapped
			} else if mapped := paymentTypeMap[ptype.String]; mapped != "" {
				order.paymentType = mapped + psubtype.String
			} else {
				order.paymentType = ptype.String + " " + psubtype.String
			}
		}
		if order.flags&model.OrderValid == 0 {
			continue
		}
		ol.order = order
		ol.prod = products[ol.pid]
		if ol.prod.ptype == model.ProdTicket {
			ol.tusage = readTicketUsage(ticketUsageStmt, olid)
		}

		// If all of the criteria match, add this line into the report
		// results.
		if lineMatches(def, &ol, critAll) != 0 {
			if !order.counted {
				result.OrderCount++
				order.counted = true
			}
			if len(def.UsedAtEvents) != 0 {
				for _, eid := range def.UsedAtEvents {
					result.ItemCount += ol.tusage[eid]
					result.TotalAmount += float64(ol.tusage[eid]*ol.price) / float64(ol.prod.tcount) / 100.0
				}
			} else if ol.tusage != nil {
				for _, c := range ol.tusage {
					result.ItemCount += c
					result.TotalAmount += float64(c*ol.price) / float64(ol.prod.tcount) / 100.0
				}
			} else {
				result.ItemCount += ol.qty * ol.prod.tcount
				result.TotalAmount += float64(ol.qty*ol.price) / 100.0
			}
			if result.Lines != nil && len(result.Lines) >= maxReportSize {
				result.Lines = nil
			}
			if result.Lines != nil {
				if len(def.UsedAtEvents) != 0 {
					// One line for each event that tickets were
					// used at (and that was requested in the
					// report), with quantity for that event.
					for _, eid := range def.UsedAtEvents {
						if c := ol.tusage[eid]; c != 0 {
							result.Lines = append(result.Lines, &model.ReportLine{
								OrderID:     ol.order.id,
								OrderTime:   ol.order.created,
								Name:        ol.order.name,
								Email:       ol.order.email,
								Quantity:    c,
								Product:     ol.prod.name,
								UsedAtEvent: eid,
								OrderSource: ol.order.source,
								PaymentType: ol.order.paymentType,
								Amount:      float64(c*ol.price) / float64(ol.prod.tcount) / 100.0,
							})
						}
					}
				} else if ol.tusage != nil {
					// One line for each event that tickets were
					// used at (including not used), with quantity.
					for eid, c := range ol.tusage {
						result.Lines = append(result.Lines, &model.ReportLine{
							OrderID:     ol.order.id,
							OrderTime:   ol.order.created,
							Name:        ol.order.name,
							Email:       ol.order.email,
							Quantity:    c,
							Product:     ol.prod.name,
							UsedAtEvent: eid,
							OrderSource: ol.order.source,
							PaymentType: ol.order.paymentType,
							Amount:      float64(c*ol.price) / float64(ol.prod.tcount) / 100.0,
						})
					}
				} else {
					// One line for the order line.
					result.Lines = append(result.Lines, &model.ReportLine{
						OrderID:     ol.order.id,
						OrderTime:   ol.order.created,
						Name:        ol.order.name,
						Email:       ol.order.email,
						Quantity:    ol.qty * ol.prod.tcount,
						Product:     ol.prod.name,
						OrderSource: ol.order.source,
						PaymentType: ol.order.paymentType,
						Amount:      float64(ol.qty*ol.price) / 100.0,
					})
				}
			}
		}

		// If all or all but one of the criteria match, add this line
		// into the statistics.
		if tcount := lineMatches(def, &ol, critAll&^critOrderSource); tcount != 0 {
			result.OrderSources[ol.order.source] += tcount
		}
		if tcount := lineMatches(def, &ol, critAll&^critOrderCoupon); tcount != 0 {
			result.OrderCoupons[ol.order.coupon] += tcount
		}
		if tcount := lineMatches(def, &ol, critAll&^critProduct); tcount != 0 {
			ol.prod.count += tcount
		}
		if tcount := lineMatches(def, &ol, critAll&^critPaymentType); tcount != 0 {
			result.PaymentTypes[ol.order.paymentType] += tcount
		}
		if tcount := lineMatches(def, &ol, critAll&^critTicketClass); tcount != 0 {
			if ol.prod.ptype == model.ProdTicket {
				result.TicketClasses[ol.prod.tclass] += tcount
			}
		}
		if lineMatches(def, &ol, critAll&^critUsedAtEvent) != 0 {
			if ol.tusage != nil {
				for eid, tcount := range ol.tusage {
					usedAtEvents[eid].Count += tcount
				}
			}
		}
	}
	panicOnError(rows.Err())

	// Generate the result.Products list.
	for pid, prod := range products {
		if prod.count == 0 {
			continue
		}
		result.Products = append(result.Products, &model.ReportProductCount{
			ID:     pid,
			Name:   prod.name,
			Series: prod.series,
			Type:   prod.ptype,
			Count:  prod.count,
		})
	}

	// Generate the result.UsedAtEvents list.
	for _, event := range usedAtEvents {
		if event.Count == 0 {
			continue
		}
		result.UsedAtEvents = append(result.UsedAtEvents, event)
	}
	return &result
}

// readTicketUsage uses the prepared statement to retrieve the ticket usage for
// a particular order line.
func readTicketUsage(stmt *sql.Stmt, olid model.OrderLineID) (usage map[model.EventID]int) {
	var (
		rows  *sql.Rows
		eid   model.EventID
		count int
		err   error
	)
	rows, err = stmt.Query(olid)
	panicOnError(err)
	for rows.Next() {
		panicOnError(rows.Scan(&count, &eid))
		if usage == nil {
			usage = make(map[model.EventID]int)
		}
		usage[eid] = count
	}
	panicOnError(rows.Err())
	return usage
}

// lineMatches determines whether an order line matches the report criteria
// defined in def and selected in crit.  It returns the number of tickets/items
// matched, with 0 meaning no match.
func lineMatches(def *model.ReportDefinition, ol *reportLine, crit byte) int {
	if crit&critOrderSource != 0 && len(def.OrderSources) != 0 {
		var found = false
		for _, os := range def.OrderSources {
			if ol.order.source == os {
				found = true
				break
			}
		}
		if !found {
			return 0
		}
	}
	if crit&critCustomer != 0 && def.Customer != "" {
		if !strings.Contains(strings.ToLower(ol.order.name), def.Customer) &&
			!strings.Contains(strings.ToLower(ol.order.email), def.Customer) {
			return 0
		}
	}
	if crit&critOrderCreated != 0 && !def.CreatedAfter.IsZero() && !ol.order.created.After(def.CreatedAfter) {
		return 0
	}
	if crit&critOrderCreated != 0 && !def.CreatedBefore.IsZero() && !ol.order.created.Before(def.CreatedBefore) {
		return 0
	}
	if crit&critOrderCoupon != 0 && len(def.OrderCoupons) != 0 {
		var found = false
		for _, oc := range def.OrderCoupons {
			if strings.EqualFold(ol.order.coupon, oc) {
				found = true
				break
			}
		}
		if !found {
			return 0
		}
	}
	if crit&critProduct != 0 && len(def.Products) != 0 {
		var found = false
		for _, p := range def.Products {
			if p == ol.pid {
				found = true
				break
			}
		}
		if !found {
			return 0
		}
	}
	if crit&critPaymentType != 0 && len(def.PaymentTypes) != 0 {
		var found = false
		for _, pt := range def.PaymentTypes {
			if pt == ol.order.paymentType {
				found = true
				break
			}
		}
		if !found {
			return 0
		}
	}
	if crit&critTicketClass != 0 && len(def.TicketClasses) != 0 {
		if ol.prod.ptype != model.ProdTicket {
			return 0
		}
		var found = false
		for _, tc := range def.TicketClasses {
			if tc == ol.prod.tclass {
				found = true
				break
			}
		}
		if !found {
			return 0
		}
	}
	if crit&critUsedAtEvent != 0 && len(def.UsedAtEvents) != 0 {
		if ol.tusage == nil {
			return 0
		}
		var count = 0
		for _, eid := range def.UsedAtEvents {
			count += ol.tusage[eid]
		}
		return count
	}
	if ol.tusage == nil {
		return ol.qty * ol.prod.tcount
	}
	var count = 0
	for _, c := range ol.tusage {
		count += c
	}
	return count
}
