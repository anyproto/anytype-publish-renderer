package utils

import (
	"context"
	"strings"

	"github.com/a-h/templ"
)

func TemplToString(component templ.Component) (string, error) {
	var sb strings.Builder
	err := component.Render(context.Background(), &sb)
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}
