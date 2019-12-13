package model

import (
	"time"

	"github.com/mailru/easyjson/jwriter"
)

func (in AuditRecord) MarshalEasyJSON(out *jwriter.Writer) {
	out.RawByte('{')
	{
		const timeFormat = "2006-01-02 15:04:05"
		out.RawString(`"timestamp":`)
		b := make([]byte, 0, len(timeFormat))
		b = in.Timestamp.In(time.Local).AppendFormat(b, timeFormat)
		out.RawText(b, nil)
	}
	if in.Username != "" {
		out.RawString(`,"username":`)
		out.String(string(in.Username))
	}
	if in.Request != "" {
		out.RawString(`,"request":`)
		out.String(string(in.Request))
	}
	if in.Card != nil {
		out.RawString(`,"card":`)
		in.Card.MarshalEasyJSON(out)
	}
	if in.Event != nil {
		out.RawString(`,"event":`)
		in.Event.MarshalEasyJSON(out)
	}
	if in.Order != nil {
		out.RawString(`,"order":`)
		in.Order.MarshalEasyJSON(out)
	}
	if in.Product != nil {
		out.RawString(`,"product":`)
		in.Product.MarshalEasyJSON(out)
	}
	if in.Session != nil {
		out.RawString(`,"session":`)
		in.Session.MarshalEasyJSON(out)
	}
	out.RawString("}\n")
}
