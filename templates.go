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
	img := f.Img("", "", "")
	js := f.JS(`
(function(panel){
    let slideIndex = 0;
    let slides = [];
    fetch(apiUrl + '/slides/slides.json')
        .then(response => response.json())
        .then(data => {
            slides = data;
            if (slides.length > 0) showSlide(0);
        })
        .catch(error => console.error('Error loading slides:', error));

    function showSlide(index) {
        if (slides.length === 0) return;
        slideIndex = ((index % slides.length) + slides.length) % slides.length;
        const img = panel.querySelector('.slides img');
        if (img) {
            img.src = apiUrl + '/slides/' + slides[slideIndex];
            img.alt = slides[slideIndex];
        }
    }
    panel.addEventListener('panelKey', (e) => {
        if (e.detail.key === 'a') showSlide(slideIndex - 1);
        else if (e.detail.key === 'd') showSlide(slideIndex + 1);
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
