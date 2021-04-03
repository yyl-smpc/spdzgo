package http

import "github.com/gin-gonic/gin"

type Handler struct {

}

func (handler Handler)Init(router *gin.Engine)  {
	router.GET("/index", func(c *gin.Context) {
		c.Writer.Write([]byte("hello index"))
	})
}

