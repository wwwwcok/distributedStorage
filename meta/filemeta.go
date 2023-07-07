package meta

import (
	"distributedStorage/db"
	"fmt"
)

//文件元信息结构
type FileMeta struct {
	FileSha1 string `json:"filehash" xorm:"file_sha1"`
	FileName string `json:"filename" xorm:"file_name"`
	FileSize int64  `json:"filesize" xorm:"file_size"`
	Location string `json:"fileaddr" xorm:"file_addr"`
	UploadAt string `xorm:"'upload_at' created comment('上传时间')"`
}

var fileMetas map[string]FileMeta = make(map[string]FileMeta)

func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

func RemoveFileMeta(fileSha1 string) {

	delete(fileMetas, fileSha1)
}

func UpdateFileMetaDB(fmeta *FileMeta) {
	db.OnFileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

//获取文件元信息:mysql方式
func GetFileMetaDB(fileSha1 string) (*FileMeta, error) {
	file, err := db.GetFileMeta(fileSha1)
	if err != nil {
		fmt.Println("GetFileMeta接口出错：", err)
		return nil, err
	}

	if file == nil {
		return nil, err
	}

	return &FileMeta{
		FileSha1: file.FileHash,
		FileName: file.FileName,
		FileSize: file.FileSize,
		Location: file.FileAddr,
	}, err
}
