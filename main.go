package frame

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

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
	f.Router.Use(f.middleware(pathlessUrl))
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

func (f *frame) middleware(pathlessUrl string) mux.MiddlewareFunc {
	origin := "http://localhost:1000"
	if pathlessUrl != "" {
		origin = "https://" + pathlessUrl
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Frame")
			w.Header().Set("Access-Control-Expose-Headers", "X-Frame, X-Frames")

			w.Header().Set("Cache-Control", "max-age=86400, public")
			w.Header().Set("Connection", "keep-alive")
			w.Header().Set("Keep-Alive", "timeout=120, max=100")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
