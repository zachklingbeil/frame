package frame

import (
	"fmt"
	"html"
	"html/template"
	"strings"
)

type One template.HTML

type Frame interface {
	Build(class string, elements ...One) *One
	JS(js string) One
	CSS(css string) One
	AddFrame(*One)
	Count() uint8
}

func NewFrame() Frame {
	f := &frame{
		element: NewElement().(*element),
		text:    NewText().(*text),
		Index:   make([]*One, 0),
	}
	return f
}

type frame struct {
	*element
	*text
	Index []*One
}

func (f *frame) Count() uint8 {
	return uint8(len(f.Index))
}

func (f *frame) AddFrame(frame *One) {
	f.Index = append(f.Index, frame)
}

func (f *frame) Build(class string, elements ...One) *One {
	var b strings.Builder
	for _, el := range elements {
		b.WriteString(string(el))
	}

	if class == "" {
		result := One(template.HTML(b.String()))
		return &result
	}

	consolidatedContent := template.HTML(b.String())
	htmlResult := fmt.Sprintf(`<div class="%s">%s</div>`, html.EscapeString(class), string(consolidatedContent))
	result := One(template.HTML(htmlResult))
	return &result
}

func (f *frame) JS(js string) One {
	var b strings.Builder
	b.WriteString(`<script>`)
	b.WriteString(js)
	b.WriteString(`</script>`)
	return One(template.HTML(b.String()))
}

func (f *frame) CSS(css string) One {
	var b strings.Builder
	b.WriteString(`<style>`)
	b.WriteString(css)
	b.WriteString(`</style>`)
	return One(template.HTML(b.String()))
}
