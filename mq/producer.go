package mq

import (
	"distributedStorage/config"
	"fmt"
	"github.com/streadway/amqp"
)

var conn *amqp.Connection
var channel *amqp.Channel

// 如果异常关闭，会接收通知
var notifyClose chan *amqp.Error

func init() {
	//是否开启异步转移功能，开启时才初始化rabbitMQ连接
	if !config.AsyncTransferEnable {
		return
	}
	if initChannel() {
		channel.NotifyClose(notifyClose)
	}
	//断线自动重连
	//go func() {
	//	for {
	//		select {
	//		case msg := <-notifyClose:
	//			conn = nil
	//			channel = nil
	//			log.Printf("onNotifyChannelClosed: %+v\n", msg)
	//			initChannel()
	//		}
	//	}
	//}()
}

func initChannel() bool {
	//1判断channel是否已经创建过
	if channel != nil {
		return true
	}
	//2获得mq的一个连接
	conn, err := amqp.Dial(config.RabbitURL)
	if err != nil {
		fmt.Println("获得rabbitmq的连接失败", err)
		return false
	}
	//3.打开一个channel,用于消息的发布和接收
	channel, err = conn.Channel()
	if err != nil {
		fmt.Println("获得rabbitmq的Channel失败", err)
		return false
	}
	return true
}

//生产者Publisher:发布消息功能
func Publish(exchange, routingKey string, msg []byte) bool {
	//1. 检查channel是否正常
	if !initChannel() {
		return false
	}

	//2 执行消息发布动作
	err := channel.Publish(exchange, routingKey, false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})
	if err != nil {
		fmt.Println("rabbitmq的执行消息发布动作失败：", err)
		return false
	}

	return true
}
