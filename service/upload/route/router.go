package route

import (
	"distributedStorage/service/upload/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Router() *gin.Engine {
	router := gin.Default()
	router.Static("/static", "./static")

	router.POST("/file/upload", api.Upload)
	router.OPTIONS("/file/upload", func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "POST,OPTIONS")
		ctx.Status(http.StatusNoContent)

	})
	return router
}
