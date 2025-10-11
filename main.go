package frame

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type One template.HTML

func NewFrame() Frame {
	f := &frame{
		element: NewElement().(*element),
		text:    NewText().(*text),
		Index:   make([]*One, 0),
	}
	return f
}

type Frame interface {
	Build(class string, elements ...*One) *One
	JS(js string) One
	CSS(css string) One
	UpdateIndex(*One)
	Count() int
	AddMarkdown(file string) *One
	Headers(w http.ResponseWriter, r *http.Request)
}

type frame struct {
	*element
	*text
	Index []*One
}

func (f *frame) Build(class string, elements ...*One) *One {
	var b strings.Builder
	for _, el := range elements {
		b.WriteString(string(*el))
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

func (f *frame) Count() int {
	return int(len(f.Index))
}

func (f *frame) UpdateIndex(frame *One) {
	f.Index = append(f.Index, frame)
}
func (f *frame) AddMarkdown(file string) *One {
	content, err := os.ReadFile(file)
	if err != nil {
		empty := One("")
		return &empty
	}

	var buf bytes.Buffer
	if err := (*f.Md).Convert(content, &buf); err != nil {
		empty := One("")
		return &empty
	}

	result := One(template.HTML(buf.String()))
	return &result
}

func (f *frame) Zero(src, alt, heading string) {
	img := f.Img(src, alt, "large")
	h1 := f.H1(heading)

	landingPage := f.Build("", img, h1)
	f.UpdateIndex(landingPage)
}

func (f *frame) Headers(w http.ResponseWriter, r *http.Request) {
	count := f.Count()
	if count == 0 {
		http.Error(w, "No frames available", http.StatusNotFound)
		return
	}

	current, err := strconv.Atoi(r.Header.Get("Y"))
	if err != nil || current < 0 || current >= count {
		current = 0
	}

	prev := (current - 1 + count) % count
	next := (current + 1) % count

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X", strconv.Itoa(prev))
	w.Header().Set("Y", strconv.Itoa(current))
	w.Header().Set("Z", strconv.Itoa(next))

	frame := f.Index[current]
	fmt.Fprint(w, *frame)
}
