package frame

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
    // Initialize state if not present
    if (!panel.state.slideIndex) panel.state.slideIndex = 0;
    if (!panel.state.slides) panel.state.slides = [];
    
    // Only fetch if we don't have slides cached
    if (panel.state.slides.length === 0) {
        fetch(apiUrl + '/slides/slides')
            .then(response => response.json())
            .then(data => {
                panel.state.slides = data;
                if (panel.state.slides.length > 0) showSlide(panel.state.slideIndex);
            })
            .catch(error => console.error('Error loading slides:', error));
    } else {
        // Restore to saved slide
        showSlide(panel.state.slideIndex);
    }

    function showSlide(index) {
        if (panel.state.slides.length === 0) return;
        panel.state.slideIndex = ((index % panel.state.slides.length) + panel.state.slides.length) % panel.state.slides.length;
        const img = panel.element.querySelector('.slides img');
        if (img) {
            img.src = apiUrl + '/slides/' + panel.state.slides[panel.state.slideIndex];
            img.alt = panel.state.slides[panel.state.slideIndex];
        }
    }
    
    panel.element.addEventListener('panelKey', (e) => {
        if (e.detail.key === 'a') showSlide(panel.state.slideIndex - 1);
        else if (e.detail.key === 'd') showSlide(panel.state.slideIndex + 1);
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
