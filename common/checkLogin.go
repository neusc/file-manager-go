package common

import (
	"../config"
	"../entity"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"time"
)

func GetLoginInfo(c *gin.Context) (entity.User, bool) {
	var user entity.User
	cookie, err := c.Cookie("filemanager")
	if err != nil {
		return user, false
	}
	var session entity.Session
	err = config.Session.DB("filemanager").C("sessions").Find(bson.M{"id": cookie}).One(&session)
	if err != nil {
		// reset valid cookie
		c.SetCookie(config.Conf.CookieName, "", -1, "/", config.Conf.CookieDomain, false, false)
		return user, false
	}
	session.LastActivity = time.Now()
	config.Session.DB("filemanager").C("sessions").Update(bson.M{"id": cookie}, session)
	err = config.Session.DB("filemanager").C("users").Find(bson.M{"name": session.UserName}).One(&user)
	if err != nil {
		return user, false
	}
	// refresh session
	c.SetCookie(config.Conf.CookieName, cookie, entity.SessionLength, "/", config.Conf.CookieDomain, false, false)
	return user, true
}

func CleanSessions() {
	var result entity.Session
	iter := config.Session.DB("filemanager").C("sessions").Find(nil).Iter()
	for iter.Next(&result) {
		if time.Now().Sub(result.LastActivity) > (time.Second * entity.SessionCleanInterval) {
			config.Session.DB("filemanager").C("sessions").Remove(bson.M{"id": result.Id})
		}
	}
	entity.SetSessionCleaned(time.Now())
}
