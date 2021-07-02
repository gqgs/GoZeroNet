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

func easyjsonC8566e17DecodeGithubComGqgsGoZeronetPkgUiwebsocket(in *jlexer.Lexer, out *serverErrorRsponse) {
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
		case "error":
			out.Error = string(in.String())
		case "cmd":
			out.CMD = string(in.String())
		case "id":
			out.ID = int64(in.Int64())
		case "to":
			out.To = int64(in.Int64())
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
func easyjsonC8566e17EncodeGithubComGqgsGoZeronetPkgUiwebsocket(out *jwriter.Writer, in serverErrorRsponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"error\":"
		out.RawString(prefix[1:])
		out.String(string(in.Error))
	}
	{
		const prefix string = ",\"cmd\":"
		out.RawString(prefix)
		out.String(string(in.CMD))
	}
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix)
		out.Int64(int64(in.ID))
	}
	{
		const prefix string = ",\"to\":"
		out.RawString(prefix)
		out.Int64(int64(in.To))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v serverErrorRsponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC8566e17EncodeGithubComGqgsGoZeronetPkgUiwebsocket(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v serverErrorRsponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC8566e17EncodeGithubComGqgsGoZeronetPkgUiwebsocket(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *serverErrorRsponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC8566e17DecodeGithubComGqgsGoZeronetPkgUiwebsocket(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *serverErrorRsponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC8566e17DecodeGithubComGqgsGoZeronetPkgUiwebsocket(l, v)
}
func easyjsonC8566e17DecodeGithubComGqgsGoZeronetPkgUiwebsocket1(in *jlexer.Lexer, out *Message) {
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
		case "to":
			out.To = int64(in.Int64())
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
func easyjsonC8566e17EncodeGithubComGqgsGoZeronetPkgUiwebsocket1(out *jwriter.Writer, in Message) {
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
	{
		const prefix string = ",\"to\":"
		out.RawString(prefix)
		out.Int64(int64(in.To))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Message) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC8566e17EncodeGithubComGqgsGoZeronetPkgUiwebsocket1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Message) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC8566e17EncodeGithubComGqgsGoZeronetPkgUiwebsocket1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Message) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC8566e17DecodeGithubComGqgsGoZeronetPkgUiwebsocket1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Message) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC8566e17DecodeGithubComGqgsGoZeronetPkgUiwebsocket1(l, v)
}
