package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/GekixD/Bookings/internal/config"
)

var app *config.AppConfig

// NEwheplers sets up app config for helpers
func NewHelepers(a *config.AppConfig) {
	app = a
}

func ClientError(res http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of", status)
	http.Error(res, http.StatusText(status), status)
}

func ServerError(res http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack()) // This allows us to trace back the error stack
	app.ErrorLog.Println(trace)
	http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
