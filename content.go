package frame

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Add a single file to the frame with a custom route path
func (f *frame) AddFile(filePath string, routePath string) error {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	base := filepath.Base(filePath)
	contentType := f.getType(base, fileData)

	// If no route path provided, use filename without extension
	if routePath == "" {
		name := base[:len(base)-len(filepath.Ext(base))]
		routePath = "/" + name
	}

	// Ensure route path starts with /
	if !strings.HasPrefix(routePath, "/") {
		routePath = "/" + routePath
	}

	f.addRoute(routePath, fileData, contentType)
	return nil
}

// Walk directory and load files into memory, determine Content-Type based on file extension. Register route/<prefix/<file without extension>.
func (f *frame) AddPath(dir string, prefix string) {
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
		routePath := "/" + strings.Trim(prefix, "/") + "/" + name

		f.addRoute(routePath, fileData, contentType)
		return nil
	})
}

func (f *frame) getType(filename string, data []byte) string {
	contentType := mime.TypeByExtension(filepath.Ext(filename))
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}
	return contentType
}

func (f *frame) addRoute(path string, data []byte, contentType string) {
	f.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		w.Write(data)
	})
}
