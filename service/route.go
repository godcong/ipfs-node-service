package service

import (
	"bufio"
	"github.com/gin-gonic/gin"
	"github.com/godcong/go-ffmpeg/util"
	"io"
	"log"
	"net/http"
	"os"
)

// Router ...
func Router(engine *gin.Engine) error {
	group := engine.Group("/")

	group.Use()

	group.POST("/uploadTransform", func(context *gin.Context) {
		defer context.Request.Body.Close()

		fileName, err := writeTo("./tmp", context.Request.Body)
		if err != nil {
			context.JSON(http.StatusOK, JSON(-1, err.Error()))
			return
		}
		ToM3U8("./tmp/" + fileName)
		context.JSON(http.StatusOK, JSON(0, "ok", gin.H{"name": fileName}))
		return
	})

	group.GET("/download", func(context *gin.Context) {
		//filePath := context.Query("path")

	})

	group.POST("/downloadTransform", func(context *gin.Context) {

	})

	return nil
}

func writeTo(path string, reader io.Reader) (string, error) {
	fileName := util.GenerateRandomString(64)
	file, err := os.OpenFile(path+"/"+fileName, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	i, err := io.Copy(writer, reader)
	if err != nil {
		return "", err
	}
	log.Println("success:", i)
	return fileName, nil
}
