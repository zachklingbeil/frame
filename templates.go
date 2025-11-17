package frame

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"
)

func (f *forge) Keyboard() {
	js := f.JS(`
(function(){
  const { panel, state } = pathless.context();
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
    keyEl.textContent = k === 'Tab' ? '⇥' : k;
    
    if (entry && entry.style) {
      keyEl.style.cssText = entry.style;
    }
    
    grid.appendChild(keyEl);
  });
  
  const updateKey = (k, pressed) => {
    const keyEl = grid.querySelector('[data-key="' + k + '"]');
    if (keyEl) {
      keyEl.classList.toggle('pressed', pressed);
    }
  };
  
  document.addEventListener('keydown', (e) => {
    if (keyMap.has(e.key)) {
      updateKey(e.key, true);
    }
  });
  
  document.addEventListener('keyup', (e) => {
    if (keyMap.has(e.key)) {
      updateKey(e.key, false);
    }
  });
})();
`)
	css := f.CSS(`
.keyboard {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    height: 100%;
    background: #111;
    border-radius: 0.75em;
    box-shadow: 0 0.25em 1.5em #000a;
    padding: 1em;
}
.grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    grid-template-rows: repeat(4, 1fr);
    gap: 0.5em;
    width: 100%;
    max-width: 600px;
}
.key {
    border: medium solid #444;
    border-radius: 0.375em;
    height: 4em;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 600;
    font-size: 1.3em;
    background: #222;
    color: white;
}
.key.pressed {
    border-color: #fff;
    background: #333;
}
.key:empty {
    opacity: 0;
    pointer-events: none;
}
`)
	html := One(template.HTML(`<div class="grid"></div>`))
	f.Build("keyboard", true, &html, &css, &js)
}

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

func (f *forge) README(file string, cssPath string) *One {
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

	var css *One
	if cssPath == "" {
		css = f.TextStyle()
	} else {
		cssContent := ""
		if b, err := os.ReadFile(cssPath); err == nil {
			cssContent = string(b)
		}
		c := f.CSS(cssContent)
		css = &c
	}

	result := f.Build("text", true, &markdown, scroll, css)
	return result
}

func (f *forge) TextStyle() *One {
	css := f.CSS(`
.text {
	display: flex;
	flex-direction: column;
	overflow-y: auto;
	height: 100%;
	padding: max(2vw, 1.5em);
	max-width: 70vw;
	margin: 0 auto;
	box-sizing: border-box;
	color: #f3f3f3;
}
.text img {
	max-width: 90%;
	max-height: 60vh;
	object-fit: contain;
	display: block;
	margin: 1.2em auto;
	border-radius: 0.4em;
}
.text h1,
.text h2,
.text h3,
.text h4 {
	font-weight: 700;
	line-height: 1.2;
	margin-top: 1.5em;
	margin-bottom: 0.5em;
	text-align: center;
	letter-spacing: 0.01em;
}
.text h1 {
	font-size: 3em;
	margin-top: 1em;
	margin-bottom: 0.3em;
}
.text h2 {
	font-size: 2em;
	border-bottom: 2px solid #ffffff33;
	padding-bottom: 0.2em;
}
.text h3 {
	font-size: 1.3em;
	text-align: left;
}
.text h4 {
	font-size: 1em;
	text-align: left;
}

.text p {
	font-size: 1.15em;
	margin: 1em 0;
	text-align: left;
	line-height: 1.7;
}
.text h1 + p,
.text h2 + p {
	text-align: center;
	margin-top: 0;
}

.text ul,
.text ol {
	margin: 0 0 1.2em 2em;
	padding-left: 1.2em;
	list-style-position: outside;
	text-align: left;
}
.text ul li,
.text ol li {
	margin-bottom: 0.4em;
}
.text code {
	border: 1px solid #ffffff58;
	padding: 0.2em 0.5em;
	border-radius: 0.3em;
	font-size: 0.98em;
	font-family: 'Fira Mono', 'Consolas', monospace;
}
.text pre {
	margin-left: auto;
	margin-right: auto;
	display: flex;
	justify-content: center;
	max-width: 90vw;
	width: fit-content;
}
.text pre code {
	display: block;
	background: none;
	border: none;
	color: inherit;
	padding: 0;
	border-radius: 0;
	font-size: 1em;
	white-space: pre;
	overflow-x: auto;
	max-width: 100vw;
}
.text table {
	border-collapse: separate;
	border-spacing: 0;
	margin: 2em auto;
	width: auto;
	min-width: 10%;
	max-width: 90%;
	font-size: 1em;
	border-radius: 0.4em;
	overflow: hidden;
	background: #202020;
	box-shadow: 0 2px 12px #0005;
}
.text th,
.text td {
	border: 1px solid #333;
	padding: 0.7em 1.2em;
	text-align: center;
	vertical-align: middle;
}
.text th {
	font-weight: 700;
	background: #232323;
}
.text tr:nth-child(even) td {
	background: #181818;
}
.text tr:hover td {
	background: #232323;
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
