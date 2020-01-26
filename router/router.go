package router

import (
	"log"
	"net/http"
	"time"

	"../constants"
	fm "./file"
	"./auth"
	_ "../config"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func InitRouter() {
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://www.shechuan.me"},
		AllowMethods:     []string{"PUT", "PATCH", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))
	user := r.Group("user")
	{
		user.POST("/signup", auth.SignUp)
		user.POST("/signin", auth.SignIn)
		user.POST("/logout", auth.LogOut)
		user.POST("/getUserInfo", auth.GetUserInfo)
	}
	file := r.Group("file")
	{
		file.POST("/upload", fm.UploadFile)
		file.POST("/list", fm.GetFileList)
		file.POST("/delete", fm.DeleteFile)
	}
	r.StaticFS("/", http.Dir(constants.StaticPath))
	r.Run(constants.StaticPort)
	log.Printf("listening on %s...", constants.StaticPort)
}
