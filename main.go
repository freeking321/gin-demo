package main

import (
	"fmt"
	"gin-demo/middleware"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// 结构体同时支持多种绑定模式

type Login struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

type LoginOut struct {
	User string `json:"user"`
	Password string `json:"password"`
}

func main() {

	// 9 写日志文件
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f) //输出日志到文件


	r := gin.Default()

	// 利用中间件写日志
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		// 你的自定义格式  -- 输出到了 gin.defaultwriter
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))


	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 1、获取路径中的参数
	// 2、获取get参数
	r.GET("/someGet", func(c *gin.Context) {
		// 两种方式
		firstName := c.DefaultQuery("first_name", "Guest")
		lastName := c.Query("last_name")
		c.String(http.StatusOK, "hello %s %s", firstName, lastName)
	})
	// 3、获取post参数
	r.POST("/somePost", func(c *gin.Context) {
		message := c.PostForm("message")
		nick := c.DefaultPostForm("nick", "any")

		c.JSON(200, gin.H{
			"status": "posted",
			"message": message,
			"nick": nick,
		})
	})
	// 4、get 和 post 参数混合  ---不推荐这样使用
	// 5、上传文件
	r.POST("/upload", func(c *gin.Context) {
		// 单文件
		file, _ := c.FormFile("file")
		log.Println(file.Filename)

		// 上传文件到指定的路径
		// c.SaveUploadedFile(file, dst)

		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})

	// 6、路由分组
	v1 := r.Group("/v1")
	{
		v1.POST("/login", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		v1.POST("/submit", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}
	// 7、无中间件
	// r := gin.New()

	// 8、全局中间件
	r.Use(gin.Logger())
	r.Use(middleware.MyLogger())
	// 路由单独设置中间件
	//r.GET("/benchmark", MyBenchLogger()， benchEndpoint)
	// 路由组中添加中间件
	// authorized := r.Group("/", AuthRequired())
	// 路由组摸部分设置中间件
	//authorized := r.Group("/")
	//authorized.Use(AuthRequired())
	//{
	//	authorized.POST("/login", loginEndpoint)
	//	authorized.POST("/submit", submitEndpoint)
	//	authorized.POST("/read", readEndpoint)
	//
	//	// nested group
	//	testing := authorized.Group("testing")
	//	testing.GET("/analytics", analyticsEndpoint)
	//}

	// 10 支持 json
	r.POST("/login/json", func(c *gin.Context) {
		var json LoginOut
		// 验证json格式
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if json.User != "yangliu" || json.Password != "123456" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})

	// 11、自定义验证器--这个先不看
	// 12、shouldbindquery 只绑定 get 参数 -- 没啥用哦
	// 13、shouldbindquery 绑定get或者post参数 -- 不大理解的
	// 14 还能绑定uri
	// 15 绑定 post 参数

	// 16 静态文件路由
	r.StaticFile("/demo1.png", "./resources/demo1.png")
	// 17 重定向
	r.GET("/test", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://www.google.com/")
	})
	// 重定向到其他路由
	r.GET("/test", func(c *gin.Context) {
		c.Request.URL.Path = "/test2"
		r.HandleContext(c)
	})
	r.GET("/test2", func(c *gin.Context) {
		c.JSON(200, gin.H{"hello": "world"})
	})

	r.GET("/test/middleware", func(c *gin.Context) {
		example := c.MustGet("example").(string)

		// it would print: "12345"
		log.Println(example)
	})

	// 18 basicauth 中间件
	//authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
	//	"foo":    "bar",
	//	"austin": "1234",
	//	"lena":   "hello2",
	//	"manu":   "4321",
	//}))

	// 19自定义 http 配置
	//router := gin.Default()
	//s := &http.Server{
	//	Addr:           ":8080",
	//	Handler:        router,
	//	ReadTimeout:    10 * time.Second,
	//	WriteTimeout:   10 * time.Second,
	//	MaxHeaderBytes: 1 << 20,
	//}
	//s.ListenAndServe()

	// 20 还能运行多个服务
	// 21 优雅的运行和停止服务


	r.Run(":8300") // listen and serve on 0.0.0.0:8080
}