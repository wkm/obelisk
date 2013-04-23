package main

import (
	"encoding/base64"
	"net/http"
	"time"
)

var (
	flashName     = "flash"
	flashDuration = 12 * time.Hour
)

func makeFlashCookie(note string) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = flashName
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(note))
	cookie.Expires = time.Now().Add(flashDuration)
	cookie.Path = "/"
	return cookie
}

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
