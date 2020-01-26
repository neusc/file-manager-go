package common

import (
	"../config"
	"../entity"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

func CheckLogin(c *gin.Context) bool {
	cookie, err := c.Cookie("uid")
	if err != nil {
		return false
	}
	if !bson.IsObjectIdHex(cookie) {
		c.SetCookie(config.Cookie["name"], "", -1, "/", config.Cookie["domain"], false, false)
		return false
	}
	user := entity.User{}
	err = config.Session.DB("filemanager").C("users").FindId(bson.ObjectIdHex(cookie)).One(&user)
	if err != nil {
		return false
	}
	return true
}
