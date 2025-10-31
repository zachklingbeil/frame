package frame

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
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
	markdown := One(template.HTML(buf.String()))
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
  const panel = frameAPI.getPanel(document.currentScript);
  const content = panel.firstElementChild;
  const state = frameAPI.getState(panel);
  
  // Restore scroll after images load
  frameAPI.waitForImages(content).then(() => {
    if (state.scroll !== undefined) {
      content.scrollTop = state.scroll;
    }
  });
  
  // Save scroll
  content.addEventListener('scroll', () => {
    frameAPI.setState(panel, 'scroll', content.scrollTop);
  });
  
  // Handle scrolling
  let scrolling = 0;
  const step = () => {
    if (!scrolling) return;
    content.scrollBy({ top: scrolling });
    requestAnimationFrame(step);
  };
  
  frameAPI.onKey(panel, (key) => {
    if (key === 'w') scrolling = -25;
    else if (key === 's') scrolling = 25;
    else if (key === 'a') scrolling = -50;
    else if (key === 'd') scrolling = 50;
    else return;
    step();
  });
  
  document.addEventListener('keyup', (e) => {
    if (['w','s','a','d'].includes(e.key)) scrolling = 0;
  });
})();
`
	result := One(template.HTML(fmt.Sprintf(`<script>%s</script>`, js)))
	return &result
}

func (f *forge) BuildSlides(dir string) *One {
	f.AddPath(dir)
	img := f.Img("", "", "large")
	js := f.JS(`
(function() {
    const panel = frameAPI.getPanel(document.currentScript);
    const state = frameAPI.getState(panel);
    
    if (state.slideIndex === undefined) {
        frameAPI.setState(panel, 'slideIndex', 0);
    }
    
    function showSlide(index, slides) {
        if (!slides?.length) return;
        const newIndex = ((index % slides.length) + slides.length) % slides.length;
        frameAPI.setState(panel, 'slideIndex', newIndex);
        
        const img = panel.querySelector('.slides img');
        if (img) {
            img.src = apiUrl + '/slides/' + slides[newIndex];
            img.alt = slides[newIndex];
        }
    }
    
    window.app.frameSource.fetchResource('slides', apiUrl + '/slides/slides')
        .then(slides => {
            const state = frameAPI.getState(panel);
            showSlide(state.slideIndex, slides);
        });
    
    frameAPI.onKey(panel, (key) => {
        const slides = window.app.frameSource.getCachedResource('slides');
        if (!slides) return;
        
        const state = frameAPI.getState(panel);
        if (key === 'a') showSlide(state.slideIndex - 1, slides);
        else if (key === 'd') showSlide(state.slideIndex + 1, slides);
    });
})();
    `)
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
