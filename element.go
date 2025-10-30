package frame

import (
	"fmt"
	"html"
	"html/template"
	"strings"

	math "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	h "github.com/yuin/goldmark/renderer/html"
)

type Element interface {
	Markdown() *goldmark.Markdown
	H1(s string) *One
	H2(s string) *One
	H3(s string) *One
	H4(s string) *One
	H5(s string) *One
	H6(s string) *One
	Paragraph(s string) *One
	Span(s string) *One
	Strong(s string) *One
	Em(s string) *One
	Small(s string) *One
	Mark(s string) *One
	Del(s string) *One
	Ins(s string) *One
	Sub(s string) *One
	Sup(s string) *One
	Kbd(s string) *One
	Samp(s string) *One
	VarElem(s string) *One
	Abbr(s string) *One
	Time(s string) *One
	Button(label string) *One
	Code(code string) *One

	Div(class string) *One
	WrapDiv(class string, children ...*One) *One
	Link(href, text string) *One
	List(items []any, ordered bool) *One
	Img(src, alt, reference string) *One
	Video(src string) *One
	Audio(src string) *One
	Iframe(src string) *One
	Embed(src string) *One
	Source(src string) *One
	Canvas(id string) *One
	Table(cols uint8, rows uint64, data [][]string) *One
}

// --- element Implementation ---
type element struct {
	Md *goldmark.Markdown
}

func NewElement() Element {
	return &element{
		Md: initGoldmark(),
	}
}
func Tag(tag, text string) *One {
	o := One(template.HTML(fmt.Sprintf("<%s>%s</%s>", tag, html.EscapeString(text), tag)))
	return &o
}

func (e *element) Markdown() *goldmark.Markdown {
	return e.Md
}

func (e *element) H1(s string) *One         { return Tag("h1", s) }
func (e *element) H2(s string) *One         { return Tag("h2", s) }
func (e *element) H3(s string) *One         { return Tag("h3", s) }
func (e *element) H4(s string) *One         { return Tag("h4", s) }
func (e *element) H5(s string) *One         { return Tag("h5", s) }
func (e *element) H6(s string) *One         { return Tag("h6", s) }
func (e *element) Paragraph(s string) *One  { return Tag("p", s) }
func (e *element) Span(s string) *One       { return Tag("span", s) }
func (e *element) Strong(s string) *One     { return Tag("strong", s) }
func (e *element) Em(s string) *One         { return Tag("em", s) }
func (e *element) Small(s string) *One      { return Tag("small", s) }
func (e *element) Mark(s string) *One       { return Tag("mark", s) }
func (e *element) Del(s string) *One        { return Tag("del", s) }
func (e *element) Ins(s string) *One        { return Tag("ins", s) }
func (e *element) Sub(s string) *One        { return Tag("sub", s) }
func (e *element) Sup(s string) *One        { return Tag("sup", s) }
func (e *element) Kbd(s string) *One        { return Tag("kbd", s) }
func (e *element) Samp(s string) *One       { return Tag("samp", s) }
func (e *element) VarElem(s string) *One    { return Tag("var", s) }
func (e *element) Abbr(s string) *One       { return Tag("abbr", s) }
func (e *element) Time(s string) *One       { return Tag("time", s) }
func (e *element) Button(label string) *One { return Tag("button", label) }
func (e *element) Code(code string) *One    { return Tag("code", code) }

func (e *element) Div(class string) *One {
	o := One(template.HTML(fmt.Sprintf(`<div class="%s"></div>`, html.EscapeString(class))))
	return &o
}
func (e *element) WrapDiv(class string, children ...*One) *One {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`<div class="%s">`, html.EscapeString(class)))
	for _, child := range children {
		if child != nil {
			b.WriteString(string(*child))
		}
	}
	b.WriteString("</div>")
	o := One(template.HTML(b.String()))
	return &o
}
func (e *element) Link(href, text string) *One {
	o := One(template.HTML(fmt.Sprintf(`<a href="%s">%s</a>`, html.EscapeString(href), html.EscapeString(text))))
	return &o
}

func (e *element) List(items []any, ordered bool) *One {
	tag := "ul"
	if ordered {
		tag = "ol"
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("<%s>", tag))
	for _, item := range items {
		b.WriteString(fmt.Sprintf("<li>%v</li>", html.EscapeString(fmt.Sprintf("%v", item))))
	}
	b.WriteString(fmt.Sprintf("</%s>", tag))
	o := One(template.HTML(b.String()))
	return &o
}

func (e *element) Img(src, alt, reference string) *One {
	styles := "width: 50vw; display: block; margin: 0 auto;"
	switch reference {
	case "large":
		styles = "width: 75vw; display: block; margin: 0 auto;"
	case "small":
		styles = "width: 25vw; display: block; margin: 0 auto;"
	}

	o := One(template.HTML(fmt.Sprintf(
		`<img src="%s" alt="%s" style="%s"/>`,
		html.EscapeString(src),
		html.EscapeString(alt),
		styles,
	)))
	return &o
}

func (e *element) Video(src string) *One {
	o := One(template.HTML(fmt.Sprintf(`<video src="%s"></video>`, html.EscapeString(src))))
	return &o
}

func (e *element) Audio(src string) *One {
	o := One(template.HTML(fmt.Sprintf(`<audio src="%s"></audio>`, html.EscapeString(src))))
	return &o
}

func (e *element) Iframe(src string) *One {
	o := One(template.HTML(fmt.Sprintf(`<iframe src="%s"></iframe>`, html.EscapeString(src))))
	return &o
}

func (e *element) Embed(src string) *One {
	o := One(template.HTML(fmt.Sprintf(`<embed src="%s"/>`, html.EscapeString(src))))
	return &o
}

func (e *element) Source(src string) *One {
	o := One(template.HTML(fmt.Sprintf(`<source src="%s"/>`, html.EscapeString(src))))
	return &o
}

func (e *element) Canvas(id string) *One {
	o := One(template.HTML(fmt.Sprintf(`<canvas id="%s"></canvas>`, html.EscapeString(id))))
	return &o
}

func (e *element) Table(cols uint8, rows uint64, data [][]string) *One {
	var b strings.Builder
	b.WriteString("<table>")
	for _, row := range data {
		b.WriteString("<tr>")
		for _, cell := range row {
			b.WriteString(fmt.Sprintf("<td>%s</td>", html.EscapeString(cell)))
		}
		b.WriteString("</tr>")
	}
	b.WriteString("</table>")
	o := One(template.HTML(b.String()))
	return &o
}

func initGoldmark() *goldmark.Markdown {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, math.MathJax),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
			parser.WithBlockParsers(),
			parser.WithInlineParsers(),
		),
		goldmark.WithRendererOptions(
			h.WithHardWraps(),
			h.WithXHTML(),
		),
	)
	return &md
}
