package frame

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type One template.HTML

func NewFrame(domain string) Frame {
	f := &frame{
		Router: mux.NewRouter(),
	}
	f.Router.Use(f.Cors(domain))
	f.Forge = NewForge(f.Router).(*forge)
	return f
}

type Frame interface {
	Headers(w http.ResponseWriter, r *http.Request)
	Serve()
	Forge
}

type frame struct {
	*mux.Router
	Forge
}

func (f *frame) Headers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	current := 0
	if v := r.Header.Get("X-Frame"); v != "" {
		if i, err := strconv.Atoi(v); err == nil && i >= 0 && i < f.Count() {
			current = i
		}
	}
	w.Header().Set("X-Frames", strconv.Itoa(f.Count()))
	w.Header().Set("X-Frame", strconv.Itoa(current))

	frame := f.GetFrame(current)
	if frame != nil {
		fmt.Fprint(w, *frame)
	}
}

func (f *frame) Serve() {
	f.HandleFunc("/", f.Headers).Methods("GET")
	go func() {
		http.ListenAndServe(":1002", f.Router)
	}()
}

func (f *frame) Cors(domain string) mux.MiddlewareFunc {
	var originValidator func(string) bool
	if domain != "" {
		subdomainPattern := regexp.MustCompile(`^https://([a-zA-Z0-9-]+\.)+` + regexp.QuoteMeta(domain) + `$`)
		originValidator = func(origin string) bool {
			return origin == "https://"+domain || subdomainPattern.MatchString(origin)
		}
	} else {
		originValidator = func(origin string) bool {
			return origin == "http://localhost:1001"
		}
	}
	return handlers.CORS(
		handlers.AllowedHeaders([]string{
			"Content-Type", "X-Frame", "X-Frames", "Cache-Control", "Connection",
		}),
		handlers.AllowedOriginValidator(originValidator),
		handlers.AllowedMethods([]string{"GET"}),
	)
}
