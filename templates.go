package frame

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"
)

func (f *forge) Zero(heading, github, x string) {
	logo := f.ApiURL() + "/img/logo"
	img := f.Img(logo, "logo", "")
	h1 := f.H1(heading)
	css := f.CSS(`
        .zero {
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            height: 100%;
            width: 100%;
            text-align: center;
            box-sizing: border-box;
            overflow: hidden;
        }
    .zero img {
            max-width: 95%;
            max-height: 30vh;
            width: auto;
            height: auto;
            display: block;
            object-fit: contain;
        }
        .zero h1 {
            color: inherit;
            width: 100%;
            white-space: nowrap;
            overflow: hidden;
            font-size: clamp(2rem, 3vw, 3rem);
            margin: 0;
        }
    `)

	var footer *One
	if github != "" || x != "" {
		footer = f.footer(github, x)
	}

	if footer != nil {
		f.Build("zero", true, &css, img, h1, footer)
	} else {
		f.Build("zero", true, &css, img, h1)
	}
}

func (f *forge) footer(github, x string) *One {
	ghURL := fmt.Sprintf("%s/img/gh", f.ApiURL())
	xURL := fmt.Sprintf("%s/img/x", f.ApiURL())

	var links string
	if github != "" {
		links += fmt.Sprintf(
			`<a href="https://github.com/%s" target="_blank" rel="noopener">
                <img src="%s" alt="GitHub" class="icon" />
            </a>`, github, ghURL)
	}
	if x != "" {
		links += fmt.Sprintf(
			`<a href="https://x.com/%s" target="_blank" rel="noopener">
                <img src="%s" alt="Twitter" class="icon" />
            </a>`, x, xURL)
	}

	footer := One(template.HTML(fmt.Sprintf(`
        <div class="footer">
            %s
        </div>
    `, links)))

	css := f.CSS(`
        .footer {
            display: flex;
            justify-content: center;
            gap: 1.5em;
            margin-top: 1.5em;
        }
        .footer img.icon {
            width: 2em;
            height: 2em;
            object-fit: contain;
        }
    `)
	return f.Build("", false, &css, &footer)
}

func (f *forge) README(file string) *One {
	content, err := os.ReadFile(file)
	if err != nil {
		empty := One("")
		return &empty
	}

	var buf bytes.Buffer
	if err := (*f.Markdown()).Convert(content, &buf); err != nil {
		empty := One("")
		return &empty
	}
	html := `
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/5.5.1/github-markdown-dark.min.css">
<div class="markdown-body">` + buf.String() + `</div>`

	markdown := One(template.HTML(html))
	scroll := f.ScrollKeybinds()
	result := f.Build("text", true, &markdown, scroll)
	return result
}

func (f *forge) BuildMarkdown(file string) *One {
	content, err := os.ReadFile(file)
	if err != nil {
		empty := One("")
		return &empty
	}

	var buf bytes.Buffer
	if err := (*f.Markdown()).Convert(content, &buf); err != nil {
		empty := One("")
		return &empty
	}
	html := buf.String()
	html = strings.ReplaceAll(html, "<p><img", "<img")
	html = strings.ReplaceAll(html, "\"></p>", "\">")
	html = strings.ReplaceAll(html, "\" /></p>", "\" />")
	html = strings.ReplaceAll(html, "\"/></p>", "\"/>")

	markdown := One(template.HTML(html))
	scroll := f.ScrollKeybinds()
	css := f.TextStyle()
	result := f.Build("text", true, &markdown, scroll, css)
	return result
}

func (f *forge) TextStyle() *One {
	css := f.CSS(`
.text {
    display: flex;
    flex-direction: column;
    align-items: center;
    overflow-y: auto;
    scroll-behavior: auto;
    height: 100%;
}
.text img {
    max-width: 90%;
    max-height: 90%;
    object-fit: contain;
    display: block;
    margin: 1em auto;
}
p {
    font-size: 1.2em;
    line-height: 1.5;
    margin: 1em;
    justify-content: flex-start;
    text-align: left;
}

h1 {
    font-size: 3.5em;
    margin-top: 0.5em;
    text-align: center;
}

h2 {
    font-size: 2em;
    border-bottom: 1px solid #33333344;
    padding-bottom: 0.2em;
    margin-top: 1.5em;
    margin-bottom: 0.75em;
    text-align: left;
}

h3 {
    font-size: 1.5em;
    margin-top: 1.2em;
    margin-bottom: 0.5em;
    text-align: left;
    font-weight: 600;
}

h4 {
    font-size: 1em;
    margin-top: 1em;
    margin-bottom: 0.5em;
    text-align: left;
    font-weight: 600;
}

.text > *:not(h1):not(h2):not(h3):not(h4) {
    align-self: stretch;
}
`)
	return &css
}

func (f *forge) ScrollKeybinds() *One {
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
	result := One(template.HTML(fmt.Sprintf(`<script>%s</script>`, js)))
	return &result
}

func (f *forge) BuildSlides(dir string) *One {
	prefix := f.AddPath(dir)
	img := f.Img("", "", "")
	js := f.JS(fmt.Sprintf(`
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
	css := f.CSS(`
.slides {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    height: 100%;
    overflow: hidden;
}
.slides img {
    max-width: 95%;
    max-height: 95%;
    object-fit: contain;
}
    `)
	return f.Build("slides", true, img, &css, &js)
}
