package server

import (
	"html/template"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"devsite/internal/content"
)

// Server handles HTTP requests
type Server struct {
	loader    *content.Loader
	baseTempl *template.Template
}

// Template functions
var funcMap = template.FuncMap{
	"add":      func(a, b int) int { return a + b },
	"safeHTML": func(s string) template.HTML { return template.HTML(s) },
}

// New creates a new server instance
func New(loader *content.Loader) *Server {
	// Parse base templates (layout and partials) that will be cloned for each page
	base := template.New("base").Funcs(funcMap)
	template.Must(base.ParseGlob("templates/partials/*.html"))
	template.Must(base.ParseFiles("templates/layout.html"))

	return &Server{
		loader:    loader,
		baseTempl: base,
	}
}

// Router returns the configured chi router
func (s *Server) Router() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Routes
	r.Get("/", s.handleHome)
	r.Get("/projects", s.handleProjects)
	r.Get("/projects/{slug}", s.handleProject)
	r.Get("/about", s.handleAbout)
	r.Get("/contact", s.handleContact)

	// Redirects for old URLs
	r.Get("/work", redirectTo("/projects"))
	r.Get("/work/{slug}", s.redirectProjectSlug)
	r.Get("/case-studies", redirectTo("/projects"))
	r.Get("/case-studies/{slug}", s.redirectProjectSlug)

	return r
}

func redirectTo(target string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target, http.StatusMovedPermanently)
	}
}

func (s *Server) redirectProjectSlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	http.Redirect(w, r, "/projects/"+slug, http.StatusMovedPermanently)
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	projects := s.loader.Projects()

	// Limit to 3 featured projects
	featured := projects
	if len(featured) > 3 {
		featured = featured[:3]
	}

	data := map[string]any{
		"FeaturedProjects": featured,
		"Year":             time.Now().Year(),
	}
	s.render(w, "home.html", data)
}

func (s *Server) handleProjects(w http.ResponseWriter, r *http.Request) {
	workItems := s.loader.ProjectsByType("work")
	caseStudies := s.loader.ProjectsByType("case-study")

	data := map[string]any{
		"WorkItems":   workItems,
		"CaseStudies": caseStudies,
	}
	s.render(w, "projects.html", data)
}

func (s *Server) handleProject(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	project := s.loader.Project(slug)
	if project == nil {
		http.NotFound(w, r)
		return
	}
	s.render(w, "project.html", project)
}

func (s *Server) handleAbout(w http.ResponseWriter, r *http.Request) {
	s.render(w, "about.html", nil)
}

func (s *Server) handleContact(w http.ResponseWriter, r *http.Request) {
	s.render(w, "contact.html", nil)
}

func (s *Server) render(w http.ResponseWriter, name string, data any) {
	// Clone base templates and parse the specific page template
	tmpl, err := s.baseTempl.Clone()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Re-apply funcs to ensure they're available for newly parsed templates
	tmpl = tmpl.Funcs(funcMap)

	_, err = tmpl.ParseFiles("templates/" + name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
