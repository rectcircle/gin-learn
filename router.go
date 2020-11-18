package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func handler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ok",
	})
}

func routerMethod(r *gin.Engine) {
	r.GET("router/method", handler)
	r.POST("router/method", handler)
	r.PUT("router/method", handler)
	r.PATCH("router/method", handler)
	r.DELETE("router/method", handler)
	r.HEAD("router/method", handler)
	r.OPTIONS("router/method", handler)
	r.Any("router/any", handler)  // 匹配所有方法
	// r.Handle // 最终调动的函数
}

func routerRequestParam(r *gin.Engine) {
	// 可以匹配 router/request/path/require/1 router/request/path/require/1/
	// 无法可以匹配 router/request/path/require/ router/request/path/require
	r.GET("router/request/path/require/:requirePathParam", func(c *gin.Context) {
		c.String(http.StatusOK, "requirePathParam = %s", c.Param("requirePathParam"))
	})
	// 可以匹配 router/request/path/remain/1 router/request/path/remain/1/
	// 可以匹配 router/request/path/remain/ router/request/path/remain
	r.GET("router/request/path/remain/*remainPathParam", func(c *gin.Context) {
		c.String(http.StatusOK, "remainPathParam = %s", c.Param("remainPathParam"))
	})

	// Query 参数
	// http://127.0.0.1:8080/router/request/query?queryParam=123&queryParam=321&queryArr=1&queryArr=2&queryMap[a]=1&queryMap[b]=2
	// 返回 queryParam = 123, queryParamWithDefault = default, queryArr = [1 2], queryMap = map[a:1 b:2]
	r.GET("router/request/query", func(c *gin.Context) {
		c.String(http.StatusOK, "queryParam = %s, queryParamWithDefault = %s, queryArr = %s, queryMap = %s", 
			c.Query("queryParam"), 
			c.DefaultQuery("queryParamWithDefault", "default"),
			c.QueryArray("queryArr"),
			c.QueryMap("queryMap"),
		)
	})

	r.GET("router/request/bind", func(c *gin.Context) {
		// curl http://127.0.0.1:8080/router/request/bind\?User\=xiaoming\&password\=312
		// user = {xiaoming  */*}
		user := struct {
			User     string `json:"user"`
			Password string
			Accept string
		} {}
		c.BindHeader(&user)
		c.Bind(&user)
		c.String(http.StatusOK, "user = %s", user)
	})

	r.GET("router/request/validate", func(c *gin.Context) {
		pagination := struct {
			PageNo  uint32 `form:"pageNo" binding:"required,min=1"`
			PageSize uint32 `form:"pageSize" binding:"required,min=1,max=100"`
		} {}
		if err := c.Bind(&pagination); err == nil {
			c.String(http.StatusOK, "pagination = %s", pagination)
		} else {
			c.String(http.StatusOK, "%s", err)
		}
	})
}

func routerResponse(r *gin.Engine) {
	responseStruct := struct {
		Name  string
		Email string
	}{
		Name:  "xiaoming",
		Email: "xiaoming@example.com",
	}

	r.GET("router/response/json", func(c *gin.Context) {
		c.JSON(http.StatusOK, responseStruct)
	})
	r.GET("router/response/yaml", func(c *gin.Context) {
		c.YAML(http.StatusOK, responseStruct)
	})
	r.GET("router/response/xml", func(c *gin.Context) {
		c.XML(http.StatusOK, responseStruct)
	})
}

func routerGroup(r *gin.Engine) {
	handler := func(c *gin.Context) {
		c.String(http.StatusOK, "Hello")
	}
	// 简单的路由组: v1
	v1 := r.Group("/router/group/v1")
	{
		v1.GET("/hello", handler)
	}

	// 简单的路由组: v2
	v2 := r.Group("/router/group/v2")
	{
		v2.GET("/hello", handler)
	}
	// curl http://127.0.0.1:8080/router/group/v1/hello
	// curl http://127.0.0.1:8080/router/group/v2/hello
}
func MyLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// 设置 example 变量
		c.Set("example", "12345")

		// 请求前

		c.Next()

		// 请求后
		latency := time.Since(t)

		// 获取发送的 status
		status := c.Writer.Status()

		log.Printf("latency=%s, status=%d", latency, status)
	}
}

func routerWithMiddleware(r *gin.Engine) {
	group := r.Group("router/group/middleware")
	group.Use(MyLogger())
	group.GET("/hello", func(c *gin.Context) {
		example := c.MustGet("example").(string)
		c.String(http.StatusOK, "example = %s", example)
	})
	// curl http://127.0.0.1:8080/router/group/middleware/hello
}

func registeRouter(r *gin.Engine) {
	routerMethod(r)
	routerRequestParam(r)
	routerResponse(r)
	routerGroup(r)
	routerWithMiddleware(r)
}