package route

import (
	"distributedStorage/service/apigw/apigwhandler"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()

	router.Static("/static", "./static")

	router.GET("/user/signup", apigwhandler.RegisterIndex)
	router.POST("/user/signup", apigwhandler.SignUp)

	//router.GET("/file/upload/suc", apigwhandler.SuccessUpload)

	router.GET("/file/upload", apigwhandler.UploadIndex)
	router.POST("/file/upload", apigwhandler.GetUploadEntry)
	router.POST("/file/donwload", apigwhandler.GetDownloadEntry)
	return router
}
