package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/timefactoryio/frame/fx"

	"github.com/timefactoryio/frame/zero"
)

type Templates interface {
	Style
	zero.Forge
	GithubLink(username string) *zero.One
	XLink(username string) *zero.One
}

type templates struct {
	Style
	zero.Forge
	fx.Fx
}

func NewTemplates(forge zero.Forge, fx fx.Fx) Templates {
	return &templates{
		Style: NewStyle(),
		Forge: forge,
		Fx:    fx,
	}
}

func (t *templates) Zero(heading, github, x string) {
	img := t.Img(t.ApiUrl()+"/img/logo", "")
	h1 := t.H1(heading)
	css := t.CSS(t.ZeroCSS())

	var footer *zero.One
	if github != "" || x != "" {
		footer = t.Footer(
			t.GithubLink(github),
			t.XLink(x),
		)
	}
	t.Build("zero", true, &css, img, h1, footer)
}

func (t *templates) Footer(links ...*zero.One) *zero.One {
	footer := t.Div("footer", links...)
	css := t.CSS(t.FooterCSS())
	return t.Build("", false, &css, footer)
}

func (t *templates) GithubLink(username string) *zero.One {
	ghURL := fmt.Sprintf("%s/img/gh", t.ApiUrl())
	return t.LinkedImg(fmt.Sprintf("https://github.com/%s", username), ghURL, "GitHub")
}

func (t *templates) XLink(username string) *zero.One {
	xURL := fmt.Sprintf("%s/img/x", t.ApiUrl())
	return t.LinkedImg(fmt.Sprintf("https://x.com/%s", username), xURL, "X")
}

func (t *templates) README(file string, cssPath string) *zero.One {
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
	scroll := t.ScrollKeybinds()

	css := t.CSS(t.TextCSS())

	result := t.Build("text", true, &markdown, scroll, &css)
	return result
}

func (t *templates) ScrollKeybinds() *zero.One {
	js := `
(function(){
  const { frame, state } = pathless.context();
  const key = 'scroll';
  
  frame.scrollTop = state[key] || 0;
  
  frame.addEventListener('scroll', () => {
    pathless.update(key, frame.scrollTop);
  });
  
  let speed = 0;
  let isScrolling = false;
  
  const scroll = () => {
    if (speed === 0) {
      isScrolling = false;
      return;
    }
    frame.scrollBy({ top: speed });
    requestAnimationFrame(scroll);
  };
  
  const speeds = { w: -20, s: 20, a: -40, d: 40 };
  pathless.onKey((k) => {
    if (speeds[k]) {
      speed = speeds[k];
      if (!isScrolling) {
        isScrolling = true;
        scroll();
      }
    }
  });
  
  document.addEventListener('keyup', (e) => {
    if (speeds[e.key]) speed = 0;
  });
})();
`
	result := zero.One(template.HTML(fmt.Sprintf(`<script>%s</script>`, js)))
	return &result
}

func (t *templates) BuildSlides(dir string) *zero.One {
	prefix := t.AddPath(dir)
	img := t.Img("", "")
	css := t.CSS(t.SlidesCSS())
	js := t.JS(fmt.Sprintf(`
(function() {
    const { frame, state } = pathless.context();

    let slides = [];
    let index = state.nav || 0;

    async function show(i) {
        if (!slides.length) return;
        index = ((i %% slides.length) + slides.length) %% slides.length;
        pathless.update("nav", index);

        const imgEl = frame.querySelector('img');
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

    pathless.onKey((k) => {
        k = k.toLowerCase();
        if (k === 'a') show(index - 1);
        else if (k === 'd') show(index + 1);
    });
})();
    `, prefix, prefix, prefix, prefix))

	return t.Build("slides", true, img, &css, &js)
}

func (t *templates) Keyboard() {
	css := t.CSS(t.KeyboardCSS())
	js := t.JS(`
(function(){
  const { panel } = pathless.context();
  const keyMap = pathless.keybinds();

  const keys = [
    ['Tab', '', ''],
    ['1', '2', '3'],
    ['q', 'w', 'e'],
    ['a', 's', 'd']
  ];

  const grid = panel.querySelector('.grid');
  if (!grid) return;

  keys.flat().forEach((k) => {
    const entry = keyMap.get(k);
    const keyEl = document.createElement('div');
    keyEl.className = 'key';
    keyEl.dataset.key = k;
    keyEl.textContent = k.toUpperCase();
    if (entry && entry.style) keyEl.style.cssText = entry.style;
    grid.appendChild(keyEl);
  });

  const updateKey = (k, pressed) => {
    const keyEl = grid.querySelector('[data-key="' + k + '"]');
    if (keyEl) keyEl.classList.toggle('pressed', pressed);
  };

  document.addEventListener('keydown', (e) => {
    if (keyMap.has(e.key)) updateKey(e.key, true);
  });

  document.addEventListener('keyup', (e) => {
    if (keyMap.has(e.key)) updateKey(e.key, false);
  });
})();
`)
	html := zero.One(template.HTML(`<div class="grid"></div>`))
	t.Build("keyboard", true, &html, &css, &js)
}
