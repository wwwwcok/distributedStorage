package apigwhandler

import (
	"context"
	"distributedStorage/service/account/proto"
	download_proto "distributedStorage/service/download/proto"

	go_micro_service_upload "distributedStorage/service/upload/proto"
	"distributedStorage/util"
	"fmt"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
)

var (
	userCli     proto.UserService
	uploadCli   go_micro_service_upload.UploadService
	downloadCli download_proto.DownloadService
)

func init() {

	consulReg := consul.NewRegistry(
		registry.Addrs("192.168.0.90:8500"),
	)

	service := micro.NewService(
		micro.Registry(consulReg),
	)
	//初始化、解析命令行参数等(V3版本不需要了，统统在NewService()传入参数完成)
	//service.Init()
	//初始化一个rpcClient; user.pb.micro.go里分别定义了生成client和handler的函数
	userCli = proto.NewUserService("go.micro.service.user", service.Client())
	//初始化一个用于访问upload服务的rpcClient
	uploadCli = go_micro_service_upload.NewUploadService("go.micro.service.upload", service.Client())

	downloadCli = download_proto.NewDownloadService("go.micro.service.download", service.Client())

}

//文件存放路径
//var curdir, _ = os.Getwd()
//var filedir = filepath.FromSlash(curdir + "/../../tmp/")

func UploadIndex(ctx *gin.Context) {
	data, err := ioutil.ReadFile("../../static/view/index.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	io.WriteString(ctx.Writer, string(data))
	return
}

func GetUploadEntry(ctx *gin.Context) {

	//处理跨域访问，权限授权
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "POST,OPTIONS")

	resp, err := uploadCli.UploadEntry(context.TODO(), &go_micro_service_upload.ReqEntry{})
	if err != nil {
		fmt.Println("获取uploadEntry地址失败", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":  0,
		"msg":   "获取uploadEntry地址成功",
		"entry": resp.Entry,
	})

}

func GetDownloadEntry(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "POST,OPTIONS")
	resp, _ := downloadCli.DownloadEntry(context.TODO(), &download_proto.ReqEntry{})

	cliResp := util.RespMsg{
		Code: 0,
		Msg:  "获取DownloadEntry成功",
		Data: resp.Entry,
	}
	ctx.JSON(200, cliResp)

}

//响应get请求
func RegisterIndex(ctx *gin.Context) {
	data, err := ioutil.ReadFile("../../static/view/signup.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	io.WriteString(ctx.Writer, string(data))
	return
}

//处理用户注册请求
func SignUp(ctx *gin.Context) {
	username := ctx.PostForm("username")
	pwd := ctx.PostForm("password")

	resp, err := userCli.Signup(context.TODO(), &proto.ReqSignup{
		Username: username,
		Password: pwd,
	})
	if err != nil {
		fmt.Println("rpc客户端处理'用户注册'请求失败", err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": resp.Code,
		"msg":  resp.Message,
	})
}
