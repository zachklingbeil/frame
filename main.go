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
	handleFrame(w http.ResponseWriter, r *http.Request)
	Serve()
	Forge
}

type frame struct {
	Router *mux.Router
	Forge
}

func NewFrame(domain string) Frame {
	f := &frame{Router: mux.NewRouter()}
	f.Router.Use(f.cors(domain))
	f.Forge = NewForge(f.Router).(*forge)
	return f
}

func (f *frame) handleFrame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	current := 0

	if v := r.Header.Get("X-Frame"); v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			http.Error(w, "Invalid X-Frame header value", http.StatusBadRequest)
			return
		}
		if i < 0 || i >= f.Count() {
			http.Error(w, fmt.Sprintf("Frame %d out of range (0-%d)", i, f.Count()-1), http.StatusBadRequest)
			return
		}
		current = i
	} else {
		w.Header().Set("X-Frames", strconv.Itoa(f.Count()))
	}

	w.Header().Set("X-Frame", strconv.Itoa(current))

	frame := f.GetFrame(current)
	if frame != nil {
		fmt.Fprint(w, *frame)
	}
}

func (f *frame) Serve() {
	f.Router.HandleFunc("/frame", f.handleFrame).Methods("GET")
	go func() {
		http.ListenAndServe(":1001", f.Router)
	}()
}

func (f *frame) cors(domain string) mux.MiddlewareFunc {
	origin := "http://localhost:1000"
	if domain != "" {
		origin = "https://" + domain
	}

	return handlers.CORS(
		handlers.AllowedHeaders([]string{"Content-Type", "X-Frame"}),
		handlers.AllowedOrigins([]string{origin}),
		handlers.AllowedMethods([]string{"GET"}),
	)
}
