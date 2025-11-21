package zero

import (
	"bytes"
	"compress/gzip"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Fx interface {
	AddFile(filePath string, prefix string) error
	AddPath(dir string) string
	PathlessUrl() string
	ApiUrl() string
	Serve()
	Router() *mux.Router
}

type fx struct {
	router      *mux.Router
	pathlessUrl string
	apiURL      string
}

func NewFx(pathlessUrl, apiUrl string) Fx {
	f := &fx{
		router:      mux.NewRouter(),
		pathlessUrl: pathlessUrl,
	}
	f.router.Use(f.cors(pathlessUrl))
	return f
}

func (f *fx) Router() *mux.Router {
	return f.router
}

func (f *fx) PathlessUrl() string {
	return f.pathlessUrl
}

func (f *fx) ApiUrl() string {
	return f.apiURL
}

func (f *fx) Serve() {
	go func() {
		http.ListenAndServe(":1001", f.Router())
	}()
}

func (f *fx) cors(pathlessUrl string) mux.MiddlewareFunc {
	origin := "http://pathless:1000"
	if pathlessUrl != "" {
		origin = "https://" + pathlessUrl
	}

	return handlers.CORS(
		handlers.AllowedHeaders([]string{"Content-Type", "X-Frame"}),
		handlers.AllowedOrigins([]string{origin}),
		handlers.AllowedMethods([]string{"GET", "OPTIONS"}),
		handlers.ExposedHeaders([]string{"X-Frame", "X-Frames"}),
	)
}

// Add a single file to the frame with a prefix path
func (f *fx) AddFile(filePath string, prefix string) error {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	base := filepath.Base(filePath)
	name := base[:len(base)-len(filepath.Ext(base))]
	contentType := f.getType(base, fileData)
	routePath := "/" + strings.Trim(prefix, "/") + "/" + name

	f.addRoute(routePath, fileData, contentType)
	return nil
}

// Walk directory and load files into memory, determine Content-Type based on file extension, register routes as /<dirname>/<file without extension>
func (f *fx) AddPath(dir string) string {
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

func (f *fx) getType(filename string, data []byte) string {
	contentType := mime.TypeByExtension(filepath.Ext(filename))
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}
	return contentType
}

func (f *fx) addRoute(path string, data []byte, contentType string) {
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	gzipWriter.Write(data)
	gzipWriter.Close()
	zipped := buf.Bytes()

	f.Router().HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", contentType)
		w.Write(zipped)
	})
}
