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
	"scholacantorum.org/orders/model"
	"scholacantorum.org/orders/ofcapi"
	"scholacantorum.org/orders/payapi"
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
	var (
		pathMatch bool
		comps     = strings.Split(r.Path[1:], "/")
	)
	switch comps[0] {
	case "ofcapi":
		switch getComponent(comps, 1) {
		case "event":
			switch getComponent(comps, 2) {
			case "":
				pathMatch = true
				switch r.Method {
				case http.MethodPost:
					return ofcapi.CreateEvent(r)
				}
			}
		case "login":
			switch getComponent(comps, 2) {
			case "":
				pathMatch = true
				switch r.Method {
				case http.MethodPost:
					return api.Login(r)
				}
			}
		case "order":
			switch oid := getOrderID(comps, 2); oid {
			case 0, -1:
				break
			default:
				pathMatch = true
				switch r.Method {
				case http.MethodGet:
					return ofcapi.GetOrder(r, oid)
				}
			}
		case "product":
			switch getComponent(comps, 2) {
			case "":
				pathMatch = true
				switch r.Method {
				case http.MethodPost:
					return ofcapi.CreateProduct(r)
				}
			}
		case "report":
			switch getComponent(comps, 2) {
			case "":
				pathMatch = true
				switch r.Method {
				case http.MethodGet:
					return ofcapi.RunReport(r)
				}
			}
		}
	case "payapi":
		switch getComponent(comps, 1) {
		case "order":
			switch getOrderID(comps, 2) {
			case 0:
				pathMatch = true
				switch r.Method {
				case http.MethodPost:
					return payapi.CreateOrder(r)
				}
			}
		case "prices":
			switch getComponent(comps, 2) {
			case "":
				pathMatch = true
				switch r.Method {
				case http.MethodGet:
					return payapi.GetPrices(r)
				}
			}
		}
	case "posapi":
		/*
			switch shiftPath(r) {
			case "event":
				switch eventID := shiftPath(r); eventID {
				case "":
					switch r.Method {
					case http.MethodGet:
						posapi.ListEvents(txh, w, r)
					default:
						methodNotAllowedError(txh, w)
					}
				default:
					switch shiftPath(r) {
					case "":
						api.NotFoundError(txh, w)
					case "orders":
						switch shiftPath(r) {
						case "":
							switch r.Method {
							case http.MethodGet:
								posapi.ListEventOrders(txh, w, r, model.EventID(eventID))
							default:
								methodNotAllowedError(txh, w)
							}
						default:
							api.NotFoundError(txh, w)
						}
					case "prices":
						switch shiftPath(r) {
						case "":
							switch r.Method {
							case http.MethodGet:
								posapi.GetEventPrices(txh, w, r, model.EventID(eventID))
							default:
								methodNotAllowedError(txh, w)
							}
						default:
							api.NotFoundError(txh, w)
						}
					case "ticket":
						switch order := shiftPath(r); order {
						case "":
							api.NotFoundError(txh, w)
						default:
							switch r.Method {
							case http.MethodGet, http.MethodPost:
								posapi.UseTicket(txh, w, r, model.EventID(eventID), order)
							default:
								methodNotAllowedError(txh, w)
							}
						}
					default:
						api.NotFoundError(txh, w)
					}
				}
			case "login":
				switch shiftPath(r) {
				case "":
					switch r.Method {
					case http.MethodPost:
						api.Login(txh, w, r)
					default:
						methodNotAllowedError(txh, w)
					}
				default:
					api.NotFoundError(txh, w)
				}
			case "order":
				switch orderID := shiftPathID(r); orderID {
				case 0:
					switch r.Method {
					case http.MethodPost:
						posapi.CreateOrder(txh, w, r)
					default:
						methodNotAllowedError(txh, w)
					}
				default:
					switch shiftPath(r) {
					case "":
						switch r.Method {
						case http.MethodDelete:
							posapi.CancelOrder(txh, w, r, model.OrderID(orderID))
						default:
							methodNotAllowedError(txh, w)
						}
					case "capturePayment":
						switch shiftPath(r) {
						case "":
							switch r.Method {
							case http.MethodPost:
								posapi.CaptureOrderPayment(txh, w, r, model.OrderID(orderID))
							default:
								methodNotAllowedError(txh, w)
							}
						default:
							api.NotFoundError(txh, w)
						}
					case "sendReceipt":
						switch shiftPath(r) {
						case "":
							switch r.Method {
							case http.MethodPost:
								posapi.SendOrderReceipt(txh, w, r, model.OrderID(orderID))
							default:
								methodNotAllowedError(txh, w)
							}
						default:
							api.NotFoundError(txh, w)
						}
					default:
						api.NotFoundError(txh, w)
					}
				case -1:
					api.NotFoundError(txh, w)
				}
			case "stripe":
				switch shiftPath(r) {
				case "connectTerminal":
					switch shiftPath(r) {
					case "":
						switch r.Method {
						case http.MethodGet:
							posapi.GetStripeConnectTerminal(txh, w, r)
						default:
							methodNotAllowedError(txh, w)
						}
					default:
						api.NotFoundError(txh, w)
					}
				default:
					api.NotFoundError(txh, w)
				}
			default:
				api.NotFoundError(txh, w)
			}
		*/
	case "ticket":
		/*
			switch token := shiftPath(r); token {
			case "":
				api.NotFoundError(txh, w)
			default:
				switch shiftPath(r) {
				case "":
					switch r.Method {
					case http.MethodGet:
						gui.ShowTicketInfo(txh, w, r, token)
					default:
						methodNotAllowedError(txh, w)
					}
				default:
					api.NotFoundError(txh, w)
				}
			}
		*/
	}
	if pathMatch {
		return api.HTTPError(http.StatusMethodNotAllowed, "405 Method Not Allowed")
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

/*
func foo() {
	// Next, initialize the logger.
	if logfile, err = os.OpenFile("server.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600); err != nil {
		fmt.Printf("Status: 500 Internal Server Error\nContent-Type: text/plain\n\nserver.log: %s\n", err)
		os.Exit(1)
	}
	log.SetOutput(logfile)
	log.SetFlags(log.Ldate | log.Ltime)
	// Next, make sure that any panic gets logged, and an error returned to
	// the caller.  Also, ensure that the transaction we're about to open,
	// below, gets rolled back if a panic occurs or if it isn't properly
	// closed by the handler.
	defer func() {
		if panicked := recover(); panicked != nil {
			txh.Rollback()
			log.Printf("PANIC: %v", panicked)
			fmt.Fprint(logfile, string(debug.Stack()))
			fmt.Print("Status: 500 Internal Server Error\nContent-Type: text/plain\n\nInternal Server Error\n")
			os.Exit(1)
		} else {
			if err = txh.Rollback(); err != sql.ErrTxDone {
				log.Print("ERROR: transaction not closed")
			}
		}
	}()
	// Next, open the database and start a transaction.
	db.Open("orders.db")
	txh = db.Begin()
	// Finally, handle the request.
	cgi.Serve(http.HandlerFunc(router))
	os.Exit(0)
}

func routerx(w http.ResponseWriter, r *http.Request) {
	r.RequestURI = r.URL.String() // net/http/cgi doesn't set it
	w.Header().Set("Access-Control-Allow-Origin", config.Get("allowOrigin"))
	switch shiftPath(r) {
	case "ofcapi":
		switch shiftPath(r) {
		case "event":
			switch eventID := shiftPath(r); eventID {
			case "":
				switch r.Method {
				case http.MethodPost:
					ofcapi.CreateEvent(txh, w, r)
				default:
					methodNotAllowedError(txh, w)
				}
			default:
				api.NotFoundError(txh, w)
			}
		case "login":
			switch shiftPath(r) {
			case "":
				switch r.Method {
				case http.MethodPost:
					api.Login(txh, w, r)
				default:
					methodNotAllowedError(txh, w)
				}
			default:
				api.NotFoundError(txh, w)
			}
		case "order":
			switch orderID := shiftPathID(r); orderID {
			case 0, -1:
				api.NotFoundError(txh, w)
			default:
				switch r.Method {
				case http.MethodGet:
					ofcapi.GetOrder(txh, w, r, model.OrderID(orderID))
					// Used by members site to validate recording orders
				default:
					methodNotAllowedError(txh, w)
				}
			}
		case "product":
			switch productID := shiftPath(r); productID {
			case "":
				switch r.Method {
				case http.MethodPost:
					ofcapi.CreateProduct(txh, w, r)
				default:
					methodNotAllowedError(txh, w)
				}
			default:
				api.NotFoundError(txh, w)
			}
		case "report":
			switch shiftPath(r) {
			case "":
				switch r.Method {
				case http.MethodGet:
					ofcapi.RunReport(txh, w, r)
				default:
					methodNotAllowedError(txh, w)
				}
			default:
				api.NotFoundError(txh, w)
			}
		default:
			api.NotFoundError(txh, w)
		}
	case "payapi":
		switch shiftPath(r) {
		case "order":
			switch orderID := shiftPathID(r); orderID {
			case 0:
				switch r.Method {
				case http.MethodPost:
					payapi.CreateOrder(txh, w, r)
				default:
					methodNotAllowedError(txh, w)
				}
			default:
				api.NotFoundError(txh, w)
			}
		case "prices":
			switch shiftPath(r) {
			case "":
				switch r.Method {
				case http.MethodGet:
					payapi.GetPrices(txh, w, r)
				default:
					methodNotAllowedError(txh, w)
				}
			default:
				api.NotFoundError(txh, w)
			}
		default:
			api.NotFoundError(txh, w)
		}
	case "posapi":
		switch shiftPath(r) {
		case "event":
			switch eventID := shiftPath(r); eventID {
			case "":
				switch r.Method {
				case http.MethodGet:
					posapi.ListEvents(txh, w, r)
				default:
					methodNotAllowedError(txh, w)
				}
			default:
				switch shiftPath(r) {
				case "":
					api.NotFoundError(txh, w)
				case "orders":
					switch shiftPath(r) {
					case "":
						switch r.Method {
						case http.MethodGet:
							posapi.ListEventOrders(txh, w, r, model.EventID(eventID))
						default:
							methodNotAllowedError(txh, w)
						}
					default:
						api.NotFoundError(txh, w)
					}
				case "prices":
					switch shiftPath(r) {
					case "":
						switch r.Method {
						case http.MethodGet:
							posapi.GetEventPrices(txh, w, r, model.EventID(eventID))
						default:
							methodNotAllowedError(txh, w)
						}
					default:
						api.NotFoundError(txh, w)
					}
				case "ticket":
					switch order := shiftPath(r); order {
					case "":
						api.NotFoundError(txh, w)
					default:
						switch r.Method {
						case http.MethodGet, http.MethodPost:
							posapi.UseTicket(txh, w, r, model.EventID(eventID), order)
						default:
							methodNotAllowedError(txh, w)
						}
					}
				default:
					api.NotFoundError(txh, w)
				}
			}
		case "login":
			switch shiftPath(r) {
			case "":
				switch r.Method {
				case http.MethodPost:
					api.Login(txh, w, r)
				default:
					methodNotAllowedError(txh, w)
				}
			default:
				api.NotFoundError(txh, w)
			}
		case "order":
			switch orderID := shiftPathID(r); orderID {
			case 0:
				switch r.Method {
				case http.MethodPost:
					posapi.CreateOrder(txh, w, r)
				default:
					methodNotAllowedError(txh, w)
				}
			default:
				switch shiftPath(r) {
				case "":
					switch r.Method {
					case http.MethodDelete:
						posapi.CancelOrder(txh, w, r, model.OrderID(orderID))
					default:
						methodNotAllowedError(txh, w)
					}
				case "capturePayment":
					switch shiftPath(r) {
					case "":
						switch r.Method {
						case http.MethodPost:
							posapi.CaptureOrderPayment(txh, w, r, model.OrderID(orderID))
						default:
							methodNotAllowedError(txh, w)
						}
					default:
						api.NotFoundError(txh, w)
					}
				case "sendReceipt":
					switch shiftPath(r) {
					case "":
						switch r.Method {
						case http.MethodPost:
							posapi.SendOrderReceipt(txh, w, r, model.OrderID(orderID))
						default:
							methodNotAllowedError(txh, w)
						}
					default:
						api.NotFoundError(txh, w)
					}
				default:
					api.NotFoundError(txh, w)
				}
			case -1:
				api.NotFoundError(txh, w)
			}
		case "stripe":
			switch shiftPath(r) {
			case "connectTerminal":
				switch shiftPath(r) {
				case "":
					switch r.Method {
					case http.MethodGet:
						posapi.GetStripeConnectTerminal(txh, w, r)
					default:
						methodNotAllowedError(txh, w)
					}
				default:
					api.NotFoundError(txh, w)
				}
			default:
				api.NotFoundError(txh, w)
			}
		default:
			api.NotFoundError(txh, w)
		}
	case "ticket":
		switch token := shiftPath(r); token {
		case "":
			api.NotFoundError(txh, w)
		default:
			switch shiftPath(r) {
			case "":
				switch r.Method {
				case http.MethodGet:
					gui.ShowTicketInfo(txh, w, r, token)
				default:
					methodNotAllowedError(txh, w)
				}
			default:
				api.NotFoundError(txh, w)
			}
		}
	default:
		// This script shouldn't get invoked for anything other than
		// these, because these are the only places the script is
		// installed.
		panic("invalid request URI: " + r.RequestURI)
	}
}

// shiftPath splits off the first component of the request path.  The returned
// component will never contain a slash, and the remaining path will always be a
// rooted path without a trailing slash.
func shiftPath(r *http.Request) (head string) {
	r.URL.Path = path.Clean("/" + r.URL.Path)
	i := strings.Index(r.URL.Path[1:], "/") + 1
	if i <= 0 {
		head, r.URL.Path = r.URL.Path[1:], "/"
	} else {
		head, r.URL.Path = r.URL.Path[1:i], r.URL.Path[i:]
	}
	return head
}

// shiftPathID splits off the first component of the path, and parses it as a
// positive integer.  It returns the integer if successful, 0 if the first path
// component is empty, and -1 if the first path component is not a positive
// integer.
func shiftPathID(r *http.Request) int {
	head := shiftPath(r)
	if head == "" {
		return 0
	}
	if id, err := strconv.Atoi(head); err == nil && id > 0 {
		return id
	}
	return -1
}

// methodNotAllowedError returns an error for a request to a valid URL with an
// invalid method.
func methodNotAllowedError(tx db.Tx, w http.ResponseWriter) {
	tx.Rollback()
	http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
}
*/
