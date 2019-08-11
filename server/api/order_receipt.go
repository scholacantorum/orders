package api

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"mime/quotedprintable"
	"net/textproto"
	"os"
	"os/exec"

	"github.com/skip2/go-qrcode"

	"scholacantorum.org/orders/config"
	"scholacantorum.org/orders/model"
)

// EmitReceipt emails a receipt for an order.  This may be for a new order or a
// revised one.  If synch is false, this is an asynchronous operation; the
// actual email delivery is handled by a separate subprocess after this function
// returns.  If synch is true, the email delivery is still handled by a
// separate subprocess but this function waits for it to finish.  Errors are
// logged.
func EmitReceipt(order *model.Order, synch bool) {
	var (
		buf      bytes.Buffer
		mw       *multipart.Writer
		hdr      textproto.MIMEHeader
		htmlw    io.Writer
		htmlqp   io.Writer
		tmpl     *template.Template
		method   string
		img      io.Writer
		b64      io.WriteCloser
		qr       []byte
		cmd      *exec.Cmd
		pipe     io.WriteCloser
		typename string
		ticket   bool
		emailTo  []string
		err      error
	)
	// Can't send a receipt if we don't have an email to send it to.
	if order.Email == "" {
		return
	}
	// Is this "Donation #23" or "Order #23"?  And does it need the QR code
	// for scanning at entry to an event?
	for _, ol := range order.Lines {
		switch ol.Product.Type {
		case model.ProdDonation:
			if typename == "" {
				typename = "Donation"
			}
		case model.ProdTicket:
			ticket = true
			fallthrough
		default:
			typename = "Order"
		}
	}
	// Start a multipart email with appropriate headers.  One part will be
	// the HTML text of the email.  Another part will be the Schola logo
	// header.  The QR code, if included, will be a third part.
	mw = multipart.NewWriter(&buf)
	fmt.Fprint(&buf, "From: Schola Cantorum <admin@scholacantorum.org>\r\n")
	fmt.Fprintf(&buf, "To: %s <%s>\r\n", order.Name, order.Email)
	fmt.Fprint(&buf, "Reply-To: info@scholacantorum.org\r\n")
	fmt.Fprintf(&buf, "Subject: Schola Cantorum %s #%d\r\n", typename, order.ID)
	fmt.Fprintf(&buf, "Content-Type: multipart/related; boundary=%s\r\n\r\n", mw.Boundary())

	// Start the HTML part, including the Schola logo.
	hdr = make(textproto.MIMEHeader)
	hdr.Set("Content-Type", "text/html; charset=UTF-8")
	hdr.Set("Content-Transfer-Encoding", "quoted-printable")
	htmlw, _ = mw.CreatePart(hdr)
	htmlqp = quotedprintable.NewWriter(htmlw)
	fmt.Fprint(htmlqp, `<!DOCTYPE html>
<html><body><body style="margin:0"><div style="width:600px;margin:0 auto"><div style="margin-bottom:24px">
<img src="cid:SCHOLA_LOGO" alt="[Schola Cantorum]" style="border-width:0"></div>`)

	// Include the QR code if we need it.
	if ticket {
		fmt.Fprintf(htmlqp, `<div style="float:right"><a href="%s/ticket/%s">
<img src="cid:ORDER_QRCODE" alt="[Ticket Barcode]" style="border-width:0"></a></div>`,
			config.Get("ordersURL"), order.Token)
	}

	// Greet the customer.
	if order.Name != "" {
		fmt.Fprintf(htmlqp, "<p>Dear %s,</p>", html.EscapeString(order.Name))
	} else {
		fmt.Fprint(htmlqp, "<p>Dear Schola Cantorum Patron,</p>")
	}

	// Add a paragraph for each line on the order.
	for _, ol := range order.Lines {
		if tmpl, err = template.New("t").Funcs(map[string]interface{}{
			"dollars": func(c int) string { return fmt.Sprintf("%.2f", float64(c)/100.0) },
		}).Parse(ol.Product.Receipt); err != nil {
			log.Printf("ERROR: receipt template for %s does not parse: %s", ol.Product.ID, err)
			return
		}
		if err = tmpl.Execute(htmlqp, ol); err != nil {
			log.Printf("ERROR: receipt template for %s failed: %s", ol.Product.ID, err)
			return
		}
	}

	// Add a paragraph with a line for each payment.
	for i, p := range order.Payments {
		if i == 0 {
			fmt.Fprint(htmlqp, "<p>")
		} else {
			fmt.Fprint(htmlqp, "<br>")
		}
		if p.Type == model.PaymentOther {
			method = p.Subtype
		} else {
			method = p.Method
		}
		if p.Amount > 0 {
			fmt.Fprintf(htmlqp, "You paid $%.2f on %s via %s.", float64(p.Amount)/100.0,
				p.Created.Format("January 2, 2006 at 3:04pm"), method)
		} else {
			fmt.Fprintf(htmlqp, "You were refunded $%.2f on %s via %s.", -float64(p.Amount)/100.0,
				p.Created.Format("January 2, 2006 at 3:04pm"), method)
		}
	}
	if len(order.Payments) != 0 {
		fmt.Fprint(htmlqp, "</p>")
	}

	// Add the signature and close out the HTML.
	fmt.Fprint(htmlqp, `<p>Sincerely yours,<br>Schola Cantorum</p>
<p>Web: <a href="https://scholacantorum.org">scholacantorum.org</a><br>
Email: <a href="mailto:info@scholacantorum.org">info@scholacantorum.org</a><br>
Phone: (650) 254-1700</p></div></body><html>
`)

	// Add the Schola logo.
	hdr = make(textproto.MIMEHeader)
	hdr.Set("Content-Type", "image/gif")
	hdr.Set("Content-Transfer-Encoding", "base64")
	hdr.Set("Content-ID", "<SCHOLA_LOGO>")
	img, _ = mw.CreatePart(hdr)
	img.Write(mailLogo)

	// If we need the QR code, add it too.
	if ticket {
		hdr = make(textproto.MIMEHeader)
		hdr.Set("Content-Type", "image/png")
		hdr.Set("Content-Transfer-Encoding", "base64")
		hdr.Set("Content-ID", "<ORDER_QRCODE>")
		img, _ = mw.CreatePart(hdr)
		if qr, err = qrcode.Encode(
			fmt.Sprintf("%s/ticket/%s", config.Get("ordersURL"), order.Token),
			qrcode.Highest, 200); err != nil {
			log.Printf("ERROR: can't create QR code for order %d: %s", order.ID, err)
			return
		}
		b64 = base64.NewEncoder(base64.StdEncoding, img)
		b64.Write(qr)
		b64.Close()
	}
	mw.Close()

	emailTo = []string{"admin@scholacantorum.org"}
	if config.Get("mode") != "development" {
		emailTo = append(emailTo, order.Email)
	}
	if config.Get("mode") == "production" {
		emailTo = append(emailTo, "info@scholacantorum.org")
	}
	cmd = exec.Command(config.Get("bin")+"/send-raw-email", emailTo...)
	if synch {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		// For asynch, they go to /dev/null, which is necessary so that
		// this parent process can exit (and the CGI caller can get a
		// response) before the child finishes.
	}
	if pipe, err = cmd.StdinPipe(); err != nil {
		log.Printf("ERROR: can't send receipt for order %d: can't send email: %s", order.ID, err)
		return
	}
	if err = cmd.Start(); err != nil {
		log.Printf("ERROR: can't send receipt for order %d: can't send email: %s", order.ID, err)
		return
	}
	pipe.Write(buf.Bytes())
	pipe.Close()
	if synch {
		if err = cmd.Wait(); err != nil {
			log.Printf("ERROR: can't send receipt for order %d: can't send email: %s", order.ID, err)
		}
	}
	// Otherwise we are intentionally not waiting for the subprocess to
	// finish.  This CGI script will exit immediately, so that the user gets
	// a fast response to their order.  The subprocess will continue as an
	// orphan until the email is sent, and its zombie will be reaped by the
	// init daemon.
}
