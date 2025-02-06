package utils

import (
	"context"
	"strings"
	"unicode"

	"github.com/a-h/templ"
)

var Colors = map[string]string{
	"default": "#252525",
	"grey":    "#b6b6b6",
	"yellow":  "#ecd91b",
	"orange":  "#ffb522",
	"red":     "#f55522",
	"pink":    "#e51ca0",
	"purple":  "#ab50cc",
	"blue":    "#3e58eb",
	"ice":     "#2aa7ee",
	"teal":    "#0fc8ba",
	"lime":    "#5dd400",
}

func TemplToString(component templ.Component) (string, error) {
	var sb strings.Builder
	err := component.Render(context.Background(), &sb)
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}

func GetColor(name string) string {
	if hex, exists := Colors[name]; exists {
		return hex
	}
	return Colors["default"]
}

// Capitalize the first letter of a string
func Capitalize(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}