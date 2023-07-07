package db

import (
	"distributedStorage/model"
	"distributedStorage/orm"
	"fmt"
	"time"
)

func OnFileUploadFinished(filehash, filename string, filesize int64, fileaddr string) bool {
	dbEngine := orm.XormMysqlEngine
	_, err := dbEngine.Table("tbl_file").Insert(&model.File{
		FileSHA1: filehash,
		FileName: filename,
		FileSize: filesize,
		FileAddr: fileaddr,
		Status:   1,
	})
	if err != nil {
		fmt.Println("tbl_file表插入失败: ", err)
	}

	return true
}

//CREATE TABLE `tbl_file` (
//`id` int(11) NOT NULL AUTO_INCREMENT,
//`file_sha1` char(40) NOT NULL DEFAULT '' COMMENT '文件hash',
//`file_name` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
//`file_size` bigint(20) DEFAULT '0' COMMENT '文件大小',
//`file_addr` varchar(1024) NOT NULL DEFAULT '' COMMENT '文件存储位置',
//`create_at` datetime default NOW() COMMENT '创建日期',
//`update_at` datetime default NOW() on update current_timestamp() COMMENT '更新日期',
//`status` int(11) NOT NULL DEFAULT '0' COMMENT '状态(可用/禁用/已删除等状态)',
//`ext1` int(11) DEFAULT '0' COMMENT '备用字段1',
//`ext2` text COMMENT '备用字段2',
//PRIMARY KEY (`id`),
//UNIQUE KEY `idx_file_hash` (`file_sha1`),
//KEY `idx_status` (`status`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8;

type File struct {
	ID        int       `xorm:"'id' autoincr"`
	FileSHA1  string    `xorm:"'file_sha1' notnull unique(idx_file_hash)"`
	FileName  string    `xorm:"'file_name'"`
	FileSize  int64     `xorm:"'file_size'"`
	FileAddr  string    `xorm:"'file_addr'"`
	CreatedAt time.Time `xorm:"'create_at' created"`
	UpdatedAt time.Time `xorm:"'update_at' updated"`
	Status    int       `xorm:"'status' index(idx_status)"`
	Ext1      int       `xorm:"'ext1'"`
	Ext2      string    `xorm:"'ext2'"`
}

func GetFileMeta(filehash string) (*model.RetFileTable, error) {
	dbEngine := orm.XormMysqlEngine
	ret := model.RetFileTable{}
	exist, err := dbEngine.Table("tbl_file").Where("file_sha1=?", filehash).Get(&ret)
	if !exist {
		fmt.Println("文件元信息不存在,filehash:", filehash)
		return nil, err
	}
	return &ret, err
}

func UpdateFileLocation(filehash string, fileaddr string) bool {
	dbEngine := orm.XormMysqlEngine
	_, err := dbEngine.Table("tbl_file").
		Cols("file_addr").
		Where("file_sha1=?", filehash).
		Update(&model.File{FileAddr: fileaddr})
	if err != nil {
		fmt.Printf("更新文件location失败, filehash:%s", filehash)
		return false
	}

	return true
}
