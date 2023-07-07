package main

import (
	"distributedStorage/service/apigw/route"
)

func main() {
	r := route.Router()
	r.Run("127.0.0.1:62201")
}
