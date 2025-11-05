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
  const { panel, frameIndex, state } = frameAPI.context();
  const content = panel.firstElementChild;
  const key = 'scroll_' + frameIndex;
  
  // Restore scroll position
  content.scrollTop = state[key] || 0;
  
  // Save scroll position
  content.addEventListener('scroll', () => {
    frameAPI.update(key, content.scrollTop);
  });
  
  // Smooth scrolling
  let speed = 0;
  let isScrolling = false;
  
  const scroll = () => {
    if (speed === 0) {
      isScrolling = false;
      return;
    }
    content.scrollBy({ top: speed });
    requestAnimationFrame(scroll);
  };
  
  const speeds = { w: -20, s: 20, a: -40, d: 40 };
  
  frameAPI.onKey((k) => {
    if (speeds[k]) {
      speed = speeds[k];
      if (!isScrolling) {
        isScrolling = true;
        scroll();
      }
    }
  });
  
  // Stop scrolling on key release (global listener)
  document.addEventListener('keyup', (e) => {
    if (speeds[e.key]) {
      const current = frameAPI.context();
      if (current.panel === panel && current.frameIndex === frameIndex) {
        speed = 0;
      }
    }
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
    const { panel, frameIndex, state } = frameAPI.context();
    const key = 'slideIndex_' + frameIndex;
    
    let slides = [];
    let index = state[key] || 0;

    async function show(i) {
        if (slides.length === 0) return;
        
        index = ((i %% slides.length) + slides.length) %% slides.length;
        frameAPI.update(key, index);

        const img = panel.querySelector('.slides img');
        const slideName = slides[index];
        const url = apiUrl + '/%s/' + slideName;
        const { data } = await frameAPI.fetch(slideName, url);
        img.src = data;
        img.alt = slideName;
    }

    // Load slides list
    frameAPI.fetch('slides-%s', apiUrl + '/%s/slides')
        .then(({ data }) => {
            slides = data;
            if (slides.length > 0) show(index);
        })
        .catch(err => console.error('Failed to load slides:', err));

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
    width: 100%%;
    height: 100%%;
    box-sizing: border-box;
    overflow: hidden;
}
.slides img {
    max-width: 100%%;
    max-height: 100%%;
    object-fit: contain;
}
    `)
	return f.Build("slides", true, img, &css, &js)
}
