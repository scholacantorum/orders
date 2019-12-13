package ofcapi

import (
	"errors"
	"strings"
	"time"

	"github.com/mailru/easyjson"

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
	easyjson.MarshalToHTTPResponseWriter(result, r)
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
