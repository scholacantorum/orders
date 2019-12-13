package model

import (
	jwriter "github.com/mailru/easyjson/jwriter"
)

func (in Event) MarshalEasyJSON(out *jwriter.Writer) {
	out.RawString(`{"id":`)
	out.String(string(in.ID))
	out.RawString(`,"membersID":`)
	out.Int(in.MembersID)
	out.RawString(`,"name":`)
	out.String(in.Name)
	out.RawString(`,"series":`)
	out.String(in.Series)
	out.RawString(`,"start":`)
	out.Raw(in.Start.MarshalJSON())
	out.RawString(`,"capacity":`)
	out.Int(in.Capacity)
	out.RawByte('}')
}
