package zero

import (
	_ "embed"
)

//go:embed embed/home.html
var homeHtml string

//go:embed embed/slides.html
var slidesHtml string

//go:embed embed/text.html
var textHtml string

//go:embed embed/keyboard.html
var keyboardHtml string

//go:embed embed/app.html
var appHtml string

type Zero struct {
	SlidesTemplate string
	TextTemplate   string
	HomeTemplate   string
	Keyboard       string
	AppTemplate    string
}

func NewZero() *Zero {
	return &Zero{
		SlidesTemplate: slidesHtml,
		TextTemplate:   textHtml,
		HomeTemplate:   homeHtml,
		Keyboard:       keyboardHtml,
		AppTemplate:    appHtml,
	}
}

func (z *Zero) BuildMaps() {}

type UniverseMap struct {
	APIURL     string
	Focused    uint8
	Layout     uint8
	Variant    uint8
	PrevLayout [2]uint8
}

func (u *UniverseMap) NewUniverseMap(url string) *UniverseMap {
	return &UniverseMap{
		APIURL:     url,
		Focused:    0,
		Layout:     0,
		Variant:    0,
		PrevLayout: [2]uint8{0, 0},
	}
}

type FrameMap struct {
	Source string
	Kb     map[string]string
	Frames [][]byte
}
