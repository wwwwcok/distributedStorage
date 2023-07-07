package controller

import (
	"distributedStorage/config"
	"distributedStorage/db"
	"distributedStorage/meta"
	"distributedStorage/mq"
	"distributedStorage/store/oss"
	"distributedStorage/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type FileController struct {
}

//文件存放路径
var curdir, _ = os.Getwd()
var filedir = curdir + "/../../tmp/"

func (f *FileController) Route(Engine *gin.Engine) {
	Engine.GET("/file/upload/suc", f.successUpload)
	Engine.GET("/file/upload", f.uploadIndex)
	Engine.POST("/file/upload", f.upload)
	Engine.GET("/file/meta", f.GetFileMeta)
	Engine.GET("/file/download", f.Download)
	Engine.POST("/file/update", f.FileMetaUpdate)
	Engine.GET("/file/delete", f.FileDelete)
	Engine.POST("/file/query", f.FileQuery)
	Engine.POST("/file/fastupload", f.TryFastUpload)

	Engine.POST("/file/downloadurl", f.DownloadURL)

}

func (f *FileController) uploadIndex(ctx *gin.Context) {
	data, err := ioutil.ReadFile("../../static/view/index.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	io.WriteString(ctx.Writer, string(data))
	return
}

func (f *FileController) upload(ctx *gin.Context) {
	formFile, err := ctx.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return
	}

	file, err := formFile.Open()
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	//timeStamp := fmt.Sprintf("%d", time.Now().UnixMicro())
	fileMeta := meta.FileMeta{
		FileName: formFile.Filename,
		Location: filedir + formFile.Filename,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	//创建文件接收数据
	created, err := os.Create(fileMeta.Location)
	if err != nil {
		fmt.Println(err)
		os.Mkdir("./tmp", 0777)
		return
	}
	defer created.Close()

	fileMeta.FileSize, _ = io.Copy(created, file)
	//光标移动到文件初始位置，下一个操作的字节
	created.Seek(0, io.SeekStart)

	fileMeta.FileSha1 = util.FileSha1(created)

	//将文件写入到oss存储
	created.Seek(0, io.SeekStart)
	ossPath := "oss/" + fileMeta.FileSha1
	//@上传到oss当中
	//err = oss.Bucket().PutObject(ossPath, created)
	//if err != nil {
	//	fmt.Println("上传oss失败")
	//	return
	//}
	//fileMeta.Location = ossPath

	//@上传到mq的生产者
	data := mq.TransferData{
		FileHash:      fileMeta.FileSha1,
		CurLocation:   fileMeta.Location,
		DestLocation:  ossPath,
		DestStoreType: 3,
	}
	pubData, _ := json.Marshal(data)
	suc := mq.Publish(
		config.TransExchangeName,
		config.TransOSSRoutingKey,
		pubData,
	)
	if !suc {
		//TODO：加入重发消息逻辑
	}

	//---保存文件信息表中(暂时不用)
	//---meta.UpdateFileMeta(fileMeta)

	//保存文件信息到唯一文件表中(mysql)
	meta.UpdateFileMetaDB(&fileMeta)

	//存到唯一文件表后还需要保存信息到用户文件表中
	username := ctx.Query("username")
	suc = db.OnUserFileUploadFinished(username, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize)

	if !suc {
		ctx.Writer.WriteString("Upload Failed.")
		return
	}

	fmt.Println("上传时的FileSha1:filehash=", fileMeta.FileSha1)
	ctx.Redirect(301, "http://127.0.0.1:62200/file/upload/suc")
}

func (f *FileController) successUpload(ctx *gin.Context) {
	//io.WriteString(ctx.Writer, "上传成功！")
	data, err := ioutil.ReadFile("../../static/view/home.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	io.WriteString(ctx.Writer, string(data))
	return

}

func (f *FileController) GetFileMeta(ctx *gin.Context) {
	fhash := ctx.Query("filehash")
	fmt.Println("fhash：", fhash)
	//表方式
	//fMeta := meta.GetFileMeta(fhash)
	//mysql方式
	fMeta, err := meta.GetFileMetaDB(fhash)
	if err != nil {
		return
	}
	fmt.Println("fMeta：", fMeta)
	ctx.JSON(200, fMeta)
}

func (f *FileController) Download(ctx *gin.Context) {
	fsha1 := ctx.Query("filehash")
	fm := meta.GetFileMeta(fsha1)
	fSrc, err := os.Open(fm.Location)
	defer fSrc.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	//data, _ := io.ReadAll(fSrc)
	//记得一定得先设置响应头,再设置响应内容
	ctx.Header("Content-Type", "application/octect-stream")
	ctx.Header("content-disposition", "attachment;filename=\""+fm.FileName+"\"")
	io.Copy(ctx.Writer, fSrc)
	//这样也行
	//ctx.Data(200, "application/octect-stream", data)

}

func (f *FileController) FileMetaUpdate(ctx *gin.Context) {
	opType := ctx.PostForm("op")
	fileSha1 := ctx.PostForm("filehash")
	newFileName := ctx.PostForm("filename")
	if opType == "0" {
		ctx.Status(http.StatusForbidden)
		return
	}

	curFile := meta.GetFileMeta(fileSha1)
	curFile.FileName = newFileName
	meta.UpdateFileMeta(curFile)

	curFile.Location = filedir + newFileName
	//修改后的文件元信息返回给浏览器
	ctx.JSONP(200, curFile)
}

//删除文件
func (f *FileController) FileDelete(ctx *gin.Context) {
	fh := ctx.Query("filehash")
	fm := meta.GetFileMeta(fh)
	meta.RemoveFileMeta(fm.FileSha1)

	os.Remove(fm.Location)

}

func (f *FileController) TryFastUpload(ctx *gin.Context) {

	//1.解析请求参数
	username := ctx.Query("username")
	filehash := ctx.Query("filehash")
	filename := ctx.Query("filename")
	filesize := ctx.Query("filesize")
	//2.从文件表里查询相同hash的文件记录
	fileMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		fmt.Println(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	//3. 查不到返回妙传失败
	if fileMeta == nil {
		fmt.Println("查不到返回妙传失败", fileMeta)
		resp := util.RespMsg{
			Code: -1,
			Msg:  "秒传失败，请访问普通上传接口",
			Data: nil,
		}
		ctx.Writer.Write(resp.JSONBytes())
	}

	fmt.Printf("秒传中查询到的文件云信息:%#v", fileMeta)

	//4. 上传过就将文件信息写入用户文件表
	size, _ := strconv.ParseInt(filesize, 10, 64)
	suc := db.OnUserFileUploadFinished(username, filehash, filename, size)
	if suc {
		resp := util.RespMsg{
			Code: 0,
			Msg:  "秒传成功",
			Data: nil,
		}
		ctx.Writer.Write(resp.JSONBytes())
	}

}

//查询批量的用户文件表元信息
func (f *FileController) FileQuery(ctx *gin.Context) {
	limitCnt, _ := strconv.Atoi(ctx.Request.FormValue("limit"))
	fmt.Println("FileQuery接口的limitCnt:", limitCnt)
	username := ctx.Query("username")
	users, err := db.QueryUserFileMetas(username, limitCnt)
	if err != nil {
		fmt.Println("批量获取用户文件元信息失败：", err)
		ctx.Status(500)
		return
	}
	//ctx.JSON(200, users)
	data, err := json.Marshal(&users)
	if err != nil {
		ctx.Status(500)
		return
	}

	ctx.Writer.Write(data)

}

//生成oss文件的下载地址
func (f *FileController) DownloadURL(ctx *gin.Context) {
	fsha1 := ctx.Query("filehash")
	fm, _ := meta.GetFileMetaDB(fsha1)

	signedURL := oss.DownloadURL(fm.Location)
	fmt.Println("oss返回的下载URL:", signedURL, "filehash:", fsha1)
	ctx.Writer.WriteString(signedURL)
}
