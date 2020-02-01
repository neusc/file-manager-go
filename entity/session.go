package entity

import "time"

const SessionLength = 10 * 60

const SessionCleanInterval = 10

var SessionCleaned time.Time

type Session struct {
	Id           string    `json:"id"`
	UserName     string    `json:"username"`
	LastActivity time.Time `json:"lastactivity"`
}

func init() {
	SessionCleaned = time.Now()
}

func SetSessionCleaned(t time.Time) {
	SessionCleaned = t
}
