package main

import (
	"distributedStorage/service/download/config"
	download_proto "distributedStorage/service/download/proto"
	"distributedStorage/service/download/route"
	"distributedStorage/service/download/rpc"
	"fmt"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
)

var consulReg registry.Registry

func init() {
	consulReg = consul.NewRegistry(
		registry.Addrs("192.168.0.90:8500"),
	)

}

func startRPCService() {
	service := micro.NewService(
		micro.Name("go.micro.service.download"),
		micro.Registry(consulReg),
	)

	service.Init()

	err := download_proto.RegisterDownloadServiceHandler(service.Server(), &rpc.Downloader{})
	if err != nil {
		fmt.Println("go.micro.service.download服务注册consul失败", err)
	}
	err = service.Run()
	if err != nil {
		fmt.Println("go.micro.service.download服务启动", err)
	}
}

func startAPIDownloadSrv() {
	r := route.Router()
	r.Run(config.DownloadServiceHost)
}

func main() {
	//开启rpc服务
	go startRPCService()
	//处理真正的服务
	startAPIDownloadSrv()

}
