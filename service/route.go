package service

import (
	"bufio"
	"fmt"
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
	group := engine.Group("/v1")

	group.Use()

	group.Static("/transfer", "./transfer")

	group.Any("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "%s", "hello world")
	})

	//上传转换，并返回id
	group.POST("/upload", UploadPost())

	//从url下载视频
	group.POST("/download", func(ctx *gin.Context) {
		//filePath := ctx.Query("URL")
	})

	//下载并转换
	group.POST("/downloadTransform", func(ctx *gin.Context) {

	})

	//从服务器下载视频
	group.GET("/list/:id", func(ctx *gin.Context) {

	})

	//查看状态
	group.GET("/status/:id", StatusGet())

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

	if i == 0 {
		return "", fmt.Errorf("upload with %d", i)
	}
	log.Println("success:", i)
	return fileName, nil
}

// UploadPost ...
func UploadPost() gin.HandlerFunc {
	// swagger:route GET /upload/{base64(url)} upload a stream
	//
	//     Responses:
	//       200: UserResponse
	return func(ctx *gin.Context) {

		defer ctx.Request.Body.Close()
		src := "./upload/"
		fileName, err := writeTo(src, ctx.Request.Body)
		if err != nil {
			ctx.JSON(http.StatusOK, JSON(-1, err.Error()))
			return
		}
		b, err := openssl.HexKey()
		if err != nil {
			ctx.JSON(http.StatusOK, JSON(-1, err.Error()))
			return
		}

		stream := NewStreamer(string(b), fileName)
		stream.SetURI("http://localhost:8080/infos" + "/" + fileName + "/key")
		stream.SetDst("./transfer/")
		stream.SetSrc(src)
		queue.Push(stream)
		client.Set(fileName, string(b), 0)
		log.Println(fileName)
		ctx.JSON(http.StatusOK, JSON(0, "ok", gin.H{"id": fileName}))
		return
	}
}

// StatusGet 获取视频转换状态
func StatusGet() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		key, err := client.Get(id).Result()
		if err != nil {
			ctx.JSON(http.StatusOK, JSON(-1, err.Error()))
			return
		}

		if key != "" {
			ctx.JSON(http.StatusOK, JSON(1, "processing"))
			return
		}

		ctx.JSON(http.StatusOK, JSON(0, "ok", gin.H{}))
		return
	}
}
