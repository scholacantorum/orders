package model

import (
	jwriter "github.com/mailru/easyjson/jwriter"
)

func (in SKU) MarshalEasyJSON(out *jwriter.Writer) {
	out.RawByte('{')
	out.RawString(`{"source":`)
	out.String(string(in.Source))
	if in.Coupon != "" {
		out.RawString(`,"coupon":`)
		out.String(in.Coupon)
	}
	if !in.SalesStart.IsZero() {
		out.RawString(`,"salesStart":`)
		out.Raw((in.SalesStart).MarshalJSON())
	}
	if !in.SalesEnd.IsZero() {
		out.RawString(`,"salesEnd":`)
		out.Raw((in.SalesEnd).MarshalJSON())
	}
	out.RawString(`,"price":`)
	out.Int(in.Price)
	out.RawByte('}')
}

func (in ProductEvent) MarshalEasyJSON(out *jwriter.Writer) {
	out.RawString(`{"priority":`)
	out.Int(in.Priority)
	out.RawString(`,"event":`)
	out.String(string(in.Event.ID))
	out.RawByte('}')
}

func (in Product) MarshalEasyJSON(out *jwriter.Writer) {
	out.RawString(`{"id":`)
	out.String(string(in.ID))
	if in.Series != "" {
		out.RawString(`,"series":`)
		out.String(in.Series)
	}
	out.RawString(`,"name":`)
	out.String(in.Name)
	out.RawString(`,"shortName":`)
	out.String(in.ShortName)
	out.RawString(`,"type":`)
	out.String(string(in.Type))
	if in.Receipt != "" {
		out.RawString(`,"receipt":`)
		out.String(in.Receipt)
	}
	if in.TicketCount != 0 {
		out.RawString(`,"ticketCount":`)
		out.Int(in.TicketCount)
	}
	if in.TicketClass != "" {
		out.RawString(`,"ticketClass":`)
		out.String(in.TicketClass)
	}
	if len(in.SKUs) != 0 {
		out.RawString(`,"skus":`)
		{
			out.RawByte('[')
			for v3, v4 := range in.SKUs {
				if v3 > 0 {
					out.RawByte(',')
				}
				v4.MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	if len(in.Events) != 0 {
		out.RawString(`,"events":`)
		{
			out.RawByte('[')
			for v5, v6 := range in.Events {
				if v5 > 0 {
					out.RawByte(',')
				}
				v6.MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}
