package model

import (
	"github.com/mailru/easyjson/jwriter"
)

func (in Session) MarshalEasyJSON(out *jwriter.Writer) {
	out.RawString(`{"token":`)
	out.String(in.Token)
	out.RawString(`,"username":`)
	out.String(in.Username)
	out.RawString(`,"expires":`)
	out.Raw(in.Expires.MarshalJSON())
	if in.Member != 0 {
		out.RawString(`,"member":`)
		out.Int(in.Member)
	}
	out.RawString(`,"privileges":`)
	out.Int(int(in.Privileges))
	out.RawByte('}')
}
