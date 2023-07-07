package controller

import (
	"distributedStorage/cache/redis"
	"distributedStorage/db"
	"distributedStorage/model"
	"distributedStorage/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type MultiuploadController struct {
}

func (m *MultiuploadController) Route(engine *gin.Engine) {

	engine.POST("/file/mpupload/init", m.InitalMultipartUpload)
	engine.POST("/file/mpupload/uppart", m.UploadPart)
	engine.POST("/file/mpupload/complete", m.CompleteUpload)
}

func (m *MultiuploadController) InitalMultipartUpload(ctx *gin.Context) {
	username := ctx.PostForm("username")
	filehash := ctx.PostForm("filehash")
	fmt.Println("InitalMultipartUpload接口的参数解析部分:", username, filehash)
	filesize, err := strconv.ParseInt(ctx.PostForm("filesize"), 10, 64)
	if err != nil {
		fmt.Println("filesize错误", err)
		return
	}

	//3. 文件分块
	upInfo := model.MultipartUploadInfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadID:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  5 * 1024 * 1024,
		ChunkCount: int(math.Ceil(float64(filesize) / (5 * 1024 * 1024))),
	}

	//4. 文件分块的元信息写入redis
	err = redis.RedisCli.HSet("MP_"+upInfo.UploadID, "chunkcount", upInfo.ChunkCount).Err()
	fmt.Println("InitalMultipartUpload接口写入redis失败:", err)
	redis.RedisCli.HSet("MP_"+upInfo.UploadID, "filehash", upInfo.FileHash)
	redis.RedisCli.HSet("MP_"+upInfo.UploadID, "filesize", upInfo.FileSize)
	ctx.Writer.Write((&util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: upInfo,
	}).JSONBytes())
}

func (m *MultiuploadController) UploadPart(ctx *gin.Context) {
	//1. 解析参数
	//username := ctx.PostForm("username")
	uploadID := ctx.Query("uploadid")
	chunkIndex := ctx.Query("index")
	fmt.Println("UploadPart接口的参数解析部分:", uploadID, chunkIndex)
	//2. 获得文件句柄，用于存储分块内容
	fpath := "./data/" + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)
	defer fd.Close()
	if err != nil {
		fmt.Println("创建文件错误", err, fpath)
		ctx.Writer.Write((&util.RespMsg{
			Code: -1,
			Msg:  "Upload part failed",
			Data: nil,
		}).JSONBytes())
	}

	buf := make([]byte, 1024*1024)
	for {
		n, err := ctx.Request.Body.Read(buf)
		defer ctx.Request.Body.Close()
		//注意！一定要先写到文件中，再判断err，因为ctx.Request.Body.Read读取到EOF也算err,最后一个切片还没来及写入文件就会提取因错误判断打断
		fd.Write(buf[:n])
		if err != nil {
			fmt.Println("读取Body错误", err)
			break
		}
	}

	//3. 更新redis缓存状态
	redis.RedisCli.HSet("MP_"+uploadID, "chkidx_"+chunkIndex, 1)
	//4. 返回信息给客户端
	ctx.Writer.Write((&util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: nil,
	}).JSONBytes())
}
func (m *MultiuploadController) CompleteUpload(ctx *gin.Context) {
	//1 .解析参数
	uploadID := ctx.PostForm("uploadid")
	username := ctx.PostForm("username")
	filehash := ctx.PostForm("filehash")
	filesize, _ := strconv.ParseInt(ctx.PostForm("filesize"), 10, 64)
	filename := ctx.PostForm("filename")
	fmt.Println("CompleteUpload接口的参数解析部分:", uploadID, username, filehash, filesize, filename)
	//2. 通过uploadid查询redis并判断是否所有分开上传完成
	totalCount := 0
	chunkCount := 0
	data := redis.RedisCli.HGetAll("MP_" + uploadID)
	res, err := data.Result()
	if err != nil {
		fmt.Println("CompleteUpload接口：从redis中获取分块信息失败")
	}
	for k, v := range res {
		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount++
		}
	}

	if totalCount != chunkCount {
		ctx.Writer.Write((&util.RespMsg{
			Code: -2,
			Msg:  "invalid request",
			Data: nil,
		}).JSONBytes())
		return
	}
	//3. TODO: 合并分块
	//4. 更新唯一文件表及用户文件表
	db.OnFileUploadFinished(filehash, filename, filesize, "")
	db.OnUserFileUploadFinished(username, filehash, filename, filesize)
	//5. 响应客户端
	ctx.Writer.Write((&util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: nil,
	}).JSONBytes())
}
