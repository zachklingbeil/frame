package fx

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Fx interface {
	AddFile(filePath string, prefix string) error
	AddPath(dir string) string
	PathlessUrl() string
	ApiUrl() string
	Serve()
}

type fx struct {
	*mux.Router
	pathlessUrl string
	apiURL      string
}

func NewFx(pathlessUrl, apiUrl string) Fx {
	f := &fx{
		Router:      mux.NewRouter(),
		pathlessUrl: pathlessUrl,
	}
	f.Use(f.cors(pathlessUrl))
	return f
}

func (f *fx) PathlessUrl() string {
	return f.pathlessUrl
}

func (f *fx) ApiUrl() string {
	return f.apiURL
}

func (f *fx) Serve() {
	go func() {
		http.ListenAndServe(":1001", f.Router)
	}()
}

func (f *fx) cors(pathlessUrl string) mux.MiddlewareFunc {
	origin := "http://localhost:1000"
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
