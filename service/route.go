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

	//上传转换，并返回id
	group.POST("/uploadTransform", func(context *gin.Context) {
		defer context.Request.Body.Close()

		fileName, err := writeTo("./tmp", context.Request.Body)
		if err != nil {
			context.JSON(http.StatusOK, JSON(-1, err.Error()))
			return
		}
		go ToM3U8("./tmp/" + fileName)
		context.JSON(http.StatusOK, JSON(0, "ok", gin.H{"name": fileName}))
		return
	})

	//从url下载视频
	group.POST("/download", func(context *gin.Context) {
		//filePath := context.Query("URL")
	})

	//下载并转换
	group.POST("/downloadTransform", func(context *gin.Context) {

	})

	//从服务器下载视频
	group.GET("/get/:id", func(context *gin.Context) {

	})

	//查看状态
	group.GET("/status/:id", func(context *gin.Context) {
		id := context.Param("id")
		context.String(http.StatusOK, "id: %s", id)
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
	err = writer.Flush()
	if err != nil {
		return "", err
	}
	log.Println("success:", i)
	return fileName, nil
}
