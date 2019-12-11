package ofcapi

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/rothskeller/json"

	"scholacantorum.org/orders/api"
	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/model"
)

// RunReport runs a report, handling GET /ofcapi/report requests.
func RunReport(r *api.Request) error {
	var (
		def    *model.ReportDefinition
		result *model.ReportResults
	)
	// Verify permissions.
	if r.Privileges&model.PrivViewOrders == 0 {
		return auth.Forbidden
	}
	// Get the report definition.
	if def = parseReportDef(r); def == nil {
		return errors.New("invalid report definition")
	}
	result = r.Tx.RunReport(def)
	r.Tx.Commit()
	// Send back the results.
	r.Header().Set("Content-Type", "application/json")
	r.Write(result.ToJSON())
	return nil
}

func parseReportDef(r *api.Request) (def *model.ReportDefinition) {
	def = new(model.ReportDefinition)
	r.ParseForm()
	for _, os := range r.Form["orderSource"] {
		switch os := model.OrderSource(os); os {
		case model.OrderFromPublic, model.OrderFromMembers, model.OrderFromGala, model.OrderFromOffice, model.OrderInPerson:
			def.OrderSources = append(def.OrderSources, os)
		default:
			return nil
		}
	}
	def.Customer = strings.ToLower(strings.TrimSpace(r.FormValue("customer")))
	if before := r.FormValue("createdBefore"); before != "" {
		if t, err := time.ParseInLocation("2006-01-02T15:04:05", before, time.Local); err != nil {
			return nil
		} else {
			def.CreatedBefore = t
		}
	}
	if after := r.FormValue("createdAfter"); after != "" {
		if t, err := time.ParseInLocation("2006-01-02T15:04:05", after, time.Local); err != nil {
			return nil
		} else {
			def.CreatedAfter = t
		}
	}
	for _, v := range r.Form["orderCoupon"] {
		if v := strings.TrimSpace(v); v != "" {
			def.OrderCoupons = append(def.OrderCoupons, strings.ToUpper(v))
		}
	}
	for _, pid := range r.Form["product"] {
		if r.Tx.FetchProduct(model.ProductID(pid)) == nil {
			return nil
		}
		def.Products = append(def.Products, model.ProductID(pid))
	}
	def.PaymentTypes = r.Form["paymentType"]
	def.TicketClasses = r.Form["ticketClass"]
	for _, eid := range r.Form["usedAtEvent"] {
		if eid != "" && r.Tx.FetchEvent(model.EventID(eid)) == nil {
			return nil
		}
		def.UsedAtEvents = append(def.UsedAtEvents, model.EventID(eid))
	}
	return def
}

func emitReport(w http.ResponseWriter, result *model.ReportResults) {
	var jw = json.NewWriter(w)
	jw.Object(func() {
		if result.Lines != nil {
			jw.Prop("lines", func() {
				jw.Array(func() {
					for _, r := range result.Lines {
						jw.Object(func() {
							jw.Prop("orderID", int(r.OrderID))
							jw.Prop("orderTime", r.OrderTime.Format(time.RFC3339))
							jw.Prop("name", r.Name)
							jw.Prop("email", r.Email)
							jw.Prop("quantity", r.Quantity)
							jw.Prop("product", r.Product)
							jw.Prop("usedAtEvent", string(r.UsedAtEvent))
							jw.Prop("orderSource", string(r.OrderSource))
							jw.Prop("paymentType", r.PaymentType)
							jw.Prop("amount", r.Amount)
						})
					}
				})
			})
		}
		jw.Prop("orderCount", result.OrderCount)
		jw.Prop("itemCount", result.ItemCount)
		jw.Prop("totalAmount", result.TotalAmount)
		jw.Prop("orderSources", func() {
			jw.Array(func() {
				for os, c := range result.OrderSources {
					jw.Object(func() {
						jw.Prop("os", string(os))
						jw.Prop("c", c)
					})
				}
			})
		})
		jw.Prop("orderCoupons", func() {
			emitStringCounts(jw, result.OrderCoupons)
		})
		jw.Prop("products", func() {
			jw.Array(func() {
				for _, prod := range result.Products {
					jw.Object(func() {
						jw.Prop("id", string(prod.ID))
						jw.Prop("name", prod.Name)
						jw.Prop("series", prod.Series)
						jw.Prop("ptype", string(prod.Type))
						jw.Prop("count", prod.Count)
					})
				}
			})
		})
		jw.Prop("paymentTypes", func() {
			emitStringCounts(jw, result.PaymentTypes)
		})
		jw.Prop("ticketClasses", func() {
			emitStringCounts(jw, result.TicketClasses)
		})
		jw.Prop("usedAtEvents", func() {
			jw.Array(func() {
				for _, event := range result.UsedAtEvents {
					jw.Object(func() {
						jw.Prop("id", string(event.ID))
						jw.Prop("series", event.Series)
						if event.Start.IsZero() {
							jw.Prop("start", "")
						} else {
							jw.Prop("start", event.Start.Format(time.RFC3339))
						}
						jw.Prop("name", event.Name)
						jw.Prop("count", event.Count)
					})
				}
			})
		})
	})
	jw.Close()
}

func emitStringCounts(jw json.Writer, sc model.StringCounts) {
	jw.Array(func() {
		for s, c := range sc {
			jw.Object(func() {
				jw.Prop("n", s)
				jw.Prop("c", c)
			})
		}
	})
}
