package frame

import (
	"fmt"
	"html"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type One template.HTML

func NewFrame() Frame {
	f := &frame{
		Element: NewElement(),
		Text:    NewText(),
		Index:   make([]*One, 0),
		Router:  mux.NewRouter(),
	}
	return f
}

type Frame interface {
	Zero(src, alt, heading string)
	Build(class string, elements ...*One) *One
	JS(js string) One
	CSS(css string) One
	UpdateIndex(*One)
	Count() int
	Headers(w http.ResponseWriter, r *http.Request)
	Serve()
}

type frame struct {
	Element
	Text
	*mux.Router
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

func (f *frame) Zero(src, alt, heading string) {
	img := f.Element.Img(src, alt, "large")
	h1 := f.Text.H1(heading)

	landingPage := f.Build("", img, h1)
	f.UpdateIndex(landingPage)
}

func (f *frame) Headers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	current := 0
	if v := r.Header.Get("X-Frame"); v != "" {
		if i, err := strconv.Atoi(v); err == nil && i >= 0 && i < f.Count() {
			current = i
		}
	}
	w.Header().Set("X-Index", strconv.Itoa(f.Count()))
	w.Header().Set("X-Frame", strconv.Itoa(current))
	fmt.Fprint(w, *f.Index[current])
}

func (f *frame) Serve() {
	f.Router.HandleFunc("/", f.Headers).Methods("GET")
	go func() {
		http.ListenAndServe(":1002", f.Router)
	}()
}
