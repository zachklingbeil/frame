package templates

import "github.com/timefactoryio/frame/zero"

type Templates interface {
	Style
	GithubLink(username string) *zero.One
	XLink(username string) *zero.One
	Landing(heading, github, x string)
	README(file string, cssPath string) *zero.One
	ScrollKeybinds() *zero.One
	BuildSlides(dir string) *zero.One
}

type templates struct {
	Style
	zero.Zero
}

func NewTemplates(zero zero.Zero) Templates {
	return &templates{
		Style: NewStyle(),
		Zero:  zero,
	}
}
