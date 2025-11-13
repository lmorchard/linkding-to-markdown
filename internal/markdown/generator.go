package markdown

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/lmorchard/linkding-to-markdown/internal/linkding"
	"github.com/lmorchard/linkding-to-markdown/internal/templates"
)

// Generator handles markdown generation from bookmarks
type Generator struct {
	template *template.Template
}

// Options for generating markdown
type Options struct {
	Title          string
	IncludeNotes   bool
	IncludeTags    bool
	GroupByDate    bool
	DateFormat     string
}

// DefaultOptions returns the default generation options
func DefaultOptions() Options {
	return Options{
		Title:          "Bookmarks",
		IncludeNotes:   true,
		IncludeTags:    true,
		GroupByDate:    true,
		DateFormat:     "2006-01-02",
	}
}

// NewGenerator creates a new markdown generator with the default template
func NewGenerator() (*Generator, error) {
	defaultTmpl, err := templates.GetDefaultTemplate()
	if err != nil {
		return nil, fmt.Errorf("failed to get default template: %w", err)
	}
	return NewGeneratorWithTemplate(defaultTmpl)
}

// NewGeneratorFromFile creates a new markdown generator from a template file
func NewGeneratorFromFile(templatePath string) (*Generator, error) {
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}
	return NewGeneratorWithTemplate(string(content))
}

// NewGeneratorWithTemplate creates a new markdown generator with a custom template
func NewGeneratorWithTemplate(tmplStr string) (*Generator, error) {
	funcMap := template.FuncMap{
		"formatDate": func(t time.Time, format string) string {
			return t.Format(format)
		},
		"join": strings.Join,
		"hasContent": func(s string) bool {
			return strings.TrimSpace(s) != ""
		},
	}

	tmpl, err := template.New("markdown").Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return &Generator{
		template: tmpl,
	}, nil
}

// templateData holds data for template rendering
type templateData struct {
	Title      string
	Generated  string
	Bookmarks  []linkding.Bookmark
	GroupedBookmarks map[string][]linkding.Bookmark
	Options    Options
}

// Generate generates markdown from bookmarks and writes it to the writer
func (g *Generator) Generate(w io.Writer, bookmarks []linkding.Bookmark, opts Options) error {
	data := templateData{
		Title:     opts.Title,
		Generated: time.Now().Format(time.RFC3339),
		Bookmarks: bookmarks,
		Options:   opts,
	}

	// Group by date if requested
	if opts.GroupByDate {
		data.GroupedBookmarks = make(map[string][]linkding.Bookmark)
		for _, bookmark := range bookmarks {
			dateKey := bookmark.DateAdded.Format(opts.DateFormat)
			data.GroupedBookmarks[dateKey] = append(data.GroupedBookmarks[dateKey], bookmark)
		}
	}

	if err := g.template.Execute(w, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}
