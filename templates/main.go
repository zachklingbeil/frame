package templates

import "github.com/timefactoryio/frame/zero"

type Templates interface {
	Style
	GithubLink(username string) *zero.One
	XLink(username string) *zero.One
	Landing(heading, github, x string)
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
