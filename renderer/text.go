package renderer

import "github.com/a-h/templ"

func (r *Renderer) RenderText(text string) templ.Component {
	return TextTemplate(r, text)
}
