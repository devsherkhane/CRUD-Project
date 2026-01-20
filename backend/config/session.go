package config

import "github.com/gorilla/sessions"

var Store = sessions.NewCookieStore([]byte("very-secret-key"))

func InitSession() {
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600, // 1 hour
		HttpOnly: true,
	}
}
