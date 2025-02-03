package renderer

import "github.com/a-h/templ"

type BlockParams struct {
	Id          string
	Classes     []string
	Content     templ.Component
	Additional  templ.Component
	ChildrenIds []string
}
