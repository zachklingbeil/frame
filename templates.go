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

func (f *forge) YouTube(videoID string) *One {

	embedURL := fmt.Sprintf("https://www.youtube.com/embed/%s", template.HTMLEscapeString(videoID))
	embed := One(template.HTML(
		fmt.Sprintf(`<embed src="%s" class="youtube-panel" type="text/html" />`, embedURL),
	))
	return f.Build("youtube", true, &embed)
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
    scroll-behavior: smooth;
    height: 100%;
}
.text img {
    max-width: 90%;
    max-height: 90%;
    object-fit: contain;
    align-items: center;
    justify-content: center;
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
  const panel = document.currentScript.closest('.panel');
  const state = panel.__frameState;
  const content = panel.firstElementChild;
  
  // Restore scroll position from state
  if (state.scrollTop !== undefined) {
    content.scrollTop = state.scrollTop;
  }
  
  // Save scroll position to state
  const saveScroll = () => {
    state.scrollTop = content.scrollTop;
  };
  content.addEventListener('scroll', saveScroll);
  
  let scrolling = 0;
  const step = () => {
    if (!scrolling) return;
    content.scrollBy({ top: scrolling });
    requestAnimationFrame(step);
  };
  
  const handleScroll = (key) => {
    if (key === 'w') scrolling = -25;
    else if (key === 's') scrolling = 25;
    else if (key === 'a') scrolling = -50;
    else if (key === 'd') scrolling = 50;
    else return false;
    step();
    return true;
  };
  
  panel.addEventListener('panelKey', (e) => {
    handleScroll(e.detail.key);
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
    const panel = document.currentScript.closest('.panel');
    const state = panel.__frameState;
    const frameSource = panel.__frameSource;
    
    // Initialize panel-specific state
    if (state.slideIndex === undefined) {
        state.slideIndex = 0;
    }
    
    function showSlide(index, slides) {
        if (!slides || slides.length === 0) return;
        
        // Update panel-specific state
        state.slideIndex = ((index % slides.length) + slides.length) % slides.length;
        
        const img = panel.querySelector('.slides img');
        if (img) {
            img.src = apiUrl + '/slides/' + slides[state.slideIndex];
            img.alt = slides[state.slideIndex];
        }
    }
    
    // Load slides using FrameSource cache (shared across all panels)
    frameSource.fetchResource('slides', apiUrl + '/slides/slides')
        .then(slides => {
            showSlide(state.slideIndex, slides);
        })
        .catch(error => {
            console.error('Error loading slides:', error);
        });
    
    // Handle navigation
    panel.addEventListener('panelKey', (e) => {
        const slides = frameSource.getCachedResource('slides');
        if (!slides) return;
        
        if (e.detail.key === 'a') showSlide(state.slideIndex - 1, slides);
        else if (e.detail.key === 'd') showSlide(state.slideIndex + 1, slides);
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
