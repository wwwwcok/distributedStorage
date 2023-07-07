package download

import (
	"distributedStorage/meta"
	"distributedStorage/store/oss"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

func Download(ctx *gin.Context) {
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

func DownloadURL(ctx *gin.Context) {
	fsha1 := ctx.Query("filehash")
	fm, _ := meta.GetFileMetaDB(fsha1)

	signedURL := oss.DownloadURL(fm.Location)
	fmt.Println("oss返回的下载URL:", signedURL, "filehash:", fsha1)
	ctx.Writer.WriteString(signedURL)
}
