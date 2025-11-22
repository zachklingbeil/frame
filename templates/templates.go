package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/timefactoryio/frame/zero"
)

func (t *templates) Landing(heading, github, x string) {
	logo := t.ApiUrl() + "/img/logo"
	img := t.Img(logo, "logo")
	h1 := t.H1(heading)
	css := t.CSS(t.ZeroCSS())
	footer := t.buildFooter(github, x)
	t.Build("zero", true, &css, img, h1, footer)
}

func (t *templates) buildFooter(github, x string) *zero.One {
	if github == "" && x == "" {
		return nil
	}

	footerCSS := t.CSS(t.FooterCSS())
	elements := []*zero.One{&footerCSS}

	if github != "" {
		elements = append(elements, t.GithubLink(github))
	}
	if x != "" {
		elements = append(elements, t.XLink(x))
	}
	return t.Build("footer", false, elements...)
}

func (t *templates) GithubLink(username string) *zero.One {
	if username == "" {
		return nil
	}
	logo := fmt.Sprintf("%s/img/gh", t.ApiUrl())
	href := fmt.Sprintf("https://github.com/%s", username)
	return t.LinkedIcon(href, logo, "GitHub")
}

func (t *templates) XLink(username string) *zero.One {
	if username == "" {
		return nil
	}
	logo := fmt.Sprintf("%s/img/x", t.ApiUrl())
	href := fmt.Sprintf("https://x.com/%s", username)
	return t.LinkedIcon(href, logo, "X")
}

func (t *templates) README(file string) *zero.One {
	content, err := os.ReadFile(file)
	if err != nil {
		empty := zero.One("")
		return &empty
	}

	var buf bytes.Buffer
	if err := (*t.Markdown()).Convert(content, &buf); err != nil {
		empty := zero.One("")
		return &empty
	}

	html := buf.String()
	html = strings.ReplaceAll(html, "<p><img", "<img")
	html = strings.ReplaceAll(html, "\"></p>", "\">")
	html = strings.ReplaceAll(html, "\" /></p>", "\" />")
	html = strings.ReplaceAll(html, "\"/></p>", "\"/>")

	markdown := zero.One(template.HTML(html))
	scroll := t.Scroll()

	css := t.CSS(t.TextCSS())

	result := t.Build("text", true, &markdown, scroll, &css)
	return result
}

func (t *templates) Scroll() *zero.One {
	js := t.JS(`pathless.scroll();`)
	return &js
}

func (t *templates) BuildSlides(dir string) *zero.One {
	prefix := t.AddPath(dir)
	img := t.Img("", "")
	css := t.CSS(t.SlidesCSS())
	js := t.JS(fmt.Sprintf(`
(function() {
    let slides = [];
    let index = pathless.state.nav || 0;

    async function show(i) {
        if (!slides.length) return;
        index = ((i %% slides.length) + slides.length) %% slides.length;
        pathless.update("nav", index);

        const imgEl = pathless.frame.querySelector('img');
        if (!imgEl) return;

        const slide = slides[index];
        const fetchKey = '%s.' + slide;
        try {
            const { data } = await pathless.fetch(apiUrl + '/%s/' + slide, { key: fetchKey });
            imgEl.src = data;
            imgEl.alt = slide;
        } catch (e) {
            imgEl.alt = "Failed to load image";
        }
    }

    pathless.fetch(apiUrl + '/%s/order', { key: '%s.order' })
        .then(({ data }) => {
            slides = data || [];
            if (slides.length) show(index);
        });

    pathless.onKey('a', { down: () => show(index - 1) });
    pathless.onKey('d', { down: () => show(index + 1) });
})();
    `, prefix, prefix, prefix, prefix))

	return t.Build("slides", true, img, &css, &js)
}
