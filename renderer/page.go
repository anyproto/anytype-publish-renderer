package renderer

func (r Renderer) RenderPage() (err error) {
	templ := PageTemplate(&r)
	err = r.templToString(templ)
	return
}
