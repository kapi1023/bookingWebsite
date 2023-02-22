package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/kapi1023/bookingWebsite/internal/models"
)

// AppConfig holds the app config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChannel   chan models.MailData
}
