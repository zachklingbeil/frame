package frame

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

type One template.HTML

func NewFrame() Frame {
	f := &frame{
		element: NewElement().(*element),
		text:    NewText().(*text),
		Index:   make([]*One, 0),
		content: map[string][]byte{},
	}
	return f
}

type Frame interface {
	Build(class string, elements ...One) *One
	JS(js string) One
	CSS(css string) One
	UpdateIndex(*One)
	Count() uint8
	AddMarkdown(file string) *One
	AddContent(filePath string, overwrite bool) error
	GetContent(key string) ([]byte, bool)
}

type frame struct {
	*element
	*text
	Index   []*One
	content map[string][]byte
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

func (f *frame) Count() uint8 {
	return uint8(len(f.Index))
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

// Load a single file into memory and store it in a map with the filename (without extension) as the key.
// If overwrite is true, replace existing entries; otherwise skip if the key exists.
func (f *frame) AddContent(filePath string, overwrite bool) error {
	value, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	base := filepath.Base(filePath)
	key := base[:len(base)-len(filepath.Ext(base))]

	if !overwrite {
		if _, exists := f.content[key]; exists {
			return nil // Skip if already exists
		}
	}
	f.content[key] = value
	return nil
}

func (f *frame) GetContent(key string) ([]byte, bool) {
	value, exists := f.content[key]
	return value, exists
}
