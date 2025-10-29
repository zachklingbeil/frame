package frame

import (
	"fmt"
	"html/template"
)

func (f *frame) Zero(src, heading string) {
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

func (t *text) ScrollKeybinds() *One {
	js := `
(function(panel){
  const content = panel.firstElementChild;
  let scrolling = 0;
  
  // Restore scroll position
  panel.addEventListener('restoreState', (e) => {
    if (e.detail.scrollTop !== undefined) {
      content.scrollTop = e.detail.scrollTop;
    }
  });
  
  // Save scroll position (debounced)
  let saveTimeout;
  content.addEventListener('scroll', () => {
    clearTimeout(saveTimeout);
    saveTimeout = setTimeout(() => {
      panel.dataset.customState = JSON.stringify({
        scrollTop: content.scrollTop
      });
    }, 100);
  });
  
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
})(panel);
`
	result := One(template.HTML(fmt.Sprintf(`<script>%s</script>`, js)))
	return &result
}

func (f *frame) BuildText(file string) *One {
	text := f.AddMarkdown(file)
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

p {
	font-size: 1.2em;
	line-height: 1.5;
	margin: 1em ;
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
	final := f.Build("text", true, text, scroll, &css)
	return final
}

func (f *frame) BuildSlides(dir string) *One {
	f.AddPath(dir, "slides")
	img := f.Img("", "", "large")
	js := f.JS(`
(function(panel){
    panel.slides = [];
    
    // Restore slide index
    panel.addEventListener('restoreState', (e) => {
        if (e.detail.slideIndex !== undefined && panel.slides.length > 0) {
            showSlide(e.detail.slideIndex);
        }
    });
    
    fetch(apiUrl + '/slides/slides')
        .then(response => response.json())
        .then(data => {
            panel.slides = data;
            if (panel.slides.length > 0) {
                // Try to restore saved state, otherwise start at 0
                const savedState = panel.dataset.customState;
                const startIndex = savedState ? 
                    JSON.parse(savedState).slideIndex : 0;
                showSlide(startIndex);
            }
        })
        .catch(error => console.error('Error loading slides:', error));

    function showSlide(index) {
        if (panel.slides.length === 0) return;
        const slideIndex = ((index % panel.slides.length) + panel.slides.length) % panel.slides.length;
        const img = panel.querySelector('.slides img');
        if (img) {
            img.src = apiUrl + '/slides/' + panel.slides[slideIndex];
            img.alt = panel.slides[slideIndex];
        }
        // Save current slide index
        panel.dataset.customState = JSON.stringify({
            slideIndex: slideIndex
        });
    }
    
    panel.addEventListener('panelKey', (e) => {
        const currentState = panel.dataset.customState;
        const currentIndex = currentState ? 
            JSON.parse(currentState).slideIndex : 0;
        
        if (e.detail.key === 'a') showSlide(currentIndex - 1);
        else if (e.detail.key === 'd') showSlide(currentIndex + 1);
    });
})(panel);
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
