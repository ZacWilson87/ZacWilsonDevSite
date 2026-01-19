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
	"year":     func() int { return time.Now().Year() },
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
	r.Get("/projects/{slug}", s.handleProjectDetail)
	r.Get("/case-study/{slug}", s.handleCaseStudyDetail)
	r.Get("/about", s.handleAbout)
	r.Get("/contact", s.handleContact)

	// Legacy redirects
	r.Get("/work", redirectTo("/projects"))
	r.Get("/work/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		http.Redirect(w, r, "/projects/"+slug, http.StatusMovedPermanently)
	})
	r.Get("/case-studies", redirectTo("/projects"))
	r.Get("/case-studies/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		http.Redirect(w, r, "/case-study/"+slug, http.StatusMovedPermanently)
	})

	return r
}

func redirectTo(target string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target, http.StatusMovedPermanently)
	}
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
	data := map[string]any{
		"Projects":    s.loader.Projects(),
		"CaseStudies": s.loader.CaseStudies(),
	}
	s.render(w, "projects.html", data)
}

func (s *Server) handleProjectDetail(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	project := s.loader.Project(slug)
	if project == nil {
		http.NotFound(w, r)
		return
	}
	s.render(w, "project.html", project)
}

func (s *Server) handleCaseStudyDetail(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	study := s.loader.CaseStudy(slug)
	if study == nil {
		http.NotFound(w, r)
		return
	}
	s.render(w, "case-study.html", study)
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
