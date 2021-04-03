package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/yyl-smpc/spdzgo/common/message"
	"github.com/yyl-smpc/spdzgo/common/ws"
	"log"
	"math"
	"math/rand"
	"time"
)

func main()  {

	address := []string{"ws://localhost:1008/ws",
		"ws://localhost:1009/ws",
	}
	clients,_ := connectTo(address)  //连接客户端
	for i := 0;i < 10000;i++ {
		time.Sleep(time.Millisecond*1000)
		sl := getSlice(address)   //获取切片
		tasks := getTask(sl, int64(i))
		for i := 0;i < len(sl);i++ {
			bys,_ := json.Marshal(tasks[i])
			clients[i].Conn.WriteMessage(websocket.TextMessage, bys)
		}
	}
}

func connectTo(address []string)  ([]*ws.Client,error){
	dialer := &websocket.Dialer{}
	n := len(address)
	clients := make([]*ws.Client,len(address))
	for i := 0;i < n;i++ {
		conn, _, err :=dialer.Dial(address[i]+"/task",nil)
		if err != nil {
			log.Println(err)
			return nil,err
		}
		clients[i]=  &ws.Client{Id:conn.RemoteAddr().String(),SendC: make(chan []byte),ReadC: make(chan []byte),Conn: conn,TimeStamp: time.Now().Unix(),Handler: nil}
		ws.Manager.AddClient(clients[i])
		go clients[i].Write()
		go clients[i].Read()
	}
	return clients,nil
}

func getSlice(address []string)  []message.SliceMessage{
	n := len(address)
	sl := make([]message.SliceMessage,n)
	alpha := int64(rand.Intn(100))
	xR := make([]int64, n)
	aR := make([]int64, n)
	cR := make([]int64, int64(math.Pow(2,float64(n))))

	addressMap := make(map[int64]string)

	for i := 0;i < n;i++{
		sl[i].AddressMap = addressMap
		addressMap[int64(i)] = address[i] + "/computer"
	}

	//随机获取x,a的分量
	for k := 0;k < n;k++ {
		sl[k].Id = int64(k)
		sl[k].X = make([]message.SecretSlice, n)
		sl[k].A = make([]message.SecretSlice, n)
		sl[k].C = make([]message.SecretSlice, int64(math.Pow(2,float64(n))))
		for i := 0;i < n;i++ {
			sl[k].X[i] = message.SecretSlice{0,0,0}
			sl[k].X[i][0] = 0
			sl[k].X[i][1] = int64(rand.Intn(100))
			xR[i] += sl[k].X[i][1]  //计算xR
			sl[k].X[i][2] = sl[k].X[i][1]*alpha
		}

		for i := 0;i < n;i++ {
			sl[k].A[i] = message.SecretSlice{0,0,0}
			sl[k].A[i][0] = 0
			sl[k].A[i][1] = int64(rand.Intn(100))
			aR[i] += sl[k].A[i][1]   //计算aR
			sl[k].A[i][2] = sl[k].A[i][1]*alpha
		}
	}

	//保存xR
	for i := 0;i < n;i++ {
		sl[i].X[i] = append(sl[i].X[i], xR[i])
	}

	//计算c
	for i := 0;i < len(sl[0].C);i++ {
		cR[i] = 1
		flag := 1
		for j := 0;j < n;j++ {
			if flag & i != 0 {
				cR[i] *= aR[j]
			}
			flag <<= 1
		}
	}

	//分割c
	for i := 0;i < n;i++ {
		for j := 0;j < len(cR);j++ {
			sl[i].C[j] = message.SecretSlice{0,0,0}
			sl[i].C[j][0] = 0
			sl[i].C[j][1] = int64(rand.Intn(100))
			sl[i].C[j][2] = sl[i].C[j][1]*alpha
			cR[j] -= sl[i].C[j][1] //更新cR
		}
	}

	for j := 0;j < len(cR);j++ {
		sl[0].C[j][1] += cR[j]
		sl[0].C[j][2] = sl[0].C[j][1]*alpha
	}
	return sl
}

func getTask(sl []message.SliceMessage,id int64) []message.Task{
	tasks := make([]message.Task, len(sl))
	for i := 0;i < len(sl);i++ {
		bys,_ := json.Marshal(sl[i])
		tasks[i] = message.Task{PId: 0,TaskId: id,Msg: bys}
	}
	return tasks
}
