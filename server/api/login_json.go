// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package api

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson9af65625DecodeScholacantorumOrgOrdersApi(in *jlexer.Lexer, out *ssoLogin) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int(in.Int())
		case "username":
			out.Username = string(in.String())
		case "privSetupOrders":
			out.PrivSetupOrders = bool(in.Bool())
		case "privViewOrders":
			out.PrivViewOrders = bool(in.Bool())
		case "privManageOrders":
			out.PrivManageOrders = bool(in.Bool())
		case "privInPersonSales":
			out.PrivInPersonSales = bool(in.Bool())
		case "privScanTickets":
			out.PrivScanTickets = bool(in.Bool())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson9af65625EncodeScholacantorumOrgOrdersApi(out *jwriter.Writer, in ssoLogin) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != 0 {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.Int(int(in.ID))
	}
	if in.Username != "" {
		const prefix string = ",\"username\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Username))
	}
	if in.PrivSetupOrders {
		const prefix string = ",\"privSetupOrders\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.PrivSetupOrders))
	}
	if in.PrivViewOrders {
		const prefix string = ",\"privViewOrders\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.PrivViewOrders))
	}
	if in.PrivManageOrders {
		const prefix string = ",\"privManageOrders\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.PrivManageOrders))
	}
	if in.PrivInPersonSales {
		const prefix string = ",\"privInPersonSales\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.PrivInPersonSales))
	}
	if in.PrivScanTickets {
		const prefix string = ",\"privScanTickets\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.PrivScanTickets))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ssoLogin) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9af65625EncodeScholacantorumOrgOrdersApi(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ssoLogin) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9af65625EncodeScholacantorumOrgOrdersApi(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ssoLogin) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9af65625DecodeScholacantorumOrgOrdersApi(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ssoLogin) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9af65625DecodeScholacantorumOrgOrdersApi(l, v)
}
func easyjson9af65625DecodeScholacantorumOrgOrdersApi1(in *jlexer.Lexer, out *loginResponse) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "token":
			out.Token = string(in.String())
		case "username":
			out.Username = string(in.String())
		case "stripePublicKey":
			out.StripePublicKey = string(in.String())
		case "privSetupOrders":
			out.PrivSetupOrders = bool(in.Bool())
		case "privViewOrders":
			out.PrivViewOrders = bool(in.Bool())
		case "privManageOrders":
			out.PrivManageOrders = bool(in.Bool())
		case "privInPersonSales":
			out.PrivInPersonSales = bool(in.Bool())
		case "privScanTickets":
			out.PrivScanTickets = bool(in.Bool())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson9af65625EncodeScholacantorumOrgOrdersApi1(out *jwriter.Writer, in loginResponse) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Token != "" {
		const prefix string = ",\"token\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.Token))
	}
	if in.Username != "" {
		const prefix string = ",\"username\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Username))
	}
	if in.StripePublicKey != "" {
		const prefix string = ",\"stripePublicKey\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.StripePublicKey))
	}
	if in.PrivSetupOrders {
		const prefix string = ",\"privSetupOrders\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.PrivSetupOrders))
	}
	if in.PrivViewOrders {
		const prefix string = ",\"privViewOrders\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.PrivViewOrders))
	}
	if in.PrivManageOrders {
		const prefix string = ",\"privManageOrders\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.PrivManageOrders))
	}
	if in.PrivInPersonSales {
		const prefix string = ",\"privInPersonSales\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.PrivInPersonSales))
	}
	if in.PrivScanTickets {
		const prefix string = ",\"privScanTickets\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.PrivScanTickets))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v loginResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9af65625EncodeScholacantorumOrgOrdersApi1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v loginResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9af65625EncodeScholacantorumOrgOrdersApi1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *loginResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9af65625DecodeScholacantorumOrgOrdersApi1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *loginResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9af65625DecodeScholacantorumOrgOrdersApi1(l, v)
}
