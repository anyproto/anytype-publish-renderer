package renderer

import "github.com/a-h/templ"

func (r *Renderer) RenderPage() templ.Component {
	return PageTemplate(r)
}
