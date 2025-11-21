package templates

import (
	"html/template"

	"github.com/timefactoryio/frame/zero"
)

func (t *templates) Keyboard() {
	css := t.CSS(t.KeyboardCSS())
	js := t.JS(`
(function(){
  const { panel } = pathless.context();
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
    keyEl.textContent = k.toUpperCase();
    if (entry && entry.style) keyEl.style.cssText = entry.style;
    grid.appendChild(keyEl);
  });

  const updateKey = (k, pressed) => {
    const keyEl = grid.querySelector('[data-key="' + k + '"]');
    if (keyEl) keyEl.classList.toggle('pressed', pressed);
  };

  document.addEventListener('keydown', (e) => {
    if (keyMap.has(e.key)) updateKey(e.key, true);
  });

  document.addEventListener('keyup', (e) => {
    if (keyMap.has(e.key)) updateKey(e.key, false);
  });
})();
`)
	html := zero.One(template.HTML(`<div class="grid"></div>`))
	t.Build("keyboard", true, &html, &css, &js)
}
