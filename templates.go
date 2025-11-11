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

func (f *forge) Landing(src, heading, githubLogo, githubLink, twitterLogo, twitterLink string) {
	img := f.Img(src, "logo", "large")
	h1 := f.H1(heading)

	footer := One(template.HTML(fmt.Sprintf(`
        <div class="footer-icons">
            <a href="%s" target="_blank" rel="noopener">
                <img src="%s" alt="GitHub" class="icon" />
            </a>
            <a href="%s" target="_blank" rel="noopener">
                <img src="%s" alt="Twitter" class="icon" />
            </a>
        </div>
    `, githubLink, githubLogo, twitterLink, twitterLogo)))

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
        .footer-icons {
            display: flex;
            justify-content: center;
            gap: 1.5em;
            margin-top: 1.5em;
        }
        .footer-icons img.icon {
            width: 2.5em;
            height: 2.5em;
            object-fit: contain;
        }
    `)
	f.Build("zero", true, &css, img, h1, &footer)
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
  const { frame, state } = pathless.context();
  const key = 'scroll_';
  
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
    const key = 'slideIndex';
    
    let slides = [];
    let index = state[key] || 0;

    async function show(i) {
        if (!slides.length) return;
        
        index = ((i %% slides.length) + slides.length) %% slides.length;
        pathless.update(key, index);

        const img = frame.querySelector('img');
        if (!img) return;
        
        const slide = slides[index];
        const { data } = await pathless.fetch(apiUrl + '/%s/' + slide, { key: slide });
        img.src = data;
        img.alt = slide;
    }

    pathless.fetch(apiUrl + '/%s/order', { key: 'order-%s' })
          .then(({ data }) => {
              slides = data;
              if (slides.length) show(index);
          });

    pathless.onKey((k) => {
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
