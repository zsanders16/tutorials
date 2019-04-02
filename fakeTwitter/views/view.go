package views

import (
	"net/http"
	"path/filepath"
	"text/template"
)

var (
	// LayoutDir is the directory where all laylout are stored
	LayoutDir = "views/layouts/"
	// TemplateDir is the directory where all views are located
	TemplateDir = "views/"
	// TemplateExt is the extension of the templates
	TemplateExt = ".gohtml"
)

// NewView returns a View struct
func NewView(layout string, files ...string) *View {
	layoutFiles, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	addTemplateAndExtPath(files)

	files = append(files, layoutFiles...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

// View holds all layout templates
type View struct {
	Template *template.Template
	Layout   string
}

// Render renders the templates with the View to the provided ResponseWriter
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := v.Render(w, nil); err != nil {
		panic(err)
	}
}

func addTemplateAndExtPath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f + TemplateExt
	}
}
