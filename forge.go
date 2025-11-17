package frame

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"html"
	"html/template"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

func NewForge(mux *mux.Router, apiUrl string) Forge {
	if apiUrl == "" {
		apiUrl = "http://localhost:1001"
	} else if !strings.HasPrefix(apiUrl, "http://") && !strings.HasPrefix(apiUrl, "https://") {
		apiUrl = "https://" + apiUrl
	}
	f := &forge{
		Element: NewElement().(*element),
		index:   make([]*One, 0),
		Router:  mux,
		apiUrl:  apiUrl,
	}
	return f
}

type forge struct {
	Element
	index []*One
	*mux.Router
	apiUrl string
}

type Forge interface {
	Build(class string, updateIndex bool, elements ...*One) *One
	JS(js string) One
	CSS(css string) One
	UpdateIndex(*One)
	Count() int
	Keyboard()
	ApiURL() string
	GetFrame(idx int) *One
	AddPath(dir string) string
	AddFile(filePath string, prefix string) error
	Element
	Zero(heading, github, x string)
	TextStyle() *One
	README(file string, cssPath string) *One
	BuildSlides(dir string) *One
	ScrollKeybinds() *One
}

func (f *forge) ApiURL() string {
	return f.apiUrl
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

func (f *forge) Count() int {
	return int(len(f.index))
}

func (f *forge) UpdateIndex(frame *One) {
	f.index = append(f.index, frame)
}

// Add a single file to the frame with a prefix path
func (f *forge) AddFile(filePath string, prefix string) error {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	base := filepath.Base(filePath)
	name := base[:len(base)-len(filepath.Ext(base))]
	contentType := f.getType(base, fileData)

	// Build route path: /prefix/filename
	routePath := "/" + strings.Trim(prefix, "/") + "/" + name

	f.addRoute(routePath, fileData, contentType)
	return nil
}

// Walk directory and load files into memory, determine Content-Type based on file extension.
// Register route using directory name as prefix: /<dirname>/<file without extension>
func (f *forge) AddPath(dir string) string {
	// Get the base directory name to use as prefix
	prefix := filepath.Base(dir)

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		fileData, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		base := filepath.Base(path)
		name := base[:len(base)-len(filepath.Ext(base))]
		contentType := f.getType(base, fileData)
		routePath := "/" + prefix + "/" + name

		f.addRoute(routePath, fileData, contentType)
		return nil
	})
	return prefix
}

func (f *forge) getType(filename string, data []byte) string {
	contentType := mime.TypeByExtension(filepath.Ext(filename))
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}
	return contentType
}

func (f *forge) addRoute(path string, data []byte, contentType string) {
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	gzipWriter.Write(data)
	gzipWriter.Close()
	zipped := buf.Bytes()

	f.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", contentType)
		w.Write(zipped)
	})
}
