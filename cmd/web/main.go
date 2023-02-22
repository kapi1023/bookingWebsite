package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/kapi1023/bookingWebsite/internal/config"
	"github.com/kapi1023/bookingWebsite/internal/driver"
	"github.com/kapi1023/bookingWebsite/internal/handlers"
	"github.com/kapi1023/bookingWebsite/internal/helpers"
	"github.com/kapi1023/bookingWebsite/internal/models"
	render "github.com/kapi1023/bookingWebsite/internal/renders"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main function
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	defer close(app.MailChannel)

	listenForMail()

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

func run() (*driver.DB, error) {
	//what is store in session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Reservation{})

	mailChan := make(chan models.MailData)
	app.MailChannel = mailChan
	// change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog
	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	//connect to database
	log.Println("connecting to database")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookingWebsite user=postgres password=zaq12#WSX")
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	log.Println("succesfully connected to database")
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Println("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))

	return db, nil
}
