package frame

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

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

func NewFrame(pathlessUrl, apiURL string) Frame {
	f := &frame{Router: mux.NewRouter()}
	f.Router.Use(f.cors(pathlessUrl), f.cache())
	f.Forge = NewForge(f.Router, apiURL).(*forge)
	return f
}

func (f *frame) handleFrame(w http.ResponseWriter, r *http.Request) {
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

func (f *frame) Serve() {
	server := &http.Server{
		Addr:         ":1001",
		Handler:      f.Router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("Server error: %v\n", err)
		}
	}()
}

func (f *frame) cors(pathlessUrl string) mux.MiddlewareFunc {
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

func (f *frame) cache() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "max-age=86400, public")
			next.ServeHTTP(w, r)
		})
	}
}
