package template

import (
	"embed"
	htmlTemplate "html/template"
	"io"
	textTemplate "text/template"
)

var text *textTemplate.Template
var html *htmlTemplate.Template

type tmpl string

const Wrapper tmpl = "wrapper.tmpl"

func (t tmpl) ExecuteHTML(w io.Writer, data interface{}) error {
	return html.ExecuteTemplate(w, string(t), data)
}

func (t tmpl) ExecuteText(w io.Writer, data interface{}) error {
	return text.ExecuteTemplate(w, string(t), data)
}

//go:embed *.tmpl
var fs embed.FS

func init() {
	var err error
	text, err = textTemplate.ParseFS(fs, "*.tmpl")
	if err != nil {
		panic(err)
	}
	html, err = htmlTemplate.ParseFS(fs, "*.tmpl")
	if err != nil {
		panic(err)
	}
}
