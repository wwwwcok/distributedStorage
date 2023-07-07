package mq

import "fmt"

var done chan bool

//StartConsume：消费者(客户端)开始监听队列，获取消息
func StartConsume(qName, CName string, callback func(msg []byte) bool) {
	//1. 通过channel.Consumme获取消息信道
	msgs, err := channel.Consume(
		qName,
		CName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println("客户端创建监听消息队列的信道失败：", err)
		return
	}
	//2.循环从信道获取消息
	done = make(chan bool)
	go func() {
		for msg := range msgs {
			//3. 获取到信息消息就调用callback方法处理
			procssSuc := callback(msg.Body)
			if !procssSuc {
				//TODO:将认为写到另一个队列，用于异常情况处理
			}
		}
	}()
	//阻塞
	<-done
	//关闭mq的队列
	channel.Close()
}
