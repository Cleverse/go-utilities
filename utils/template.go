package utils

import (
	"bytes"
	"fmt"
	htmltemplate "html/template"
	"io"
	texttemplate "text/template"
)

type TemplateConfig struct {
	Funcs  FuncMap
	Option []string
}

type FuncMap map[string]interface{}

// Template is a wrapper of text/template and html/template.
type Template struct {
	textTemplate   *texttemplate.Template
	htmlTemplate   *htmltemplate.Template
	config         TemplateConfig
	name           string
	templateString string
}

// NewTemplate is a wrapper of text/template and html/template.
// used to define a template at initialization time without error checking.
//
// Example:
//
//	package main
//
//	import "github.com/Cleverse/go-utilities/utils"
//
//	var HelloWorldTemplate = utils.NewTemplate("hello", "Hello {{.Name}}", utils.TemplateConfig{})
//
//	func main() { ... }
func NewTemplate(name string, tmpl string, config TemplateConfig) *Template {
	return &Template{
		name:           name,
		templateString: tmpl,
		config:         config,
	}
}

func (t *Template) ParsedTemplate() (*texttemplate.Template, error) {
	if t.textTemplate != nil {
		return t.textTemplate, nil
	}

	tmpl, err := texttemplate.
		New(t.name).
		Funcs(texttemplate.FuncMap(t.config.Funcs)).
		Option(t.config.Option...).
		Parse(t.templateString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse text template: %w", err)
	}

	t.textTemplate = tmpl
	return t.textTemplate, nil
}

func (t *Template) ParsedTemplateHTML() (*htmltemplate.Template, error) {
	if t.htmlTemplate != nil {
		return t.htmlTemplate, nil
	}

	tmpl, err := htmltemplate.
		New(t.name).
		Funcs(htmltemplate.FuncMap(t.config.Funcs)).
		Option(t.config.Option...).
		Parse(t.templateString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse html template: %w", err)
	}

	t.htmlTemplate = tmpl
	return t.htmlTemplate, nil
}

// String returns the template string.
func (t *Template) String() string {
	return t.templateString
}

// Text applies the template to the given args and returns the result as a string.
func (t *Template) Text(args interface{}) (string, error) {
	buf := new(bytes.Buffer)
	if err := t.Execute(buf, args); err != nil {
		return "", fmt.Errorf("failed to execute text template: %w", err)
	}
	return buf.String(), nil
}

// HTML applies a parsed template to the specified data object and returns the output as a string.
func (t *Template) HTML(args interface{}) (string, error) {
	buf := new(bytes.Buffer)
	if err := t.ExecuteHTML(buf, args); err != nil {
		return "", fmt.Errorf("failed to execute html template: %w", err)
	}
	return buf.String(), nil
}

// Execute applies a parsed template to the specified data object and writes the output to wr.
func (t *Template) Execute(wr io.Writer, args interface{}) error {
	tmpl, err := t.ParsedTemplate()
	if err != nil {
		return err
	}

	if err := tmpl.Execute(wr, args); err != nil {
		return fmt.Errorf("failed to execute text template: %w", err)
	}

	return err
}

// ExecuteHTML applies a parsed template to the specified data object and writes the output to wr.
func (t *Template) ExecuteHTML(wr io.Writer, args interface{}) error {
	tmpl, err := t.ParsedTemplateHTML()
	if err != nil {
		return err
	}

	if err := tmpl.Execute(wr, args); err != nil {
		return fmt.Errorf("failed to execute text template: %w", err)
	}

	return err
}
