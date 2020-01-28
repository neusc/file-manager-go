package auth

import (
	"../../common"
	"../../config"
	"../../entity"
	
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

type SignUpForm struct {
	Name       string `json:"name" validate:"min=5,max=20"`
	Password   string `json:"password" validate:"min=6,max=12"`
	Repassword string `json:"repassword" validate:"eqfield=Password"`
}

func SignUp(c *gin.Context) {
	if _, ok := common.CheckLogin(c); ok {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 2,
			"msg":        "redirect",
			"data":       "list",
		})
		return
	}
	var form SignUpForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("%+v\n", form)
	err := validateForm(form)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 1,
			"msg":        err.Error(),
		})
		return
	}

	user := entity.User{}
	err = config.Session.DB("filemanager").C("users").Find(bson.M{"name": form.Name}).One(&user)
	// user existed
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 1,
			"msg":        "username has existed! please change another one!",
		})
		return
	}

	originPwd := []byte(form.Password)
	hashPwd, _ := bcrypt.GenerateFromPassword(originPwd, bcrypt.DefaultCost)

	user.Id = bson.NewObjectId()
	user.Name = form.Name
	user.Password = string(hashPwd)
	config.Session.DB("filemanager").C("users").Insert(user)
	// c.SetCookie(config.Conf.CookieName, user.Id.Hex(), 3600, "/", config.Conf.CookieDomain, false, false)
	c.JSON(http.StatusCreated, gin.H{
		"statusCode": 2,
		"msg":        "success",
		"data":       "signin",
	})
}

var validate *validator.Validate

func validateForm(form SignUpForm) error {
	validate = validator.New()
	err := validate.Struct(form)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return err
		}

		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err.Field())
			fmt.Println(err.Value())
			fmt.Print(err.Tag())
		}
		return err
	}
	return nil
}

func SignIn(c *gin.Context) {
	if _, ok := common.CheckLogin(c); ok {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 2,
			"msg":        "redirect",
			"data":       "list",
		})
		return
	}
	var form entity.User
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user entity.User
	err := config.Session.DB("filemanager").C("users").Find(bson.M{"name": form.Name}).One(&user)
	if err == nil {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": 1,
				"msg":        "password is not correct!",
			})
			return
		}
		c.SetCookie(config.Conf.CookieName, user.Id.Hex(), 3600, "/", config.Conf.CookieDomain, false, false)
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 2,
			"msg":        "redirect",
			"data":       "list",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"statusCode": 1,
		"msg":        "user don't exist! please go to sign up!",
	})
}

func LogOut(c *gin.Context) {
	cookie, err := c.Cookie("uid")
	if err != nil {
		c.String(http.StatusOK, "current user is not login!")
		return
	}
	user := entity.User{}
	err = config.Session.DB("filemanager").C("users").FindId(bson.ObjectIdHex(cookie)).One(&user)
	if err != nil {
		c.String(http.StatusOK, "current cookie info error!")
		return
	}
	c.SetCookie(config.Conf.CookieName, "", -1, "/", config.Conf.CookieDomain, false, false)
	c.JSON(http.StatusOK, gin.H{
		"statusCode": 2,
		"msg":        "redirect",
		"data":       "signin",
	})

}

type UserInfo struct {
	Uid  string `json:"uid"`
	Name string `json:"name"`
}

func GetUserInfo(c *gin.Context) {
	var params UserInfo
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !bson.IsObjectIdHex(params.Uid) {
		c.SetCookie(config.Conf.CookieName, "", -1, "/", config.Conf.CookieDomain, false, false)
		return
	}
	user := entity.User{}
	err := config.Session.DB("filemanager").C("users").FindId(bson.ObjectIdHex(params.Uid)).One(&user)
	if err != nil {
		c.SetCookie(config.Conf.CookieName, "", -1, "/", config.Conf.CookieDomain, false, false)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	params.Name = user.Name
	c.JSON(http.StatusOK, gin.H{
		"statusCode": 0,
		"msg":        "success",
		"data":       params,
	})
}
