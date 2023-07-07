package model

import "time"

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

type RetFileTable struct {
	FileHash string `json:"filehash" xorm:"file_sha1"`
	FileName string `json:"filename" xorm:"file_name"`
	FileSize int64  `json:"filesize" xorm:"file_size"`
	FileAddr string `json:"fileaddr" xorm:"file_addr"`
}

type User struct {
	ID             int       `xorm:"'id' autoincr pk"`
	UserName       string    `xorm:"'user_name' notnull default('') unique(idx_username)"`
	UserPwd        string    `xorm:"'user_pwd' notnull default('')"`
	Email          string    `xorm:"'email' default('')"`
	Phone          string    `xorm:"'phone' default('')"`
	EmailValidated bool      `xorm:"'email_validated' default(0)"`
	PhoneValidated bool      `xorm:"'phone_validated' default(0)"`
	SignupAt       time.Time `xorm:"'signup_at' created"`
	LastActive     time.Time `xorm:"'last_active' updated"`
	Profile        string    `xorm:"'profile'"`
	Status         int       `xorm:"'status' notnull default(0) index(idx_status)"`
}

type UserToken struct {
	ID        int    `xorm:"'id' notnull autoincr pk"`
	UserName  string `xorm:"'user_name' notnull default '' comment('用户名') unique(idx_username)"`
	UserToken string `xorm:"'user_token' notnull default '' comment('用户登录token')"`
}

type RetUser struct {
	Username     string `xorm:"user_name"`
	Email        string `xorm:"email"`
	Phone        string `xorm:"phone"`
	SignupAt     string `xorm:"signup_at"`
	LastActiveAt string `xorm:"last_active"`
	Status       int    `xorm:"status"`
}

type UserFile struct {
	ID         int       `xorm:"'id' notnull pk autoincr"`
	UserName   string    `xorm:"'user_name' notnull index(idx_user_id) unique(idx_user_file)"`
	FileSHA1   string    `xorm:"'file_sha1' notnull default '' comment('文件hash') index(user_file) unique(idx_user_file)"`
	FileSize   int64     `xorm:"'file_size' default 0 comment('文件大小')"`
	FileName   string    `xorm:"'file_name' notnull default '' comment('文件名')"`
	UploadAt   time.Time `xorm:"'upload_at' created comment('上传时间')"`
	LastUpdate time.Time `xorm:"'last_update' created updated comment('最后修改时间')"`
	Status     int       `xorm:"'status' notnull default 0 comment('文件状态(0正常1已删除2禁用)') index(idx_status)"`
}

type RetUserFile struct {
	UserName    string `xorm:"user_name"`
	FileHash    string `xorm:"file_sha1"`
	FileName    string `xorm:"file_name"`
	FileSize    int64  `xorm:"file_size"`
	UploadAt    string `xorm:"upload_at"`
	LastUpdated string `xorm:"last_update"`
}

type MultipartUploadInfo struct {
	FileHash   string `xorm:"file_hash"`
	FileSize   int64  `xorm:"file_hash"`
	UploadID   string `xorm:"UploadID"`
	ChunkSize  int    `xorm:"ChunkSize"`
	ChunkCount int    `xorm:"ChunkCount"`
}
