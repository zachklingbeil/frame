package zero

type Zero interface {
	Fx
	Forge
	Element
}

type zeroImpl struct {
	Fx
	Forge
	Element
}

func NewZero(pathlessUrl, apiUrl string) Zero {
	z := &zeroImpl{
		Fx:      NewFx(pathlessUrl, apiUrl).(*fx),
		Forge:   NewForge().(*forge),
		Element: NewElement().(*element),
	}
	z.Router().HandleFunc("/frame", z.Forge.HandleFrame).Methods("GET", "OPTIONS")
	return z
}
