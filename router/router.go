package router

import (
	"log"
	"net/http"
	"time"

	"../config"
	"./auth"
	fm "./file"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
		MaxAge:           12 * time.Hour,
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
	r.StaticFS("/", http.Dir(config.Conf.StaticPath))
	r.Run(config.Conf.StaticPort)
	log.Printf("listening on %s...", config.Conf.StaticPort)
}
