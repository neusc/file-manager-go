package common

import (
	"../config"
	"../entity"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/satori/go.uuid"
	"time"
)

func GetUser(c *gin.Context) entity.User {
	cookie, err := c.Cookie("filemanager")
	value := cookie
	if err != nil {
		sID, _ := uuid.NewV4()
		value = sID.String()
	}
	c.SetCookie(config.Conf.CookieName, value, entity.SessionLength, "/", config.Conf.CookieDomain, false, false)

	var user entity.User
	var session entity.Session
	err = config.Session.DB("filemanager").C("sessions").Find(bson.M{"id": value}).One(&session)
	if err != nil {
		return user
	}
	session.LastActivity = time.Now()
	config.Session.DB("filemanager").C("sessions").Update(bson.M{"id": value}, session)
	config.Session.DB("filemanager").C("users").Find(bson.M{"name": session.UserName}).One(&user)
	return user
}

func CheckLogin(c *gin.Context) bool {
	cookie, err := c.Cookie("filemanager")
	if err != nil {
		return false
	}
	var session entity.Session
	err = config.Session.DB("filemanager").C("sessions").Find(bson.M{"id": cookie}).One(&session)
	if err != nil {
		return false
	}
	session.LastActivity = time.Now()
	config.Session.DB("filemanager").C("sessions").Update(bson.M{"id": cookie}, session)

	var user entity.User
	err = config.Session.DB("filemanager").C("users").Find(bson.M{"name": session.UserName}).One(&user)
	if err != nil {
		return false
	}
	// refresh session
	c.SetCookie(config.Conf.CookieName, cookie, entity.SessionLength, "/", config.Conf.CookieDomain, false, false)
	return true
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
