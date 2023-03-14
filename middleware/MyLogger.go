package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

// 自定义的一个日志中间件

func MyLogger() gin.HandlerFunc  {
	return func(c *gin.Context) {
		t := time.Now() //当前时间戳

		c.Set("example", "123456") // set参数
		// before request

		c.Next()

		// after request
		latency := time.Since(t)
		log.Print(latency)

		status := c.Writer.Status()
		log.Println(status)
	}
}