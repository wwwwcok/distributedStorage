package main

import (
	"distributedStorage/service/upload/config"
	go_micro_service_upload "distributedStorage/service/upload/proto"
	"distributedStorage/service/upload/route"
	"distributedStorage/service/upload/rpc"
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
		micro.Name("go.micro.service.upload"),
		micro.Registry(consulReg),
	)
	service.Init()
	err := go_micro_service_upload.RegisterUploadServiceHandler(service.Server(), new(rpc.Upload))
	if err != nil {
		fmt.Println("go_micro_service_upload服务注册失败", err)
	}
	err = service.Run()
	if err != nil {
		fmt.Println("go_micro_service_upload服务启动失败", err)
	}

}

func startAPIService() {
	r := route.Router()
	r.Run(config.UploadServiceHost)
}

func main() {
	//开启rpc服务
	go startRPCService()
	//真正的上传服务
	startAPIService()
}
