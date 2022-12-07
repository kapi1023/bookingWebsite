package handlers

import (
	"net/http"

	"github.com/kapi1023/bookingWebsite/packages/config"
	"github.com/kapi1023/bookingWebsite/packages/models"
	"github.com/kapi1023/bookingWebsite/packages/renders"
)

//Repo is repository used by the handlers

var Repo *Repository

// Repository is repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the about page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	renders.RenderTemplate(w, "home.page.html", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	//perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	//send the data
	renders.RenderTemplate(w, "about.page.html", &models.TemplateData{
		StringMap: stringMap})
}
