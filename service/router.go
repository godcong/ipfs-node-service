package service

import (
	"github.com/gin-gonic/gin"
	"github.com/rakyll/statik/fs"
	"log"
	"net/http"
)

// Router ...
func Router(engine *gin.Engine) error {
	//api document
	//engine.Static("/doc", "./doc")
	st, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	engine.StaticFS("/doc", st)
	engine.Static("/transfer", "./transfer")
	engine.Static("/upload", "./upload")

	ver := "v1"
	group := engine.Group(ver)
	group.Use(func(ctx *gin.Context) {
		log.Println("visited:", ctx.Request.URL.String())
	})

	group.Any("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "%s", "hello world")
	})

	//上传文件
	group.POST("/upload", UploadPost(ver))
	//获取文件
	group.POST("/rd", RemoteDownloadPost(ver))
	//视频转换
	group.POST("/transfer", TransferPost(ver))
	//服务器视频列表
	group.GET("/list", ListGet(ver))
	group.POST("/commit", CommitPost(ver))
	//查看状态
	group.GET("/status/:id", StatusGet(ver))

	return nil
}
