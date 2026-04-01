package web

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin/render"
	"go.rtnl.ai/x/rlog"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

//go:embed all:static
//go:embed all:templates
var content embed.FS

var (
	// These files are included with all matched page templates
	includes = []string{"base.html", "components/*.html"}

	// Any file matching one of these patterns will be a page with a named template and
	// all of the includes files will be included in the page template.
	pages = []string{"errors/*.html", "errors/*/*.html", "pages/*.html", "pages/*/*.html"}

	// Any file matching one of these patterns will be a partial template and will be
	// available to render as a template but with no includes or other template files.
	partials = []string{"partials/*.html", "partials/*/*.html"}
)

// Static creates the StaticFS that contains our static assets to be served.
func Static() http.FileSystem {
	staticFiles, err := fs.Sub(content, "static")
	if err != nil {
		panic(fmt.Errorf("failed to create static file system: %w", err))
	}
	return http.FS(staticFiles)
}

// Templates returns the FileSystem that contains our HTML templates to be rendered.
func Templates() fs.FS {
	templatesFiles, err := fs.Sub(content, "templates")
	if err != nil {
		panic(fmt.Errorf("failed to create templates file system: %w", err))
	}
	return templatesFiles
}

// Creates a new template renderer from the default templates.
func NewRender(fsys fs.FS) (render *Render, err error) {
	render = &Render{
		templates: make(map[string]*template.Template),
	}

	// Add the pages patterns
	for _, pattern := range pages {
		if !globExists(fsys, pattern) {
			continue
		}

		if err = render.AddPattern(fsys, pattern, includes...); err != nil {
			return nil, err
		}
	}

	// Add the partials patterns
	for _, pattern := range partials {
		if !globExists(fsys, pattern) {
			continue
		}

		if err = render.AddPattern(fsys, pattern); err != nil {
			return nil, err
		}
	}

	return render, nil
}

// Implements the render.HTMLRender interface from the gin framework.
type Render struct {
	templates map[string]*template.Template
	funcs     template.FuncMap
}

var _ render.HTMLRender = (*Render)(nil)

func (r *Render) Instance(name string, data any) render.Render {
	return &render.HTML{
		Template: r.templates[name],
		Name:     filepath.Base(name),
		Data:     data,
	}
}

func (r *Render) AddPattern(fsys fs.FS, pattern string, includes ...string) (err error) {
	var names []string
	if names, err = fs.Glob(fsys, pattern); err != nil {
		return err
	}

	for _, name := range names {
		patterns := make([]string, 0, len(includes)+1)
		patterns = append(patterns, includes...)
		patterns = append(patterns, name)

		tmpl := template.New(name).Funcs(r.FuncMap())
		if r.templates[name], err = tmpl.ParseFS(fsys, patterns...); err != nil {
			return err
		}

		rlog.TraceAttrs(
			context.Background(),
			"template parsed and added to renderer",
			slog.String("name", name),
			slog.Any("patterns", patterns),
		)
	}
	return nil
}

// Create functions available to all templates.
func (r *Render) FuncMap() template.FuncMap {
	if r.funcs == nil {
		r.funcs = template.FuncMap{
			"static":      static,
			"titlecase":   titlecase,
			"lowercase":   lowercase,
			"uppercase":   uppercase,
			"datetime":    datetime,
			"currentYear": currentYear,
		}
	}
	return r.funcs
}

// ===========================================================================
// Template Functions
// ===========================================================================

func titlecase(s string) string {
	return cases.Title(language.English).String(s)
}

func lowercase(s string) string {
	return strings.ToLower(s)
}

func uppercase(s string) string {
	return strings.ToUpper(s)
}

func datetime(t time.Time) string {
	return t.Format("January 2, 2006 3:04 PM")
}

func static(path string) string {
	return filepath.Join("/static", path)
}

func currentYear() int {
	return time.Now().Year()
}

// ===========================================================================
// Helper Functions
// ===========================================================================

func globExists(fsys fs.FS, pattern string) (exists bool) {
	names, _ := fs.Glob(fsys, pattern)
	return len(names) > 0
}
