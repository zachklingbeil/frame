package frame

import (
	"fmt"
	"html"
	"html/template"
	"strings"
)

type Element interface {
	Div(class string) *One
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
type element struct{}

func NewElement() Element {
	return &element{}
}
func Tag(tag, text string) *One {
	o := One(template.HTML(fmt.Sprintf("<%s>%s</%s>", tag, html.EscapeString(text), tag)))
	return &o
}

func (e *element) Div(class string) *One {
	o := One(template.HTML(fmt.Sprintf(`<div class="%s"></div>`, html.EscapeString(class))))
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
