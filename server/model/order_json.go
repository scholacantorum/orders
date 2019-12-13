package model

import (
	"github.com/mailru/easyjson/jwriter"
)

func (in Card) MarshalEasyJSON(out *jwriter.Writer) {
	out.RawString(`{"card":`)
	out.String(in.Card)
	out.RawString(`,"name":`)
	out.String(in.Name)
	out.RawString(`,"email":`)
	out.String(in.Email)
	out.RawByte('}')
}

func (in Ticket) MarshalEasyJSON(out *jwriter.Writer) {
	out.RawString(`{"id":`)
	out.Int(int(in.ID))
	out.RawString(`,"event":`)
	if in.Event != nil {
		out.String(string(in.Event.ID))
	} else {
		out.RawString("null")
	}
	out.RawString(`,"used":`)
	if !in.Used.IsZero() {
		out.Raw(in.Used.MarshalJSON())
	} else {
		out.RawString("null")
	}
	out.RawByte('}')
}

func (in Payment) MarshalEasyJSON(out *jwriter.Writer) {
	out.RawString(`{"id":`)
	out.Int(int(in.ID))
	out.RawString(`,"type":`)
	out.String(string(in.Type))
	if in.Subtype != "" {
		out.RawString(`,"subtype":`)
		out.String(in.Subtype)
	}
	if in.Method != "" {
		out.RawString(`,"method":`)
		out.String(in.Method)
	}
	if in.Stripe != "" {
		out.RawString(`,"stripe":`)
		out.String(in.Stripe)
	}
	out.RawString(`,"created":`)
	out.Raw(in.Created.MarshalJSON())
	out.RawString(`,"amount":`)
	out.Int(in.Amount)
	out.RawByte('}')
}

func (in OrderLine) MarshalEasyJSON(out *jwriter.Writer) {
	out.RawString(`{"id":`)
	out.Int(int(in.ID))
	out.RawString(`,"product":`)
	out.String(string(in.Product.ID))
	out.RawString(`,"quantity":`)
	out.Int(in.Quantity)
	out.RawString(`,"price":`)
	out.Int(in.Price)
	if len(in.Tickets) != 0 {
		out.RawString(`,"tickets":`)
		{
			out.RawByte('[')
			for v8, v9 := range in.Tickets {
				if v8 > 0 {
					out.RawByte(',')
				}
				v9.MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	if in.Error != "" {
		out.RawString(`,"error":`)
		out.String(in.Error)
	}
	out.RawByte('}')
}

func (in Order) MarshalEasyJSON(out *jwriter.Writer) {
	out.RawString(`{"id":`)
	out.Int(int(in.ID))
	out.RawString(`,"token":`)
	out.String(in.Token)
	out.RawString(`,"valid":`)
	out.Bool(in.Valid)
	out.RawString(`,"source":`)
	out.String(string(in.Source))
	if in.Name != "" {
		out.RawString(`,"name":`)
		out.String(in.Name)
	}
	if in.Email != "" {
		out.RawString(`,"email":`)
		out.String(in.Email)
	}
	if in.Address != "" {
		out.RawString(`,"address":`)
		out.String(in.Address)
	}
	if in.City != "" {
		out.RawString(`,"city":`)
		out.String(in.City)
	}
	if in.State != "" {
		out.RawString(`,"state":`)
		out.String(in.State)
	}
	if in.Zip != "" {
		out.RawString(`,"zip":`)
		out.String(in.Zip)
	}
	if in.Phone != "" {
		out.RawString(`,"phone":`)
		out.String(in.Phone)
	}
	if in.Customer != "" {
		out.RawString(`,"customer":`)
		out.String(in.Customer)
	}
	if in.Member != 0 {
		out.RawString(`,"member":`)
		out.Int(in.Member)
	}
	out.RawString(`,"created":`)
	out.Raw(in.Created.MarshalJSON())
	if in.CNote != "" {
		out.RawString(`,"cNote":`)
		out.String(in.CNote)
	}
	if in.ONote != "" {
		out.RawString(`,"oNote":`)
		out.String(in.ONote)
	}
	if in.InAccess {
		out.RawString(`,"inAccess":`)
		out.Bool(in.InAccess)
	}
	if in.Coupon != "" {
		out.RawString(`,"coupon":`)
		out.String(in.Coupon)
	}
	if len(in.Lines) != 0 {
		out.RawString(`,"lines":`)
		{
			out.RawByte('[')
			for v12, v13 := range in.Lines {
				if v12 > 0 {
					out.RawByte(',')
				}
				v13.MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	if len(in.Payments) != 0 {
		out.RawString(`,"payments":`)
		{
			out.RawByte('[')
			for v14, v15 := range in.Payments {
				if v14 > 0 {
					out.RawByte(',')
				}
				v15.MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}
