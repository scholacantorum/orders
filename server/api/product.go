package api

import (
	"strings"
	"time"

	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
)

// ProductHasCapacity returns false if the product is a ticket for an event that
// is sold out.  It returns true for non-ticket products.
func ProductHasCapacity(tx db.Tx, product *model.Product) bool {
	for _, pe := range product.Events {
		if pe.Priority == 0 {
			if pe.Event.Capacity == 0 {
				return true
			}
			return tx.FetchTicketCount(pe.Event) < pe.Event.Capacity
		}
	}
	return true
}

// MatchingSKU returns true if the SKU's criteria are met.  If future is true,
// SKUs whose sales ranges have not yet started will be accepted.
func MatchingSKU(sku *model.SKU, coupon string, source model.OrderSource, future bool) bool {
	if source != sku.Source {
		return false
	}
	if sku.Coupon != "" && !strings.EqualFold(sku.Coupon, coupon) {
		return false
	}
	return sku.InSalesRange(time.Now()) <= 0
}

// BetterSKU returns the "better" of two SKUs, that is, the one that should be
// used in an order when both are applicable.
func BetterSKU(sku1, sku2 *model.SKU) *model.SKU {
	var (
		isr1, isr2 int
		now        = time.Now()
	)
	// Either of the SKUs might be nil, in which case the other is better.
	if sku2 == nil {
		return sku1
	}
	if sku1 == nil {
		return sku2
	}
	// Check which SKU(s) are in their valid date range.
	isr1 = sku1.InSalesRange(now)
	isr2 = sku2.InSalesRange(now)
	switch {
	case isr1 == 0 && isr2 == 0: // in the date range for both SKUs
		// If one has a better price, return it.
		if sku1.Price < sku2.Price {
			return sku1
		}
		if sku1.Price > sku2.Price {
			return sku2
		}
		// If one has a coupon and the other doesn't, the one with a
		// coupon is better.
		if sku1.Coupon != "" && sku2.Coupon == "" {
			return sku1
		}
		return sku2
	case isr1 == 0 && isr2 != 0: // only sku1 is in range
		return sku1
	case isr1 != 0 && isr2 == 0: // only sku2 is in range
		return sku2
	case isr1 < 0 && isr2 > 0: // sku1 is in the future and sku2 is in the past
		return sku1
	case isr1 > 0 && isr2 < 0: // sku1 is in the past and sku2 is in the future
		return sku2
	case isr1 < 0 && isr2 < 0: // both are in the future
		if sku1.SalesStart.Before(sku2.SalesStart) {
			return sku1
		}
		return sku2
	default: // both are in the past
		if sku1.SalesEnd.After(sku2.SalesEnd) {
			return sku1
		}
		return sku2
	}
}
