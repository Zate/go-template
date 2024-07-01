package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

// TODO restructure all this to support text/template for prompts and use FS for templates
// see prompt_templates.go for more info

// TemplateRegistry is a struct to hold the templates
type TemplateRegistry struct {
	templates *template.Template
}

// Render is a function to render a template
func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
