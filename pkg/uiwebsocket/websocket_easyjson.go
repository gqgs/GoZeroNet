// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package uiwebsocket

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

func easyjsonC8566e17DecodeGithubComGqgsGoZeronetPkgUiserver(in *jlexer.Lexer, out *Message) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int64(in.Int64())
		case "cmd":
			out.CMD = string(in.String())
		case "wrapper_nonce":
			out.WrapperNonce = string(in.String())
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
func easyjsonC8566e17EncodeGithubComGqgsGoZeronetPkgUiserver(out *jwriter.Writer, in Message) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.ID))
	}
	{
		const prefix string = ",\"cmd\":"
		out.RawString(prefix)
		out.String(string(in.CMD))
	}
	{
		const prefix string = ",\"wrapper_nonce\":"
		out.RawString(prefix)
		out.String(string(in.WrapperNonce))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Message) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC8566e17EncodeGithubComGqgsGoZeronetPkgUiserver(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Message) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC8566e17EncodeGithubComGqgsGoZeronetPkgUiserver(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Message) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC8566e17DecodeGithubComGqgsGoZeronetPkgUiserver(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Message) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC8566e17DecodeGithubComGqgsGoZeronetPkgUiserver(l, v)
}