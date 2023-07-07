package main

import (
	"distributedStorage/service/account/handler"
	"distributedStorage/service/account/proto"
	"fmt"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/selector"
	"time"
)

//注册服务功能
var consulReg registry.Registry

// 负载均衡器
var rrSelector selector.Selector

func init() {
	//新建一个consul注册的地址，也就是我们consul服务启动的机器ip+端口
	consulReg = consul.NewRegistry(
		registry.Addrs("192.168.0.90:8500"),
	)

	rrSelector = selector.NewSelector(
		//策略为RoundRobin
		selector.SetStrategy(selector.RoundRobin),
	)

}

func main() {

	go startAccountService("192.168.0.1:8004")
	//go startAccountService("192.168.142.1:8005")
	//go startAccountService("8006")

	select {}
}

func startAccountService(addr string) {
	//consulReg := consul.NewRegistry(
	//	registry.Addrs("192.168.0.90:8500"),
	//)
	//
	//rrSelector := selector.NewSelector(
	//	//策略为RoundRobin
	//	selector.SetStrategy(selector.RoundRobin),
	//)

	service := micro.NewService(
		micro.Name("go.micro.service.user"),
		micro.Registry(consulReg),
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
		micro.Selector(rrSelector),
	)
	//service.Init()

	err := proto.RegisterUserServiceHandler(service.Server(), new(handler.User))
	if err != nil {
		fmt.Println("go.micro.service.user注册失败", err)
	}

	//for i := 0; i < 3; i++ {
	//	port := 8000 + i // 每个实例使用不同的端口号
	//	go func(port int) {
	//		service.Init(
	//			micro.Address(fmt.Sprintf(":%d", port)), // 设置不同的端口号
	//		)
	//		if err := service.Run(); err != nil {
	//			//log.Fatal(err)
	//		}
	//	}(port)
	//}

	//go func() {
	//	service.Init(
	//		micro.Address(fmt.Sprintf(":%d", 8001)), // 设置不同的端口号
	//	)
	//	err = service.Run()
	//	if err != nil {
	//		fmt.Println("go.micro.service.user启动失败", err)
	//	}
	//}()
	//

	service.Init(
		micro.Address(fmt.Sprintf(addr)), // 设置不同的端口号
	)
	err = service.Run()
	if err != nil {
		fmt.Println("go.micro.service.user启动失败", err)
	}

}
