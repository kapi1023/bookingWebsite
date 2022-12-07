package renders

import (
	"bytes"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/kapi1023/bookingWebsite/packages/config"
	"github.com/kapi1023/bookingWebsite/packages/models"
)

var functions = template.FuncMap{}
var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a

}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {

	var tc map[string]*template.Template
	//get the template cache from the app config
	if app.UseCache {
		tc = app.TemplateCashe
	} else {
		tc, _ = CreateTemplateCashe()
	}

	//get requested templaet from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("could not get template from template cache")
	}
	buf := new(bytes.Buffer)

	_ = t.Execute(buf, td)

	//render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}

}

func CreateTemplateCashe() (map[string]*template.Template, error) {
	myCashe := map[string]*template.Template{}
	//get all of the filest name *.page.html
	pages, err := filepath.Glob("./templates/*page.html")
	if err != nil {
		return myCashe, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCashe, err
		}
		matches, err := filepath.Glob(("./templates/*.layout.html"))
		if err != nil {
			return myCashe, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return myCashe, err
			}
		}
		myCashe[name] = ts
	}
	return myCashe, nil
}
