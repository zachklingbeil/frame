package fx

import (
	"bytes"
	"compress/gzip"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

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

	f.Router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", contentType)
		w.Write(zipped)
	})
}
