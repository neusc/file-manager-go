package file

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"../../common"
	"../../constants"
	"../../entity"
	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	c.Request.ParseMultipartForm(32 << 20)
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	files := form.File["files"]
	for _, file := range files {
		if err := c.SaveUploadedFile(file , constants.StaticPath+file.Filename); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
	}
	c.String(http.StatusOK, fmt.Sprintf("Uploaded successfully %d files .", len(files)))
}

func GetFileList(c *gin.Context) {
	if !common.CheckLogin(c) {
		c.JSON(http.StatusOK, gin.H{
			"statusCode": 2,
			"msg":        "redirect",
			"data":       "signin",
		})
		return
	}
	fileInfo, err := ioutil.ReadDir(constants.StaticPath)
	if err != nil {
		fmt.Println("Read Dir error", err)
	}
	var fileList []entity.File
	for _, file := range fileInfo {
		fileItem := entity.File{Name: file.Name(), Path: constants.FilePath + file.Name(), Size: file.Size(), ModTime: file.ModTime().Unix()}
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
	var params deleteParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, name := range params.Names {
		deleteErr := os.Remove(constants.StaticPath + name)
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
