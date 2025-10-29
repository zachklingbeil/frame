package frame

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	math "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	h "github.com/yuin/goldmark/renderer/html"
)

type Text interface {
	AddMarkdown(file string) *One
	H1(s string) *One
	H2(s string) *One
	H3(s string) *One
	H4(s string) *One
	H5(s string) *One
	H6(s string) *One
	Paragraph(s string) *One
	Span(s string) *One
	Strong(s string) *One
	Em(s string) *One
	Small(s string) *One
	Mark(s string) *One
	Del(s string) *One
	Ins(s string) *One
	Sub(s string) *One
	Sup(s string) *One
	Kbd(s string) *One
	Samp(s string) *One
	VarElem(s string) *One
	Abbr(s string) *One
	Time(s string) *One
	Button(label string) *One
	Code(code string) *One
	ScrollKeybinds() *One
}

// --- text Implementation ---
type text struct {
	Md *goldmark.Markdown
}

func NewText() Text {
	return &text{
		Md: initGoldmark(),
	}
}

func initGoldmark() *goldmark.Markdown {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, math.MathJax),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
			parser.WithBlockParsers(),
			parser.WithInlineParsers(),
		),
		goldmark.WithRendererOptions(
			h.WithHardWraps(),
			h.WithXHTML(),
		),
	)
	return &md
}

func (t *text) AddMarkdown(file string) *One {
	content, err := os.ReadFile(file)
	if err != nil {
		empty := One("")
		return &empty
	}

	var buf bytes.Buffer
	if err := (*t.Md).Convert(content, &buf); err != nil {
		empty := One("")
		return &empty
	}

	result := One(template.HTML(buf.String()))
	return &result
}

func (t *text) ScrollKeybinds() *One {
	js := `
(function(panel){
  const content = panel.firstElementChild;
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
})(panel);
`
	result := One(template.HTML(fmt.Sprintf(`<script>%s</script>`, js)))
	return &result
}

func (t *text) H1(s string) *One         { return Tag("h1", s) }
func (t *text) H2(s string) *One         { return Tag("h2", s) }
func (t *text) H3(s string) *One         { return Tag("h3", s) }
func (t *text) H4(s string) *One         { return Tag("h4", s) }
func (t *text) H5(s string) *One         { return Tag("h5", s) }
func (t *text) H6(s string) *One         { return Tag("h6", s) }
func (t *text) Paragraph(s string) *One  { return Tag("p", s) }
func (t *text) Span(s string) *One       { return Tag("span", s) }
func (t *text) Strong(s string) *One     { return Tag("strong", s) }
func (t *text) Em(s string) *One         { return Tag("em", s) }
func (t *text) Small(s string) *One      { return Tag("small", s) }
func (t *text) Mark(s string) *One       { return Tag("mark", s) }
func (t *text) Del(s string) *One        { return Tag("del", s) }
func (t *text) Ins(s string) *One        { return Tag("ins", s) }
func (t *text) Sub(s string) *One        { return Tag("sub", s) }
func (t *text) Sup(s string) *One        { return Tag("sup", s) }
func (t *text) Kbd(s string) *One        { return Tag("kbd", s) }
func (t *text) Samp(s string) *One       { return Tag("samp", s) }
func (t *text) VarElem(s string) *One    { return Tag("var", s) }
func (t *text) Abbr(s string) *One       { return Tag("abbr", s) }
func (t *text) Time(s string) *One       { return Tag("time", s) }
func (t *text) Button(label string) *One { return Tag("button", label) }
func (t *text) Code(code string) *One    { return Tag("code", code) }
