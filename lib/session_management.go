package lib

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/op/go-logging"
	"io"
	"net/http"
	"strings"
	"time"
)

// Contains every user and the user's authentication Cookie.
var cookieStorage map[string]http.Cookie

// Init the cookieStorage-map.
func InitSessionManagement() {
	logging.MustGetLogger("").Debug("Initializing Session-Management...")
	cookieStorage = make(map[string]http.Cookie)
}

// Logs the user in by returning a Cookie containing a randomId.
func LoginAndGetCookie(username string) http.Cookie {
	// sessionValue: username + randomString
	randomId := getRandomId()
	sessionValue := strings.TrimSpace(username + ":" + randomId)

	// Build Cookie
	cookie := http.Cookie{Name: "session", Value: sessionValue, Path: "/", Expires: time.Now().AddDate(0, 0, 14), HttpOnly: true}

	// Save and return Cookie
	cookieStorage[username] = cookie
	return cookie
}

// Checks if the given Request contains a Cookie and if the Cookie authenticates a specific user.
func IsLoggedIn(r *http.Request) bool {
	// Get Cookie from Request
	rCookie, err := r.Cookie("session")
	if err != nil {
		return false
	}

	// Get data from received Cookie
	rCookieData := strings.Split(rCookie.Value, ":")

	// Get data from saved Cookie
	if _, ok := cookieStorage[rCookieData[0]]; !ok {
		return false
	}
	sCookie := cookieStorage[rCookieData[0]]
	sCookieData := strings.Split(sCookie.Value, ":")

	// Do not allow expired Cookies in Storage
	if sCookie.Expires.Before(time.Now()) {
		delete(cookieStorage, rCookieData[0])
		return false
	}

	// Check if the saved Cookie's randomId equals the received Cookie's randomId
	return rCookieData[1] == sCookieData[1]
}

// Logs the user out by returning a Cookie, which expired a day ago.
func LogoutAndDestroyCookie(r *http.Request) http.Cookie {
	cookie, _ := r.Cookie("session")

	// Remove the saved Cookie
	delete(cookieStorage, strings.Split(cookie.Value, ":")[0])
	logging.MustGetLogger("").Info("Logout successful.")

	// Return useless Cookie
	return http.Cookie{Name: "session", Value: "", Path: "/", Expires: time.Now().AddDate(0, 0, -1), HttpOnly: true}
}

// Returns a random, base64-encoded, string.
func getRandomId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
