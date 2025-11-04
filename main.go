package frame

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type One template.HTML

type Frame interface {
	HandleFrame(w http.ResponseWriter, r *http.Request)
	Serve()
	Forge
}

type frame struct {
	*mux.Router
	Forge
}

func NewFrame(domain string) Frame {
	f := &frame{Router: mux.NewRouter()}
	f.Router.Use(f.cors(domain))
	f.Forge = NewForge(f.Router).(*forge)
	return f
}

func (f *frame) HandleFrame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Parse frame index from header, default to 0
	index := 0
	if v := r.Header.Get("X-Frame"); v != "" {
		if i, err := strconv.Atoi(v); err == nil && i >= 0 && i < f.Count() {
			index = i
		}
	}
	// Send total frames only on frame 0
	if index == 0 {
		w.Header().Set("X-Frames", strconv.Itoa(f.Count()))
	}

	frame := f.GetFrame(index)
	if frame == nil {
		http.Error(w, "Frame not found", http.StatusNotFound)
		return
	}
	fmt.Fprint(w, *frame)
}

func (f *frame) Serve() {
	f.HandleFunc("/frame", f.HandleFrame).Methods("GET")
	go func() {
		http.ListenAndServe(":1002", f.Router)
	}()
}

func (f *frame) cors(domain string) mux.MiddlewareFunc {
	origin := "http://localhost:1001"
	if domain != "" {
		origin = "https://" + domain
	}

	return handlers.CORS(
		handlers.AllowedHeaders([]string{"Content-Type", "X-Frame"}),
		handlers.AllowedOrigins([]string{origin}),
		handlers.AllowedMethods([]string{"GET"}),
	)
}
