package config

import (
	"log"
	"text/template"

	"github.com/alexedwards/scs/v2"
)

// AppConfig holds the app config
type AppConfig struct {
	UseCache      bool
	TemplateCashe map[string]*template.Template
	infoLog       *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
}
