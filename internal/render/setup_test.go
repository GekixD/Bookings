package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/GekixD/Bookings/internal/config"
	"github.com/GekixD/Bookings/internal/models"
	"github.com/alexedwards/scs/v2"
)

var session *scs.SessionManager // session manager to simulate "real" http requests
var testApp config.AppConfig    // we need to create our own app for tests, since the render package already uses it

func TestMain(m *testing.M) {
	// What we want our session to contain:
	gob.Register(models.Reservation{}) //What do I want to store in the session
	testApp.Prod = false               // whether the web app is in producton or development

	session = scs.New()
	session.Lifetime = 24 * time.Hour              // set the lifetime for the session to 24 hours
	session.Cookie.Persist = true                  // whether the session will persis if they close the window
	session.Cookie.SameSite = http.SameSiteLaxMode // how strict is the cookie enforcement in the site
	session.Cookie.Secure = false                  // whether the cookies are encrupted (http vs https)

	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}

// myWriter is a dummy struct for a ResponseWriter
type myWriter struct{}

// Header() method for implementing the ResponseWritere interface for testing
func (mw *myWriter) Header() http.Header {
	var h http.Header
	return h
}

// WriteHeader() method for implementing the ResponseWritere interface for testing
func (mw *myWriter) WriteHeader(i int) {}

// Writer() method for implementing the ResponseWritere interface for testing, returns len([]byte)
func (mw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
