package frame

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"
)

func (f *forge) Zero(src, heading string) {
	img := f.Img(src, "logo", "large")
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
			max-height: 95%;
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
	f.Build("zero", true, &css, img, h1)
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
    justify-content: center;
}

h2, h3, h4, code {
    margin: 0.5em;
}

h1 {
    font-size: 3.5em;
    margin-top: 0.5em;
}
h2 {
    font-size: 2em;
}
h3 {
    font-size: 1.5em;
}
h4 {
    font-size: 1em;
}
`)
	result := f.Build("text", true, &markdown, scroll, &css)
	return result
}
func (f *forge) ScrollKeybinds() *One {
	js := `
(function(){
  const { panel, frameIndex, state } = frameAPI.getContext();
  const content = panel.firstElementChild;
  const key = 'scroll_' + frameIndex;
  
  // Restore scroll position
  content.scrollTop = state[key] || 0;
  
  // Save scroll position
  content.addEventListener('scroll', () => {
    frameAPI.setState(key, content.scrollTop);
  });
  
  // Smooth scrolling
  let speed = 0;
  const scroll = () => {
    if (speed === 0) return;
    content.scrollBy({ top: speed });
    requestAnimationFrame(scroll);
  };
  
  frameAPI.onKey((k) => {
    const speeds = { w: -20, s: 20, a: -40, d: 40 };
    if (speeds[k]) {
      if (speed === 0) scroll();
      speed = speeds[k];
    }
  });
  
  document.addEventListener('keyup', (e) => {
    if ('wsad'.includes(e.key)) speed = 0;
  });
})();
`
	result := One(template.HTML(fmt.Sprintf(`<script>%s</script>`, js)))
	return &result
}

func (f *forge) BuildSlides(dir string) *One {
	prefix := f.AddPath(dir)
	img := f.Img("", "", "large")
	js := f.JS(fmt.Sprintf(`
(function() {
    const { panel, state } = frameAPI.getContext();
    const key = 'slideIndex';
    
    let slides = [];
    let index = state[key] || 0;

    async function show(i) {
        index = ((i %% slides.length) + slides.length) %% slides.length;
        frameAPI.setState(key, index);

        const img = panel.querySelector('.slides img');
        if (!img) return;

        const url = apiUrl + '/%s/' + slides[index];
        img.src = await frameAPI.fetch(url, url);
        img.alt = slides[index];
    }

    // Load slides and show first
    frameAPI.fetch(apiUrl + '/%s/slides', apiUrl + '/%s/slides')
        .then(data => {
            slides = data;
            show(index);
        });

    // Navigate slides
    frameAPI.onKey((k) => {
        if (k === 'a') show(index - 1);
        else if (k === 'd') show(index + 1);
    });
})();
    `, prefix, prefix, prefix))
	css := f.CSS(`
.slides {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
    width: 100%;
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
