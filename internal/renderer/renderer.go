package renderer

import (
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type Renderer struct {
	funcs       template.FuncMap
	templates   map[string]*template.Template
	defaultData map[string]any
}

var rendererFuncs = template.FuncMap{
	"formatDate": func(t time.Time) string {
		return t.Format("2006-01-02")
	},
	"formatDateTime": func(t time.Time) string {
		return t.Format("2006-01-02T15:04")
	},
	"formatDateTimeLocal": func(t time.Time) string {
		return t.Format("2006-01-02 15:04")
	},
	"add": func(a, b int) int {
		return a + b
	},
	"sub": func(a, b int) int {
		return a - b
	},
	"seq": func(start, end int) []int {
		n := end - start + 1
		if n <= 0 {
			return nil
		}
		s := make([]int, n)
		for i := range s {
			s[i] = start + i
		}
		return s
	},
	"join": strings.Join,
	"isAfter": func(checkAfter, base time.Time) bool {
		return checkAfter.After(base)
	},
}

func New(templateFS fs.FS, templateSuffix string, defaultData map[string]any) (*Renderer, error) {
	r := &Renderer{
		funcs:       rendererFuncs,
		defaultData: defaultData,
		templates:   make(map[string]*template.Template),
	}

	err := r.initTemplates(templateFS, templateSuffix)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Renderer) initTemplates(templateFS fs.FS, templateSuffix string) error {
	baseTemplate := filepath.Join("base" + templateSuffix)
	partials, err := fs.Glob(templateFS, filepath.Join("partials", "*"+templateSuffix))
	if err != nil {
		return err
	}
	baseFiles := append([]string{baseTemplate}, partials...)

	foundLayoutFiles, err := fs.Glob(templateFS, filepath.Join("layouts", "*"+templateSuffix))
	if err != nil {
		return err
	}

	for _, layoutFile := range foundLayoutFiles {
		layoutName := strings.TrimSuffix(filepath.Base(layoutFile), templateSuffix)
		layoutFiles := append(baseFiles, layoutFile)

		slog.Debug("Parsing layout", "layoutFiles", layoutFiles)

		templ := template.New(layoutName).Funcs(r.funcs)
		templ, err := templ.ParseFS(templateFS, layoutFiles...)
		if err != nil {
			return err
		}

		r.templates[layoutName] = templ
	}

	return nil
}

// Render a template without additional data.
// Data passed as defaultData to the Renderer constructor will always be used.
func (r *Renderer) RenderWithoutData(w http.ResponseWriter, templateName string) {
	r.Render(w, templateName, map[string]any{})
}

func (r *Renderer) Render(w http.ResponseWriter, templateName string, data map[string]any) {
	combinedData := r.defaultData

	for k, v := range data {
		combinedData[k] = v
	}

	t, ok := r.templates[templateName]
	if !ok {
		r.handleError(w, fmt.Errorf("Template '%s' doesn't exist!", templateName))
	}

	err := t.Execute(w, data)
	if err != nil {
		r.handleError(w, err)
	}
}

func (r *Renderer) handleError(w http.ResponseWriter, err error) {
	slog.Error("Renderer error", "err", err)
	http.Error(w, "Internal Server Error", 500)
}
