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
	"strconv"
)

// StatusUploaded 已上传
const StatusUploaded = "uploaded"

// StatusQueuing 队列中
const StatusQueuing = "queuing"

// StatusTransferring 转换中
const StatusTransferring = "transferring"

// StatusFileWrong 文件错误
const StatusFileWrong = "wrong file"

// StatusFinished 完成
const StatusFinished = "finished"

/**
 * @apiDefine Success
 * @apiSuccess {string} msg 返回具体消息
 * @apiSuccess {int} code 返回状态码：【正常：0】，【失败，-1】
 * @apiSuccess {json} [detail]  如正常则返回detail
 */
/**
 * @apiDefine Failed
 * @apiErrorExample {json} Error-Response:
 *     {
 *       "code":-1,
 *       "msg":"error message",
 *     }
 */
const _ = "apiDefine"

//UploadPost 文件上传接口
/**
* @api {post} /v1/upload 上传文件接口
* @apiName upload
* @apiGroup Upload
* @apiVersion  0.0.1
*
* @apiUse Success
* @apiParam  {Binary} binary 媒体文件二进制文件
* @apiParamExample  {Binary} Request-Example:
*    upload a binary file from local
* @apiSuccess (detail) {string} id 文件名ID
* @apiSuccessExample {json} Success-Response:
*     {
*       "code":0,
*       "msg":"ok",
*       "detail":{
*			"id":"9FCp2x2AeEWNobvzKA3vRgqzZNqFWEJTMpLAz2hLhQGEd3URD5VTwDdTwrjTu2qm"
*		 }
*     }
* @apiUse Failed
* @apiSampleRequest /v1/upload
 */
func UploadPost(vertion string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer ctx.Request.Body.Close()
		src := "./upload/"
		fileName, err := writeTo(src, ctx.Request.Body)
		if err != nil {
			resultFail(ctx, err.Error())
			return
		}

		client.Set(fileName, StatusUploaded, 0)
		log.Println(fileName)
		ctx.JSON(http.StatusOK, JSON(0, "ok", gin.H{"id": fileName}))
		return
	}
}

// TransferPost 视频转换处理
/**
* @api {post} /v1/transfer 视频转换处理
* @apiName transfer
* @apiGroup Transfer
* @apiVersion  0.0.1
*
* @apiUse Success
* @apiParam  {string} id 文件名ID
* @apiParam  (string) url key存放url地址
* @apiParam  (string) m3u8 m3u8文件名（暂不支持）
* @apiParam  (string) key key文件名（暂不支持）
* @apiParam  (string) keyInfo keyInfo文件名（暂不支持）
* @apiParamExample  {string} Request-Example:
* {
*     "id":"9FCp2x2AeEWNobvzKA3vRgqzZNqFWEJTMpLAz2hLhQGEd3URD5VTwDdTwrjTu2qm",
*     "url":"http://localhost:8080/transfer/xxx/key"
* }
*
* @apiSuccess (detail) {string} id 文件名ID
* @apiSuccessExample {json} Success-Response:
*     {
*       "code":0,
*       "msg":"ok",
*       "detail":{
*			"id":"9FCp2x2AeEWNobvzKA3vRgqzZNqFWEJTMpLAz2hLhQGEd3URD5VTwDdTwrjTu2qm"
*		 }
*     }
* @apiUse Failed
* @apiSampleRequest /v1/transfer
 */
func TransferPost(version string) gin.HandlerFunc {

	src := config.Upload + "/"
	dst := config.Transfer + "/"
	return func(ctx *gin.Context) {
		b, err := openssl.HexKey()
		if err != nil {
			resultFail(ctx, err.Error())
			return
		}
		//probe := ffprobe.New(src + id)
		//probe.IsH264AndAAC()
		//TODO:file is not a media file
		id := ctx.PostForm("id")
		if id == "" {
			resultFail(ctx, "wrong id request")
			return
		}
		url := ctx.PostForm("url")
		if url == "" {
			url = config.URL + "/" + config.Transfer + "/" + id + "/key"
		}

		stream := NewStreamer(string(b), id)
		stream.SetURI(url)
		stream.SetDst(dst)
		stream.SetSrc(src)
		client.Set(id, StatusQueuing, 0)
		queue.Push(stream)
		resultOK(ctx, gin.H{"id": id})
	}
}

// InfoGet 获取视频转换状态
/**
*
* @api {get} /v1/info/:id 获取视频转换状态
* @apiName info
* @apiGroup Info
* @apiVersion  0.0.1
*
* @apiParam  {String} id 文件名ID
*
* @apiUse Success
* @apiSuccess  {string} code 返回状态码：【异常错误：-1】，【正常：0】，【文件不存在：1】,【处理中：2】，【文件异常：3】，【队列中：4】，
* @apiSuccess  {json} [detail] 正常则返回detail
* @apiSuccess (detail) {string} uri 视频存放的相对地址
* @apiSuccess (detail) {string} m3u8 m3u8存放的文件名
* @apiSuccess (detail) {string} key key存放的文件名
* @apiSuccess (detail) {string} keyInfo keyInfo存放的文件名
*
* @apiSampleRequest /v1/info/:id
* @apiParamExample  {string} Request-Example:
* 	http://localhost:8080/v1/info/9FCp2x2AeEWNobvzKA3vRgqzZNqFWEJTMpLAz2hLhQGEd3URD5VTwDdTwrjTu2qm
*
* @apiSuccessExample {json} Success-Response OK:
* {
*       "code":0,
*       "msg":"ok",
*       "detail":{
*			"uri":"transfer/xxx",
*			"m3u8":"media.m3u8",
*			"key":"key"
*			"keyInfo":"KeyInfo",
*		}
*
* }
* @apiSuccessExample {json} Success-Response NoData:
* {
*       "code":1,
*       "msg":"data not found",
* }
* @apiSuccessExample {json} Success-Response Transferring:
* {
*       "code":2,
*       "msg":"transferring",
* }
* @apiSuccessExample {json} Success-Response FileWrong:
* {
*       "code":3,
*       "msg":"wrong file",
* }
* @apiUse Failed
 */
func InfoGet(version string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		val, err := client.Get(id).Result()
		if err != nil {
			ctx.JSON(http.StatusOK, JSON(1, "data not found"))
			return
		}
		if val != StatusTransferring {
			ctx.JSON(http.StatusOK, JSON(2, val))
			return
		} else if val == StatusFileWrong {
			ctx.JSON(http.StatusOK, JSON(3, val))
			return
		} else if val != StatusQueuing {
			ctx.JSON(http.StatusOK, JSON(4, val))
			return
		}
		resultOK(ctx, gin.H{
			"uri":     config.Transfer + "/" + id,
			"m3u8":    config.M3U8,
			"key":     config.KeyFile,
			"keyInfo": config.KeyInfoFile,
		})
		return
	}
}

// ListGet 获取所有视频列表
/**
*
* @api {get} /v1/list 获取所有视频列表
* @apiName list
* @apiGroup List
* @apiVersion  0.0.1
*
* @apiParam  {string} start 列表开始Index,从0开始计数
* @apiParam  {string} end 列表结束Index,默认取100
*
* @apiUse Success
* @apiSuccess  {string} code 返回状态码：【正常：0】，【处理中：1】
*
* @apiSampleRequest /v1/list
* @apiParamExample  {string} Request-Example:
* 	http://localhost:8080/v1/list
*
* @apiSuccessExample {json} Success-Response OK:
* {
*       "code":0,
*       "msg":"ok",
* }
* @apiUse Failed
 */
func ListGet(ver string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//client.Append()
		//client.
		list, err := findDir(ctx.Query("start"), ctx.Query("end"))
		if err != nil {
			resultFail(ctx, err.Error())
			return
		}
		resultOK(ctx, gin.H{
			"start": list.Start,
			"end":   list.End,
			"max":   list.Max,
			"list":  list.List,
		})
	}
}

// DirList ...
type DirList struct {
	Start int
	End   int
	Max   int
	List  []string
}

const limit = 100

func findDir(start, end string) (*DirList, error) {
	//var err error
	st, err := strconv.Atoi(start)
	if err != nil {
		st = 0
	}
	ed, err := strconv.Atoi(end)
	if err != nil {
		ed = st + limit
	}
	d, err := os.Open(config.Transfer)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer d.Close()
	fi, err := d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var dirs []string
	max := len(fi)
	if st >= max {
		st = max
	}
	if ed >= max {
		ed = max
	}

	for i := st; i < ed; i++ {
		log.Println(fi[i].Name(), fi[i].IsDir())
		if fi[i].IsDir() {
			dirs = append(dirs, fi[i].Name())
		}
	}
	return &DirList{
		Start: st,
		End:   ed,
		Max:   max,
		List:  dirs,
	}, nil
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

func resultOK(ctx *gin.Context, h ...gin.H) {
	if h != nil {
		ctx.JSON(http.StatusOK, JSON(0, "ok", h...))
		return
	}
	ctx.JSON(http.StatusOK, JSON(0, "ok"))
}

func resultFail(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, JSON(-1, msg))
}
