package md2html

import (
	"html/template"

	"gopkg.in/russross/blackfriday.v2"
)

// Md2Html is return html
func Md2Html(markdown string) template.HTML {
	html := blackfriday.Run([]byte(markdown))
	return template.HTML(html)
}
