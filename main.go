package frame

import (
	"html/template"

	"github.com/gorilla/mux"
	"github.com/timefactoryio/frame/fx"
	"github.com/timefactoryio/frame/templates"
	"github.com/timefactoryio/frame/zero"
)

type One template.HTML

type Frame interface {
	fx.Fx
	templates.Templates
	zero.Forge
}

type frame struct {
	fx.Fx
	templates.Templates
	zero.Forge
	*mux.Router
}

func NewFrame(pathlessUrl, apiURL string) Frame {
	frame := &frame{
		Fx:     fx.NewFx(pathlessUrl, apiURL),
		Router: mux.NewRouter(),
		Forge:  zero.NewForge(),
	}
	frame.Templates = templates.NewTemplates(frame.Forge, frame.Fx)
	frame.HandleFunc("/frame", frame.).Methods("GET", "OPTIONS")

	return f
}
