package controllers

import (
	"net/http"
	"time"
)

type Session struct {
	username string
	expiry   time.Time
}

// Track user sessions in memory
var Sessions = map[string]*Session{}

func (s *Session) SetTime(time time.Time) {
	s.expiry = time
}

func (s *Session) GetUsername() string {
	return s.username
}

func Expired(s *Session) bool {
	return s.expiry.Before(time.Now())
}

func GetSession(r *http.Request) (string, *Session) {
	sessionToken := GetSessionToken(r)
	if sessionToken == "" {
		return "", &Session{}
	}
	return sessionToken, Sessions[sessionToken]
}

func GetSessionToken(r *http.Request) string {

	c, err := r.Cookie("session_token")
	if err != nil {
		return ""
	}
	sessionToken := c.Value

	return sessionToken
}
