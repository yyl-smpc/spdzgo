package ws

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

//处理消息的方法
type MessageHandler interface {
	Handler([]byte)
}

type Client struct {
	Id string
	SendC  chan []byte
	ReadC  chan []byte
	Conn   *websocket.Conn
	TimeStamp  int64
	Handler MessageHandler
}


func (client *Client)Read(){
	defer client.Close()
	for {
		_,msg,err := client.Conn.ReadMessage()
		if err != nil{
			log.Println(err)
			break
		}
		client.TimeStamp = time.Now().Unix()  //更新时间戳
		client.Handler.Handler(msg)
	}
}

func (client *Client)Write()  {
	defer client.Close()
	for {
		msg := <- client.SendC
		err := client.Conn.WriteMessage(websocket.TextMessage,msg)
		if err != nil {
			log.Println(err)
			break
		}
		client.TimeStamp = time.Now().Unix()  //更新时间戳
	}
}

func (client *Client)Close() {
	log.Println("断开连接--"+client.Id)
	Manager.DeleteClient(client.Id)
	client.Conn.Close()
}
