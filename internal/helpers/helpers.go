package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/kapi1023/bookingWebsite/internal/config"
)

var app *config.AppConfig

func NewHelpers(a *config.AppConfig) {
	app = a
}
func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status: ", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Print("Server error with trace: ", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

}
