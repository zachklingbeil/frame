package zero

import (
	"fmt"
	"html"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type One template.HTML

func NewForge() Forge {
	f := &forge{
		index:   make([]*One, 0),
		Element: NewElement().(*element),
	}
	return f
}

type forge struct {
	index []*One
	Element
}

type Forge interface {
	Build(class string, updateIndex bool, elements ...*One) *One
	JS(js string) One
	CSS(css string) One
	UpdateIndex(*One)
	GetFrame(idx int) *One
	Frames() int
	Count() int
	HandleFrame(w http.ResponseWriter, r *http.Request)
	Element
}

func (f *forge) GetFrame(idx int) *One {
	if idx < 0 || idx >= len(f.index) {
		return nil
	}
	return f.index[idx]
}

func (f *forge) Build(class string, updateIndex bool, elements ...*One) *One {
	var b strings.Builder
	for _, el := range elements {
		b.WriteString(string(*el))
	}

	var htmlOut string
	if class == "" {
		htmlOut = b.String()
	} else {
		consolidatedContent := b.String()
		htmlOut = fmt.Sprintf(`<div class="%s">%s</div>`, html.EscapeString(class), consolidatedContent)
	}
	cleaned := f.consolidateAssets(htmlOut)
	result := One(template.HTML(cleaned))

	if updateIndex {
		f.UpdateIndex(&result)
	}
	return &result
}

func (f *forge) consolidateAssets(html string) string {
	styleRe := regexp.MustCompile(`(?s)<style>(.*?)</style>`)
	styleMatches := styleRe.FindAllStringSubmatch(html, -1)
	var styleBlock string
	if len(styleMatches) > 1 {
		for _, m := range styleMatches {
			styleBlock += m[1] + "\n"
		}
		html = styleRe.ReplaceAllString(html, "")
		if styleBlock != "" {
			html = fmt.Sprintf("<style>%s</style>%s", styleBlock, html)
		}
	}
	scriptRe := regexp.MustCompile(`(?s)<script>(.*?)</script>`)
	scriptMatches := scriptRe.FindAllStringSubmatch(html, -1)
	var scriptBlock string
	if len(scriptMatches) > 1 {
		for _, m := range scriptMatches {
			scriptBlock += m[1] + "\n"
		}
		html = scriptRe.ReplaceAllString(html, "")
		if scriptBlock != "" {
			html = fmt.Sprintf("%s<script>%s</script>", html, scriptBlock)
		}
	}
	return html
}

func (f *forge) JS(js string) One {
	var b strings.Builder
	b.WriteString(`<script>`)
	b.WriteString(js)
	b.WriteString(`</script>`)
	return One(template.HTML(b.String()))
}

func (f *forge) CSS(css string) One {
	var b strings.Builder
	b.WriteString(`<style>`)
	b.WriteString(css)
	b.WriteString(`</style>`)
	return One(template.HTML(b.String()))
}

func (f *forge) Frames() int {
	return int(len(f.index))
}

func (f *forge) UpdateIndex(frame *One) {
	f.index = append(f.index, frame)
}

func (f *forge) Count() int {
	return int(len(f.index))
}

func (f *forge) HandleFrame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Frames", strconv.Itoa(f.Count()))

	current := 0
	if v := r.Header.Get("X-Frame"); v != "" {
		i, err := strconv.Atoi(v)
		if err == nil && i >= 0 && i < f.Count() {
			current = i
		}
	}
	w.Header().Set("X-Frame", strconv.Itoa(current))
	frame := f.GetFrame(current)
	if frame != nil {
		fmt.Fprint(w, *frame)
	}
}
