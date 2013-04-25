package main

import (
	"encoding/base64"
	"net/http"
	"time"
)

var (
	// the name of the flash cookie
	flashName = "flash"

	// how long for the flash to last
	flashDuration = 12 * time.Hour
)

// prepare a cookie with the content of the flash
func makeFlashCookie(note string) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = flashName
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(note))
	cookie.Expires = time.Now().Add(flashDuration)
	cookie.Path = "/"
	return cookie
}

// prepare a cookie which has no content and has already expired
func makeDeleteFlash() *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = flashName
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(""))
	cookie.Expires = time.Unix(0, 0)
	cookie.Path = "/"
	return cookie
}

// write a flash into a cookie (eww, but it works)
func setFlash(rw http.ResponseWriter, note string) {
	http.SetCookie(rw, makeFlashCookie(note))
}

// delete a flash from a cookie
func deleteFlash(rw http.ResponseWriter) {
	http.SetCookie(rw, makeDeleteFlash())
}

// extract the flash from a request
func getFlash(req *http.Request) string {
	cookie, err := req.Cookie(flashName)
	if err != nil {
		// no cookie
	} else {
		if cookie.Value != "" {
			value, _ := base64.URLEncoding.DecodeString(cookie.Value)
			return string(value)
		}
	}

	// there was no flash
	return ""
}
