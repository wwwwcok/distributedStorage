package main

import (
	"bufio"
	"distributedStorage/config"
	"distributedStorage/db"
	"distributedStorage/mq"
	"distributedStorage/store/oss"
	"encoding/json"
	"fmt"
	"os"
)

func ProcessTransfer(msg []byte) bool {
	//1 解析msg
	pubData := mq.TransferData{}
	err := json.Unmarshal(msg, &pubData)
	if err != nil {
		fmt.Println("回调函数ProcessTransfer解析msg失败：", err)
		return false
	}
	//2 根据msg找到临时存储文件路径，创建文件句柄
	filed, err := os.Open(pubData.CurLocation)
	if err != nil {
		pwd, _ := os.Getwd()
		fmt.Println("ProcessTransfer服务：创建文件句柄失败：", err, pubData.CurLocation, "当前目录：", pwd)
		return false
	}
	//3 通过文件句柄将文件读出并且上传到OSS
	err = oss.Bucket().PutObject(
		pubData.DestLocation,
		bufio.NewReader(filed),
	)
	if err != nil {
		fmt.Println("ProcessTransfer服务：将文件上传到OSS失败：", err)
		return false
	}
	//4 更新文件表信息，将文件存储路径修改为oss上的存储路径
	suc := db.UpdateFileLocation(
		pubData.FileHash,
		pubData.DestLocation,
	)
	if !suc {
		fmt.Println("ProcessTransfer服务：将文件存储路径修改为oss上的存储路径：", err)
		return false
	}
	return true
}

func main() {
	mq.StartConsume(
		config.TransOSSQueueName,
		"transfer_oss",
		ProcessTransfer,
	)

}
