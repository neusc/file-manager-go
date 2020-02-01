package auth

import (
	"../../common"
	"../../config"
	"../../entity"

	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"time"
)

type SignUpForm struct {
	Name       string `json:"name" validate:"min=3,max=8"`
	Password   string `json:"password" validate:"min=6"`
	Repassword string `json:"repassword" validate:"eqfield=Password"`
}

func SignUp(c *gin.Context) {
	if common.CheckLogin(c) {
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

	var user entity.User
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
			fmt.Print(err.Tag() + "\n")
		}
		return err
	}
	return nil
}

func SignIn(c *gin.Context) {
	if common.CheckLogin(c) {
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
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 1,
			"msg":        "user don't exist! please go to sign up!",
		})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 1,
			"msg":        "password is not correct!",
		})
		return
	}
	// create session
	sID, _ := uuid.NewV4()
	c.SetCookie(config.Conf.CookieName, sID.String(), entity.SessionLength, "/", config.Conf.CookieDomain, false, false)
	session := entity.Session{
		Id:           sID.String(),
		UserName:     form.Name,
		LastActivity: time.Now(),
	}
	// store session
	config.Session.DB("filemanager").C("sessions").Insert(session)

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 2,
		"msg":        "redirect",
		"data":       "list",
	})
}

func LogOut(c *gin.Context) {
	if !common.CheckLogin(c) {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 2,
			"msg":        "redirect",
			"data":       "signin",
		})
		return
	}
	cookie, _ := c.Cookie("filemanager")
	config.Session.DB("filemanager").C("sessions").Remove(bson.M{"id": cookie})
	c.SetCookie(config.Conf.CookieName, "", -1, "/", config.Conf.CookieDomain, false, false)
	// clean dbSessions
	if time.Now().Sub(entity.SessionCleaned) > (time.Second * entity.SessionCleanInterval) {
		go common.CleanSessions()
	}
	c.JSON(http.StatusOK, gin.H{
		"statusCode": 2,
		"msg":        "redirect",
		"data":       "signin",
	})
}

type userInfo struct {
	SessionID  string `json:"sessionid"`
	Name       string `json:"name"`
}

// GetUserInfo return userInfo
func GetUserInfo(c *gin.Context) {
	if !common.CheckLogin(c) {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 2,
			"msg":        "redirect",
			"data":       "signin",
		})
		return
	}
	user := common.GetUser(c)
	cookie, _ := c.Cookie("filemanager")
	params := userInfo{
		SessionID: cookie,
		Name: user.Name,
	}
	c.JSON(http.StatusOK, gin.H{
		"statusCode": 0,
		"msg":        "success",
		"data":       params,
	})
}
