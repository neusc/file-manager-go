package common

import (
	"../config"
	"../entity"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

func CheckLogin(c *gin.Context) (entity.User, bool) {
	cookie, err := c.Cookie("uid")
	user := entity.User{}
	if err != nil {
		return user, false
	}
	if !bson.IsObjectIdHex(cookie) {
		c.SetCookie(config.Conf.CookieName, "", -1, "/", config.Conf.CookieDomain, false, false)
		return user, false
	}
	err = config.Session.DB("filemanager").C("users").FindId(bson.ObjectIdHex(cookie)).One(&user)
	if err != nil {
		return user, false
	}
	return user, true
}
