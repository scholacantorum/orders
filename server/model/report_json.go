package model

import (
	jwriter "github.com/mailru/easyjson/jwriter"
)

func (in ReportResults) MarshalEasyJSON(out *jwriter.Writer) {
	out.RawString(`{"orderCount":`)
	out.Int(in.OrderCount)
	out.RawString(`,"itemCount":`)
	out.Int(in.ItemCount)
	out.RawString(`,"totalAmount":`)
	out.Float64(float64(in.TotalAmount))
	if len(in.Lines) != 0 {
		out.RawString(`,"lines":`)
		{
			out.RawByte('[')
			for v8, v9 := range in.Lines {
				if v8 > 0 {
					out.RawByte(',')
				}
				v9.MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawString(`,"orderSources":`)
	{
		out.RawByte('[')
		v10First := true
		for v10Name, v10Value := range in.OrderSources {
			if v10First {
				v10First = false
			} else {
				out.RawByte(',')
			}
			out.RawString(`{"os":`)
			out.String(string(v10Name))
			out.RawString(`,"c":`)
			out.Int(v10Value)
			out.RawByte('}')
		}
		out.RawByte(']')
	}
	out.RawString(`,"orderCoupons":`)
	in.OrderCoupons.MarshalEasyJSON(out)
	out.RawString(`,"products":`)
	{
		out.RawByte('[')
		for v12, v13 := range in.Products {
			if v12 > 0 {
				out.RawByte(',')
			}
			v13.MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
	out.RawString(`,"paymentTypes":`)
	in.PaymentTypes.MarshalEasyJSON(out)
	out.RawString(`,"ticketClasses":`)
	in.TicketClasses.MarshalEasyJSON(out)
	out.RawString(`,"usedAtEvents":`)
	{
		out.RawByte('[')
		for v16, v17 := range in.UsedAtEvents {
			if v16 > 0 {
				out.RawByte(',')
			}
			v17.MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
	out.RawByte('}')
}

func (in StringCounts) MarshalEasyJSON(out *jwriter.Writer) {
	out.RawByte('[')
	first := true
	for n, c := range in {
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.RawString(`{"n":`)
		out.String(n)
		out.RawString(`,"c":`)
		out.Int(c)
		out.RawByte('}')
	}
	out.RawByte(']')
}

func (in ReportProductCount) MarshalEasyJSON(out *jwriter.Writer) {
	out.RawString(`{"id":`)
	out.String(string(in.ID))
	out.RawString(`,"name":`)
	out.String(in.Name)
	out.RawString(`,"series":`)
	out.String(in.Series)
	out.RawString(`,"ptype":`)
	out.String(string(in.Type))
	out.RawString(`,"count":`)
	out.Int(in.Count)
	out.RawByte('}')
}

func (in ReportLine) MarshalEasyJSON(out *jwriter.Writer) {
	out.RawString(`{"orderID":`)
	out.Int(int(in.OrderID))
	out.RawString(`,"orderTime":`)
	out.Raw((in.OrderTime).MarshalJSON())
	out.RawString(`,"name":`)
	out.String(in.Name)
	out.RawString(`,"email":`)
	out.String(in.Email)
	out.RawString(`,"quantity":`)
	out.Int(in.Quantity)
	out.RawString(`,"product":`)
	out.String(string(in.Product))
	out.RawString(`,"usedAtEvent":`)
	out.String(string(in.UsedAtEvent))
	out.RawString(`,"orderSource":`)
	out.String(string(in.OrderSource))
	out.RawString(`,"paymentType":`)
	out.String(string(in.PaymentType))
	out.RawString(`,"amount":`)
	out.Float64(in.Amount)
	out.RawByte('}')
}

func (in ReportEventCount) MarshalEasyJSON(out *jwriter.Writer) {
	out.RawString(`{"id":`)
	out.String(string(in.ID))
	out.RawString(`,"start":`)
	out.Raw((in.Start).MarshalJSON())
	out.RawString(`,"name":`)
	out.String(in.Name)
	out.RawString(`,"series":`)
	out.String(in.Series)
	out.RawString(`,"count":`)
	out.Int(in.Count)
	out.RawByte('}')
}
