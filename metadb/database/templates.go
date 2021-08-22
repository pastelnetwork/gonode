package database

import (
	"bytes"
	"io/ioutil"
	"path"
	"strings"
	"text/template"

	"github.com/pastelnetwork/gonode/common/errors"
)

// TemplateKeeper is a query template keeper
type TemplateKeeper struct {
	mp map[string]*template.Template
}

func substituteTemplate(tmpl *template.Template, data interface{}) (string, error) {
	var templateBuffer bytes.Buffer
	if tmpl == nil || data == nil {
		return "", errors.Errorf("input nil template or data")
	}
	if err := tmpl.Execute(&templateBuffer, data); err != nil {
		return "", err
	}
	return templateBuffer.String(), nil
}

// GetCommand get a sql command
func (k *TemplateKeeper) GetCommand(key string, data interface{}) (string, error) {
	template := k.GetTemplate(key)
	if template == nil {
		return "", errors.Errorf("no template for key: %s", key)
	}
	return substituteTemplate(template, data)
}

// GetTemplate get a sql template
func (k *TemplateKeeper) GetTemplate(key string) *template.Template {
	if tmpl, found := k.mp[key]; found {
		return tmpl
	}
	return nil
}

// NewTemplateKeeper create new template keeper instance
func NewTemplateKeeper(templateDir string) (*TemplateKeeper, error) {
	files, err := ioutil.ReadDir(templateDir)
	if err != nil {
		return nil, err
	}

	mp := make(map[string]*template.Template)
	for _, file := range files {
		filename := file.Name()
		if path.Ext(filename) != ".tmpl" {
			continue
		}
		filepath := path.Join(templateDir, filename)
		parts := strings.Split(filename, ".")
		key := parts[0]

		tmpl, err := template.ParseFiles(filepath)
		if err != nil {
			return nil, err
		}
		mp[key] = tmpl
	}

	return &TemplateKeeper{
		mp: mp,
	}, nil
}
