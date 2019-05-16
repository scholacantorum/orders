// Test server for the orders CGI handler.
//
// This program listens on HTTP port 8100, and redirects all requests to the
// orders CGI server (invoked as "./orders").
package main

import (
	"net/http"
	"net/http/cgi"
)

func main() {
	var handler = cgi.Handler{Path: "./orders"}
	http.ListenAndServe("localhost:8100", &handler)
}
