package db

import (
	"distributedStorage/model"
	"distributedStorage/orm"
	"fmt"
)

////用户信息表
//type RetUserFile struct {
//	UserName    string
//	FileHash    string
//	FileName    string
//	FileSize    string
//	UploadAt    string
//	LastUpdated string
//}

//插入一条数据至用户文件表
func OnUserFileUploadFinished(username, filehash, filename string, filesize int64) bool {
	dbEngine := orm.XormMysqlEngine
	_, err := dbEngine.Table("tbl_user_file").Insert(&model.UserFile{
		UserName: username,
		FileSHA1: filehash,
		FileSize: filesize,
		FileName: filename,
	})
	if err != nil {
		fmt.Println("插入数据至用户文件表失败:", err)
		return false
	}
	return true
}

func QueryUserFileMetas(username string, limit int) ([]model.RetUserFile, error) {
	dbEngine := orm.XormMysqlEngine
	userFiles := []model.RetUserFile{}
	err := dbEngine.Table("tbl_user_file").Where("user_name = ?", username).Limit(limit).Find(&userFiles)
	if err != nil {
		fmt.Println("查询用户文件表失败:", err)
		return []model.RetUserFile{}, err
	}
	return userFiles, err
}
