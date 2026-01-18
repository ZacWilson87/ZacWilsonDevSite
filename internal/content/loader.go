package content

import (
	"bytes"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
)

// Markdown renderer
var md = goldmark.New()

// frontmatter represents the YAML frontmatter in content files
type frontmatter struct {
	Title  string   `yaml:"title"`
	Date   string   `yaml:"date"`
	Status string   `yaml:"status"`
	Tags   []string `yaml:"tags"`
}

// Item represents a piece of content (project or case study)
type Item struct {
	Slug    string
	Title   string
	Date    time.Time
	Status  string
	Tags    []string
	Body    string // rendered HTML
	RawBody string // original markdown
}

// Loader loads and caches content from markdown files
type Loader struct {
	basePath    string
	mu          sync.RWMutex
	projects    map[string]*Item
	caseStudies map[string]*Item
}

// NewLoader creates a new content loader
func NewLoader(basePath string) *Loader {
	return &Loader{
		basePath:    basePath,
		projects:    make(map[string]*Item),
		caseStudies: make(map[string]*Item),
	}
}

// Load reads all content from disk and caches it in memory
func (l *Loader) Load() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Load projects
	if err := l.loadDir("projects", l.projects); err != nil {
		return err
	}

	// Load case studies
	if err := l.loadDir("case-studies", l.caseStudies); err != nil {
		return err
	}

	return nil
}

func (l *Loader) loadDir(dir string, dest map[string]*Item) error {
	path := filepath.Join(l.basePath, dir)

	entries, err := os.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // directory doesn't exist yet, that's ok
		}
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		item, err := l.parseFile(filepath.Join(path, entry.Name()))
		if err != nil {
			return err
		}

		dest[item.Slug] = item
	}

	return nil
}

func (l *Loader) parseFile(path string) (*Item, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	slug := filepath.Base(path)
	slug = slug[:len(slug)-len(filepath.Ext(slug))]

	// Parse frontmatter
	fm, body := parseFrontmatter(data)

	// Render markdown to HTML
	var htmlBuf bytes.Buffer
	if err := md.Convert(body, &htmlBuf); err != nil {
		return nil, err
	}

	item := &Item{
		Slug:    slug,
		Title:   fm.Title,
		Status:  fm.Status,
		Tags:    fm.Tags,
		RawBody: string(body),
		Body:    htmlBuf.String(),
	}

	// Use slug as title fallback
	if item.Title == "" {
		item.Title = slug
	}

	// Parse date if present
	if fm.Date != "" {
		if t, err := time.Parse("2006-01-02", fm.Date); err == nil {
			item.Date = t
		}
	}

	return item, nil
}

// parseFrontmatter extracts YAML frontmatter from content
// Returns the parsed frontmatter and the remaining body
func parseFrontmatter(data []byte) (frontmatter, []byte) {
	var fm frontmatter

	// Check for frontmatter delimiter
	if !bytes.HasPrefix(data, []byte("---\n")) {
		return fm, data
	}

	// Find closing delimiter
	rest := data[4:] // skip opening "---\n"
	end := bytes.Index(rest, []byte("\n---"))
	if end == -1 {
		return fm, data
	}

	// Parse YAML
	yaml.Unmarshal(rest[:end], &fm)

	// Return body after closing delimiter
	body := rest[end+4:] // skip "\n---"
	if len(body) > 0 && body[0] == '\n' {
		body = body[1:]
	}

	return fm, body
}

// Projects returns all loaded projects
func (l *Loader) Projects() []*Item {
	l.mu.RLock()
	defer l.mu.RUnlock()

	items := make([]*Item, 0, len(l.projects))
	for _, item := range l.projects {
		items = append(items, item)
	}
	return items
}

// Project returns a single project by slug
func (l *Loader) Project(slug string) *Item {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.projects[slug]
}

// CaseStudies returns all loaded case studies
func (l *Loader) CaseStudies() []*Item {
	l.mu.RLock()
	defer l.mu.RUnlock()

	items := make([]*Item, 0, len(l.caseStudies))
	for _, item := range l.caseStudies {
		items = append(items, item)
	}
	return items
}

// CaseStudy returns a single case study by slug
func (l *Loader) CaseStudy(slug string) *Item {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.caseStudies[slug]
}
