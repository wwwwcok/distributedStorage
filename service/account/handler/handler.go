package handler

import (
	"context"
	"distributedStorage/common"
	"distributedStorage/config"
	"distributedStorage/db"
	"distributedStorage/service/account/proto"
	"distributedStorage/util"
)

type User struct {
}

func (U *User) Signup(ctx context.Context, req *proto.ReqSignup, res *proto.RespSignup) error {
	username := req.Username
	pwd := req.Password

	encPwd := util.Sha1([]byte(pwd + config.PasswordSalt))
	suc := db.UserSignUp(username, encPwd)

	if suc {
		//ctx.Writer.WriteString("SUCCESS")
		res.Code = common.StatusOK
		res.Message = "注册成功"
	} else {
		//ctx.Writer.WriteString("FAILED")
		res.Code = common.StatusRegisterFailed
		res.Message = "注册失败"
	}
	return nil
}
