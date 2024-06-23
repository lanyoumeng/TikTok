package main

import (
	"net/http"
)

func main() {
	// 创建一个默认的 Gin 路由器
	r := gin.Default()

	// 设置一个简单的 GET 路由
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	// 启动 HTTP 服务器，默认监听在 0.0.0.0:8080
	r.Run()
}
