package api

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/mailru/easyjson/jwriter"
)

// NotFound is the error returned when a page is not found.
var NotFound = HTTPError(http.StatusNotFound, "404 Not Found")

// SendError sends an error message as a JSON object with an "error" key.  For
// convenience, it also rolls back the transaction.
func SendError(r *Request, message string) {
	r.Header().Set("Content-Type", "application/json")
	w := jwriter.Writer{}
	w.RawString(`{"error":`)
	w.String(message)
	w.RawByte('}')
	w.DumpTo(r)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// NewToken generates a random token string.
func NewToken() string {
	tval := rand.Intn(1000000000000)
	return fmt.Sprintf("%04d-%04d-%04d", tval/100000000, tval/10000%10000, tval%10000)
}
