// resend-receipt resends the receipt for the order with the specified number.
//
// usage: resend-receipt order-number

package main

import (
	"fmt"
	"os"
	"strconv"

	"scholacantorum.org/orders/api"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
)

func main() {
	var (
		onum  int
		tx    db.Tx
		order *model.Order
		err   error
	)
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: resend-receipt order-number\n")
		os.Exit(2)
	}
	if onum, err = strconv.Atoi(os.Args[1]); err != nil || onum < 1 {
		fmt.Fprintf(os.Stderr, "usage: resend-receipt order-number\n")
		os.Exit(2)
	}
	db.Open("orders.db")
	tx = db.Begin()
	if order = tx.FetchOrder(model.OrderID(onum)); order == nil {
		fmt.Fprintf(os.Stderr, "ERROR: no such order %d\n", onum)
		os.Exit(1)
	}
	if order.Email == "" {
		fmt.Fprintf(os.Stderr, "ERROR: order has no email\n")
		os.Exit(1)
	}
	api.EmitReceipt(order, true)
}
