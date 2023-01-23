package render

import (
	"net/http"
	"testing"

	"github.com/kapi1023/bookingWebsite/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	session.Put(r.Context(), "flash", 123)
	result := AddDefaultData(&td, r)

	if result.Flash == "123" {
		t.Error("flash value 123 not found")
	}

}

func TestRenderTempolate(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
	app.TemplateCache = tc
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	var ww myWriter

	err = Template(&ww, "home.page.html", &models.TemplateData{}, r)
	if err != nil {
		t.Error("error writing template to browser")
	}
	err = Template(&ww, "non-existing-template.page.html", &models.TemplateData{}, r)
	if err == nil {
		t.Error("rendered template that should not exist")
	}

}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}

func TestNewTemplate(t *testing.T) {
	NewRenderer(app)

}

func TestCreateTemplateCashe(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

}
