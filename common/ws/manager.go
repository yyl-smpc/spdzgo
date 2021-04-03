package ws

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

type ClientManager struct {
	Clients map[string]*Client
	Handler map[int64]chan []byte
	cLock sync.Mutex
	hLock sync.Mutex
}

func (clientManager *ClientManager)AddHandler(taskId int64 ,handlerC chan []byte)  {
	defer clientManager.hLock.Unlock()
	clientManager.hLock.Lock()
	clientManager.Handler[taskId] = handlerC
}

func (clientManager *ClientManager)GetHandler(taskId int64 )  chan []byte{
	defer clientManager.hLock.Unlock()
	clientManager.hLock.Lock()
	if _,ok := clientManager.Handler[taskId];ok {
		return clientManager.Handler[taskId]
	}
	return nil
}


func (clientManager *ClientManager)DeleteHandler(taskId int64 ) {
	defer clientManager.hLock.Unlock()
	clientManager.hLock.Lock()
	if _,ok := clientManager.Handler[taskId];ok {
		delete(clientManager.Handler, taskId)
	}
}

func (clientManager *ClientManager)AddClient(client *Client)  {
	defer clientManager.cLock.Unlock()
	clientManager.cLock.Lock()
	clientManager.Clients[client.Id] = client
}

func (clientManager *ClientManager)GetClient(address string)  *Client{
	defer clientManager.cLock.Unlock()
	clientManager.cLock.Lock()
	if _,ok := clientManager.Clients[address];ok {
		return clientManager.Clients[address]
	} else {
		dialer := &websocket.Dialer{}
		conn, _,err := dialer.Dial(address, nil)
		if err != nil {
			log.Printf(err.Error())
			return nil
		}
		client := &Client{Id: address,SendC: make(chan []byte),ReadC:make(chan []byte),Conn: conn, TimeStamp: time.Now().Unix(),Handler: nil} //创建连接客户端
		Manager.Clients[address] = client                                                                                               //注册客户端
		go client.Read()
		go client.Write()
		return client
	}
}

func (clientManager *ClientManager)DeleteClient(address string) {
	defer clientManager.cLock.Unlock()
	clientManager.cLock.Lock()
	if _,ok := clientManager.Clients[address];ok {
		delete(clientManager.Clients, address)
	}
}

func (clientManager *ClientManager)Broadcast(clients []*Client, msg []byte)  {
	for _,v := range clients{
		v.SendC <- msg
	}
}



//创建一个管理器实例
var Manager = &ClientManager{
	Clients: make(map[string]*Client),
	Handler: make(map[int64]chan []byte),
	cLock: sync.Mutex{},
	hLock: sync.Mutex{},
}


