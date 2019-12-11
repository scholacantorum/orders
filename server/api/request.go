package api

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path"
	"runtime/debug"
	"syscall"
	"time"

	"scholacantorum.org/orders/db"
	"scholacantorum.org/orders/model"
)

const logfilename = "request.log"

// RunRequest handles an incoming HTTP request, described by the Go standard
// library http.Request and http.ResponseWriter parameters.  It generates a
// Request structure with the request details, passes it to the specified
// RequestHandler, and logs the result.
//
// If the RequestHandler did not emit a response, RunRequest will generate one
// and emit it, as follows:
//   - If the handler returned nil, the generated response is a 204 No Content.
//   - If the handler returned an error that has an HTTPStatusCode() method,
//     such as an error generated by calling HTTPError(), the generated response
//     uses the status code returned by that method, and the error message from
//     the error.
//   - If the handler returned an error that does not have an HTTPStatusCode()
//     method, the generated response is a 400 Bad Request, with the error
//     message from the error.
//   - If the handler panicked, the generated response is a 500 Internal Server
//     Error.
func RunRequest(w http.ResponseWriter, r *http.Request, handler RequestHandler) {
	var (
		request *Request
		err     error
	)
	defer func() {
		var (
			panicked interface{}
			logfile  *os.File
			logmsg   string
			logbuf   bytes.Buffer
			code     int
			end      time.Time
		)
		if panicked = recover(); panicked != nil {
			logmsg = "PANIC"
			code = http.StatusInternalServerError
			err = HTTPError(http.StatusInternalServerError, "Internal Server Error")
		} else if hsce, ok := err.(hasStatusCode); ok {
			logmsg = err.Error()
			code = hsce.StatusCode()
		} else if err != nil {
			logmsg = err.Error()
			code = http.StatusBadRequest
		} else {
			code = http.StatusNoContent
		}
		if request.StatusCode == 0 {
			if code == http.StatusNoContent {
				request.WriteHeader(http.StatusNoContent)
			} else {
				http.Error(request.Response, err.Error(), code)
			}
		}
		if f, ok := request.Response.ResponseWriter.(http.Flusher); ok {
			f.Flush()
		}
		end = time.Now().In(time.Local)
		fmt.Fprintf(&logbuf, "%s %s %s", request.Start.Format("2006-01-02 15:04:05"), request.Method, request.Path)
		if request.Session != nil {
			fmt.Fprintf(&logbuf, " [u=%s]", request.Session.Username)
		}
		if request.Form != nil {
			for k, vs := range request.Form {
				if k == "auth" || k == "password" {
					continue
				}
				for _, v := range vs {
					fmt.Fprintf(&logbuf, " %s=%q", k, v)
				}
			}
		}
		fmt.Fprintf(&logbuf, " => %d", request.StatusCode)
		if logmsg != "" {
			fmt.Fprintf(&logbuf, " (%s)", logmsg)
		}
		fmt.Fprintf(&logbuf, ", %dms\n", end.Sub(request.Start)/time.Millisecond)
		if panicked != nil {
			fmt.Fprintf(&logbuf, string(debug.Stack()))
		}
		if logfile, err = os.OpenFile(logfilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600); err != nil {
			panic(err)
		}
		if err = syscall.Flock(int(logfile.Fd()), syscall.LOCK_EX); err != nil {
			panic(err)
		}
		logbuf.WriteTo(logfile)
		logfile.Close()
	}()
	request = &Request{
		Request: r,
		Response: Response{
			ResponseWriter: w,
		},
		Start: time.Now().In(time.Local),
		Path:  path.Clean("/" + r.URL.Path),
	}
	err = handler(request)
}

// A Request represents an in-progress web request.  It contains the request
// data, the response data, the caller's session data if any, and similar
// tracking data.
type Request struct {
	*http.Request
	Response
	Session    *model.Session
	Privileges model.Privilege
	Start      time.Time
	Path       string
	Tx         db.Tx
}

// Header and Write on Request resolve the ambiguities between http.Request
// and http.ResponseWriter, in favor of the latter.  This allows Request to be
// used in the context of an http.ResponseWriter, such as in calls to
// http.Error.
func (r Request) Header() http.Header           { return r.Response.Header() }
func (r Request) Write(buf []byte) (int, error) { return r.Response.Write(buf) }

// Response is an implementation of http.ResponseWriter that records the status
// code sent in the response.
type Response struct {
	http.ResponseWriter
	StatusCode int
}

func (r Response) Write(buf []byte) (int, error) {
	if r.StatusCode == 0 {
		r.StatusCode = http.StatusOK
	}
	return r.ResponseWriter.Write(buf)
}

// WriteHeader implements http.ResponseWriter.
func (r Response) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// A RequestHandler is a function that takes a Request and returns an error.
// See RunRequest for how this is used.
type RequestHandler func(*Request) error

// HTTPError returns an error containing an HTTP status code and a message.
// This is a common return from a RequestHandler, as the RunRequest function can
// use it to generate an error response to the caller.
func HTTPError(statusCode int, message string) error {
	return &httpError{code: statusCode, msg: message}
}

type httpError struct {
	msg  string
	code int
}

func (he httpError) Error() string   { return he.msg }
func (he httpError) StatusCode() int { return he.code }

type hasStatusCode interface {
	StatusCode() int
}
