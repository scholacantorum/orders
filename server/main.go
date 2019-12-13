// Main program for the orders.scholacantorum.org server.
//
// This program handles requests to
// https://orders{,-test}.scholacantorum.org/{api,ticket}/*, for management of
// Schola Cantorum product orders and ticket tracking.  It is invoked as a CGI
// "script" by the Dreamhost web server.
//
// This program expects to be run in the web root directory, which must contain
// a mode-700 "data" subdirectory.  The data subdirectory must contain the
// orders.db database and the config.json configuration file.  The server.log
// log file will be created there.
package main

import (
	"fmt"
	"net/http"
	"net/http/cgi"
	"os"
	"strconv"
	"strings"

	"scholacantorum.org/orders/api"
	"scholacantorum.org/orders/auth"
	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/gui"
	"scholacantorum.org/orders/model"
	"scholacantorum.org/orders/ofcapi"
	"scholacantorum.org/orders/payapi"
	"scholacantorum.org/orders/posapi"
)

var (
	txh db.Tx
)

func main() {
	// Change working directory to the data subdirectory of the CGI script
	// location.  This directory should be mode 700 so that it not directly
	// readable by the web server.
	if err := os.Chdir("data"); err != nil {
		fmt.Printf("Status: 500 Internal Server Error\nContent-Type: text/plain\n\n%s\n", err)
		os.Exit(1)
	}
	// Run the request.
	cgi.Serve(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.RunRequest(w, r, txWrapper)
	}))
}

// txWrapper opens the database and wraps the request in a transaction.
func txWrapper(r *api.Request) error {
	// Open the database and start a transaction.
	db.Open("orders.db")
	r.Tx = db.Begin()
	defer func() {
		r.Tx.Rollback()
	}()
	r.Tx.SetRequest(r.Method + " " + r.Path)
	return authWrapper(r)
}

// authWrapper looks for an Auth header in the request and, if present,
// validates the session.
func authWrapper(r *api.Request) error {
	if err := auth.ValidateSession(r); err != nil {
		return err
	}
	return router(r)
}

// router sends the request to the appropriate handler given its method and
// path.
func router(r *api.Request) error {
	c := strings.Split(r.Path[1:], "/")
	for len(c) < 6 {
		c = append(c, "")
	}
	switch {
	case r.Method == "POST" && c[0] == "ofcapi" && c[1] == "event" && c[2] == "":
		return ofcapi.CreateEvent(r)
	case r.Method == "POST" && c[0] == "ofcapi" && c[1] == "login" && c[2] == "":
		return api.Login(r)
	case r.Method == "GET" && c[0] == "ofcapi" && c[1] == "order" && c[2] != "" && c[3] == "":
		if oid, err := strconv.Atoi(c[2]); err == nil && oid > 0 {
			return ofcapi.GetOrder(r, model.OrderID(oid))
		}
	case r.Method == "POST" && c[0] == "ofcapi" && c[1] == "product" && c[2] == "":
		return ofcapi.CreateProduct(r)
	case r.Method == "GET" && c[0] == "ofcapi" && c[1] == "report" && c[2] == "":
		return ofcapi.RunReport(r)
	case r.Method == "POST" && c[0] == "payapi" && c[1] == "order" && c[2] == "":
		return payapi.CreateOrder(r)
	case r.Method == "GET" && c[0] == "payapi" && c[1] == "prices" && c[2] == "":
		return payapi.GetPrices(r)
	case r.Method == "GET" && c[0] == "posapi" && c[1] == "event" && c[2] == "":
		return posapi.ListEvents(r)
	case r.Method == "GET" && c[0] == "posapi" && c[1] == "event" && c[2] != "" && c[3] == "orders" && c[4] == "":
		return posapi.ListEventOrders(r, model.EventID(c[2]))
	case r.Method == "GET" && c[0] == "posapi" && c[1] == "event" && c[2] != "" && c[3] == "prices" && c[4] == "":
		return posapi.GetEventPrices(r, model.EventID(c[2]))
	case r.Method == "GET" && c[0] == "posapi" && c[1] == "event" && c[2] != "" && c[3] == "ticket" && c[4] != "" && c[5] == "":
		return posapi.UseTicket(r, model.EventID(c[2]), c[4])
	case r.Method == "POST" && c[0] == "posapi" && c[1] == "event" && c[2] != "" && c[3] == "ticket" && c[4] != "" && c[5] == "":
		return posapi.UseTicket(r, model.EventID(c[2]), c[4])
	case r.Method == "POST" && c[0] == "posapi" && c[1] == "login" && c[2] == "":
		return api.Login(r)
	case r.Method == "POST" && c[0] == "posapi" && c[1] == "order" && c[2] == "":
		return posapi.CreateOrder(r)
	case r.Method == "DELETE" && c[0] == "posapi" && c[1] == "order" && c[2] != "" && c[3] == "":
		if oid, err := strconv.Atoi(c[2]); err == nil && oid > 0 {
			return posapi.CancelOrder(r, model.OrderID(oid))
		}
	case r.Method == "POST" && c[0] == "posapi" && c[1] == "order" && c[2] != "" && c[3] == "capturePayment" && c[4] == "":
		if oid, err := strconv.Atoi(c[2]); err == nil && oid > 0 {
			return posapi.CaptureOrderPayment(r, model.OrderID(oid))
		}
	case r.Method == "POST" && c[0] == "posapi" && c[1] == "order" && c[2] != "" && c[3] == "sendReceipt" && c[4] == "":
		if oid, err := strconv.Atoi(c[2]); err == nil && oid > 0 {
			return posapi.SendOrderReceipt(r, model.OrderID(oid))
		}
	case r.Method == "GET" && c[0] == "posapi" && c[1] == "stripe" && c[2] == "connectTerminal" && c[3] == "":
		return posapi.GetStripeConnectTerminal(r)
	case r.Method == "GET" && c[0] == "ticket" && c[1] != "" && c[2] == "":
		return gui.ShowTicketInfo(r, c[1])
	}
	return api.NotFound
}

func getComponent(comps []string, index int) string {
	if index >= len(comps) {
		return ""
	}
	return comps[index]
}

func getOrderID(comps []string, index int) model.OrderID {
	if index >= len(comps) {
		return 0
	}
	if val, err := strconv.Atoi(comps[index]); err == nil && val > 0 {
		return model.OrderID(val)
	}
	return model.OrderID(-1)
}
