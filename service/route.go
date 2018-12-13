package service

import (
	"bufio"
	"github.com/gin-gonic/gin"
	"github.com/godcong/go-ffmpeg/openssl"
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

	group.Static("/stream", "./split")

	group.Any("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "%s", "hello world")
	})

	//上传转换，并返回id
	group.POST("/uploadTransform", func(ctx *gin.Context) {
		defer ctx.Request.Body.Close()
		src := "./upload/"
		fileName, err := writeTo(src, ctx.Request.Body)
		if err != nil {
			ctx.JSON(http.StatusOK, JSON(-1, err.Error()))
			return
		}
		b, err := openssl.Base64Key()
		if err != nil {
			ctx.JSON(http.StatusOK, JSON(-1, err.Error()))
			return
		}

		stream := NewStreamer(string(b), fileName)
		stream.SetURI("localhost:8080/stream")
		stream.SetDst("./transfer/")
		stream.SetSrc(src)

		queue.Push(stream)
		log.Println(fileName)
		ctx.JSON(http.StatusOK, JSON(0, "ok", gin.H{"name": fileName}))
		return
	})

	//从url下载视频
	group.POST("/download", func(ctx *gin.Context) {
		//filePath := ctx.Query("URL")
	})

	//下载并转换
	group.POST("/downloadTransform", func(ctx *gin.Context) {

	})

	//从服务器下载视频
	group.GET("/get/:id", func(ctx *gin.Context) {

	})

	//查看状态
	group.GET("/status/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		ctx.String(http.StatusOK, "id: %s", id)
	})

	return nil
}

func writeTo(path string, reader io.Reader) (string, error) {
	fileName := util.GenerateRandomString(64)
	file, err := os.OpenFile(path+fileName, os.O_CREATE|os.O_RDWR, os.ModePerm)
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
