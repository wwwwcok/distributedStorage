package db

import (
	"distributedStorage/model"
	"distributedStorage/orm"
	"fmt"
)

//用户通过用户名密码注册
func UserSignUp(username string, passwd string) bool {
	dbEngine := orm.XormMysqlEngine
	_, err := dbEngine.Table("tbl_user").Insert(&model.User{
		UserName: username,
		UserPwd:  passwd,
	})
	if err != nil {
		fmt.Println("用户注册失败: ", err)
		return false
	}
	return true
}

func UserSignIn(username string, encpwd string) bool {
	dbEngine := orm.XormMysqlEngine
	user := model.User{}
	bool, err := dbEngine.Table("tbl_user").Where("user_name=?", username).Get(&user)
	if err != nil {
		fmt.Println("UserSignIn接口查询用户过程中出错", err)
		return false
	}
	if !bool {
		fmt.Println("没有当前用户")
		return false
	}
	if user.UserPwd != encpwd {
		fmt.Println("密码错误")
		return false
	}
	return true
}

//刷新用户登录token
func UpdateToken(username string, token string) bool {
	dbEngine := orm.XormMysqlEngine
	_, err := dbEngine.Table("tbl_user_token").Cols("user_token").Where("user_name=?", username).Update(&model.UserToken{UserToken: token})
	if err != nil {
		fmt.Println("更新用户token过程中出错", err)
		return false
	}
	return true
}

func GetUserInfo(username string) (*model.RetUser, error) {
	dbEngine := orm.XormMysqlEngine
	user := model.RetUser{}

	bool, err := dbEngine.Table("tbl_user").Where("user_name = ?", username).Get(&user)
	if err != nil {
		fmt.Println("GetUserInfo接口查询用户信息过程中出错", err)
		return nil, err
	}
	if !bool {
		fmt.Println("GetUserInfo接口没有查询到用户信息", "用户名："+username)
		return nil, err
	}

	return &user, err
}
