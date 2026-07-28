package renderer

import (
	"embed"
	"encoding/hex"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"maps"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/blake2s"
)

type Renderer struct {
	funcs              template.FuncMap
	templates          map[string]*template.Template
	defaultData        map[string]any
	staticHashCache    map[string]string
	webFS              embed.FS
	staticFilesRoot    string
	webStaticFilesRoot string
	templateSuffix     string
}

func New(webFS embed.FS, staticFilesRoot string, webStaticFilesRoot string, templateFilesRoot string, templateSuffix string, defaultData map[string]any) (*Renderer, error) {
	r := &Renderer{
		defaultData:        defaultData,
		templates:          make(map[string]*template.Template),
		staticHashCache:    make(map[string]string),
		webFS:              webFS,
		staticFilesRoot:    strings.TrimSuffix(staticFilesRoot, "/"),
		webStaticFilesRoot: strings.TrimSuffix(webStaticFilesRoot, "/"),
		templateSuffix:     templateSuffix,
	}

	r.funcs = template.FuncMap{
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
		"blake2sSum256":      hashData,
		"hashStatic":         r.hashStatic,
		"staticPathWithHash": r.staticPathWithHash,
	}

	templateFS, err := fs.Sub(webFS, templateFilesRoot)
	if err != nil {
		return nil, err
	}

	err = r.initTemplates(templateFS)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Renderer) staticPathWithHash(filePath string) string {
	return fmt.Sprintf("%s?v=%s", filePath, r.hashStatic(filePath))
}

func (r *Renderer) hashStatic(filePath string) string {
	filePath = strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(filePath, "/"), r.webStaticFilesRoot), "/")

	if hash, ok := r.staticHashCache[filePath]; ok {
		return hash
	}
	fullFilePath := filepath.Join(r.staticFilesRoot, filePath)
	data, err := r.webFS.ReadFile(fullFilePath)
	if err != nil {
		panic(err)
	}
	hash := hashData(data)
	r.staticHashCache[filePath] = hash
	return hash
}

func hashData(data []byte) string {
	sum := blake2s.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func (r *Renderer) initTemplates(templateFS fs.FS) error {
	templateSuffix := r.templateSuffix
	baseTemplateName := filepath.Join("base" + templateSuffix)
	partials, err := fs.Glob(templateFS, filepath.Join("partials", "*"+templateSuffix))
	if err != nil {
		return err
	}
	baseFiles := append([]string{baseTemplateName}, partials...)

	foundLayoutFiles, err := fs.Glob(templateFS, filepath.Join("layouts", "*"+templateSuffix))
	if err != nil {
		return err
	}

	slog.Debug("Parsing base template", "baseTemplateName", baseTemplateName, "baseFiles", baseFiles)
	baseTemplate, err := template.New(baseTemplateName).Funcs(r.funcs).ParseFS(templateFS, baseFiles...)
	if err != nil {
		return err
	}
	// https://stackoverflow.com/questions/50842389/parsing-multiple-templates-in-go
	for _, layoutFile := range foundLayoutFiles {
		layoutName := filepath.Base(layoutFile)
		slog.Debug("Parsing layout", "templateName", layoutName)

		templ, err := baseTemplate.Clone()
		if err != nil {
			return err
		}
		templ, err = templ.ParseFS(templateFS, layoutFile)
		if err != nil {
			return err
		}

		r.templates[layoutName] = templ
	}

	return nil
}

// Render a template without additional data.
// Data passed as defaultData to the Renderer constructor will always be used.
func (r *Renderer) RenderWithoutData(w http.ResponseWriter, templateName string) error {
	return r.Render(w, templateName, map[string]any{})
}

func (r *Renderer) Render(w http.ResponseWriter, templateName string, data map[string]any) error {
	combinedData := r.defaultData
	templateName = templateName + r.templateSuffix

	maps.Copy(combinedData, data)

	t, ok := r.templates[templateName]
	if !ok {
		return fmt.Errorf("Template '%s' doesn't exist!", templateName)
	}

	err := t.Execute(w, combinedData)
	if err != nil {
		return err
	}

	return nil
}

func (r *Renderer) handleError(w http.ResponseWriter, err error) {
	slog.Error("Renderer error", "err", err)
	http.Error(w, "Internal Server Error", 500)
}
