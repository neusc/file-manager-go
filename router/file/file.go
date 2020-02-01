package file

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"../../common"
	"../../entity"
	"../../config"
	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	user, ok := common.GetLoginInfo(c);
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 2,
			"msg":        "redirect",
			"data":       "signin",
		})
		return
	}
	c.Request.ParseMultipartForm(32 << 20)
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	dir := config.Conf.StaticPath + user.Id.Hex()
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, os.ModePerm)
	}
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("create dir err: %s", err.Error()))
		return
	}
	files := form.File["files"]
	for _, file := range files {
		if err := c.SaveUploadedFile(file, dir+"/"+file.Filename); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"statusCode": 0,
		"msg":        fmt.Sprintf("Uploaded successfully %d files .", len(files)),
	})
}

func GetFileList(c *gin.Context) {
	user, ok := common.GetLoginInfo(c);
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 2,
			"msg":        "redirect",
			"data":       "signin",
		})
		return
	}
	fileInfo, err := ioutil.ReadDir(config.Conf.StaticPath + user.Id.Hex())
	var fileList []entity.File
	if err != nil {
		fmt.Println("Read Dir error", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 0,
			"msg":        "success",
			"data":       fileList,
		})
		return
	}
	for _, file := range fileInfo {
		fileItem := entity.File{Name: file.Name(), Path: config.Conf.FilePath + user.Id.Hex() + "/" + file.Name(), Size: file.Size(), ModTime: file.ModTime().Unix()}
		fileList = append(fileList, fileItem)
	}
	c.JSON(http.StatusOK, gin.H{
		"statusCode": 0,
		"msg":        "success",
		"data":       fileList,
	})
	// response := entity.ResponseData{StatusCode: 0, Msg: "success", Data: fileList}
	// json.NewEncoder(w).Encode(response)
}

type deleteParams struct {
	Names []string `json:"names"`
}

func DeleteFile(c *gin.Context) {
	user, ok := common.GetLoginInfo(c);
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 2,
			"msg":        "redirect",
			"data":       "signin",
		})
		return
	}
	var params deleteParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, name := range params.Names {
		deleteErr := os.Remove(config.Conf.StaticPath + user.Id.Hex() + "/" + name)
		if deleteErr != nil {
			fmt.Println("delete file err", deleteErr.Error())
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"statusCode": 0,
		"msg":        "success",
		"data":       nil,
	})
	// response := entity.ResponseData{StatusCode: 0, Msg: "success", Data: nil}
	// json.NewEncoder(w).Encode(response)
}
