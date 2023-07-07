package controller

import (
	"distributedStorage/db"
	"distributedStorage/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type MemberController struct {
}

const (
	pwd_salt = "*#890"
)

func (m *MemberController) Route(engine *gin.Engine) {
	engine.GET("/user/signup", m.RegisterIndex)
	engine.POST("/user/signup", m.SignUp)
	engine.POST("/user/signin", m.SignIn)
	engine.GET("/user/signin", m.SignInIndex)
	engine.GET("/user/signin/home", m.HomeIndex)

	engine.POST("/user/info", m.Userlnfo)
}

func (m *MemberController) RegisterIndex(ctx *gin.Context) {
	data, err := ioutil.ReadFile("../../static/view/signup.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	io.WriteString(ctx.Writer, string(data))
	return
}

//处理用户注册请求
func (m *MemberController) SignUp(ctx *gin.Context) {
	username := ctx.PostForm("username")
	pwd := ctx.PostForm("password")

	encPwd := util.Sha1([]byte(pwd + pwd_salt))
	suc := db.UserSignUp(username, encPwd)

	if suc {
		ctx.Writer.WriteString("SUCCESS")
	} else {
		ctx.Writer.WriteString("FAILED")
	}
}

func (m *MemberController) SignInIndex(ctx *gin.Context) {
	data, err := ioutil.ReadFile("../../static/view/signin.html")
	if err != nil {
		fmt.Println("读取home.html数据错误", err)
		return
	}
	ctx.Writer.WriteString(string(data))
	return
}

func (m *MemberController) SignIn(ctx *gin.Context) {
	username := ctx.PostForm("username")
	pwd := ctx.PostForm("passwd")
	encPwd := util.Sha1([]byte(pwd + pwd_salt))
	checked := db.UserSignIn(username, encPwd)
	if !checked {
		ctx.Writer.WriteString("FAILED")
		return
	}

	//2. 生成token
	token := GenToken(username)
	upRes := db.UpdateToken(username, token)
	if !upRes {
		ctx.Writer.WriteString("FAILED")
		return
	}

	//3. 登录成功跳转到首页
	//ctx.Writer.WriteString("http://127.0.0.1:62200/user/signin/home")
	//ctx.Redirect(301, "http://127.0.0.1:62200/user/signin/home")

	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://127.0.0.1:62200/user/signin/home",
			Username: username,
			Token:    token,
		},
	}

	ctx.Writer.Write(resp.JSONBytes())
	return
}

func (m *MemberController) Userlnfo(ctx *gin.Context) {
	//j解析请求参数
	//username := ctx.PostForm("username")
	//token := ctx.PostForm("token")
	username := ctx.Query("username")
	token := ctx.Query("token")

	//判断token是否有效
	isValid := IsTokenVaild(token)
	if !isValid {
		ctx.Status(http.StatusForbidden)
		return
	}

	//查询用户信息
	user, err := db.GetUserInfo(username)
	if err != nil {
		ctx.Status(http.StatusForbidden)
		return
	}
	if user == nil {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "Failed",
			Data: nil,
		}
		ctx.Writer.Write(resp.JSONBytes())
		return
	}
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	ctx.Writer.Write(resp.JSONBytes())
}

func (m *MemberController) HomeIndex(ctx *gin.Context) {
	data, err := ioutil.ReadFile("../../static/view/home.html")
	if err != nil {
		fmt.Println("读取home.html数据错误", err)
		return
	}
	//ctx.Writer.WriteString(string(data))
	io.WriteString(ctx.Writer, string(data))
	return

}

func GenToken(username string) string {
	//40位toekn:username+timestamp+token_salt+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	return util.MD5([]byte(username + ts + "_tokensalt" + ts[:8]))

}

func IsTokenVaild(token string) bool {

	//判断token是否过期
	//从数据库里查询两个token是否一致
	return true
}
