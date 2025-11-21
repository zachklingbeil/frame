package templates

import _ "embed"

//go:embed css/zero.css
var zeroCSS string

//go:embed css/slides.css
var slidesCSS string

//go:embed css/footer.css
var footerCSS string

//go:embed css/text.css
var textCSS string

//go:embed css/keyboard.css
var keyboardCSS string

type Style interface {
	ZeroCSS() string
	SlidesCSS() string
	FooterCSS() string
	TextCSS() string
	KeyboardCSS() string
}

type style struct{}

func NewStyle() Style {
	return &style{}
}

func (s *style) ZeroCSS() string {
	return zeroCSS
}

func (s *style) SlidesCSS() string {
	return slidesCSS
}

func (s *style) FooterCSS() string {
	return footerCSS
}

func (s *style) TextCSS() string {
	return textCSS
}

func (s *style) KeyboardCSS() string {
	return keyboardCSS
}
