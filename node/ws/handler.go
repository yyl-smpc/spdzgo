package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/yyl-smpc/spdzgo/common/ws"
	"github.com/yyl-smpc/spdzgo/node/task"
	"log"
	"net/http"
	"time"
)

type Handler struct {

}

func (handler Handler)Init(router *gin.Engine) {
	router.GET("ws/computer", func(c *gin.Context) {
		conn,err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
			return true
		}}).Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println(err)
			return
		}
		client := &ws.Client{Id: conn.RemoteAddr().String(),SendC: make(chan []byte),ReadC:make(chan []byte),Conn: conn, TimeStamp: time.Now().Unix(),Handler: &task.Executor{}} //创建连接客户端
		ws.Manager.AddClient(client)                                                                                             //注册客户端
		go client.Read()
		go client.Write()
	})


	router.GET("ws/task", func(c *gin.Context) {
		conn,err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
			return true
		}}).Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println(err)
			return
		}
		client := &ws.Client{Id:conn.RemoteAddr().String(),SendC: make(chan []byte),ReadC:make(chan []byte),Conn: conn, TimeStamp: time.Now().Unix(), Handler: &task.Task{}} //创建连接客户端
		ws.Manager.AddClient(client)
		go client.Read()
		go client.Write()
	})
}
