package task

import (
	"encoding/json"
	"github.com/yyl-smpc/spdzgo/common/message"
	"github.com/yyl-smpc/spdzgo/common/ws"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type Executor struct {
	TaskId int64//该执行者执行的任务号
	Reader chan []byte //读取数据的通道
	Sender chan byte//发送数据的通道
	SecretMsg *message.SliceMessage //保存任务初始消息
	messageMap map[int64][]*message.Task  //保存计算过程消息
	lock sync.Mutex
}

func (executor *Executor)Read(end chan byte)  {
	for {
		select {
		case msg := <-executor.Reader:
			executor.addMessage(msg)
		case <-end:
			return
		}
	}
}
func (executor *Executor)Handler(msg []byte)  {
	taskMessage := &message.Task{}
	json.Unmarshal(msg, taskMessage)
	for i := 0;i < 10;i++ {
		sender := ws.Manager.GetHandler(taskMessage.TaskId)
		if sender != nil{
			sender <- msg
			break
		}
		time.Sleep(time.Millisecond*100)  //等待一百毫秒
	}
	return
}

func (executor *Executor)addMessage(msg []byte)  {
	executor.lock.Lock()
	taskMessage := &message.Task{}
	json.Unmarshal(msg, taskMessage)
	if executor.messageMap[taskMessage.Seq] == nil {
		executor.messageMap[taskMessage.Seq] = make([]*message.Task,0)
	}
	executor.messageMap[taskMessage.Seq] = append(executor.messageMap[taskMessage.Seq], taskMessage)
	executor.lock.Unlock()
}

func (executor *Executor)PopMessage(seq int64)  []*message.Task{
	executor.lock.Lock()
	defer executor.lock.Unlock()
	msg := executor.messageMap[seq]
	delete(executor.messageMap,seq)
	return msg
}

func (executor *Executor)Mul()  {
	end := make(chan byte)
	defer func() {
		end<-byte(1)
	}()
	go executor.Read(end)  //开始读取数据
	x := int64(rand.Intn(100))  //本地的数据
	localId := executor.SecretMsg.Id
	n := len(executor.SecretMsg.X)  //人数

	clients := make([]*ws.Client, 0)

	for i := 0; i < n;i++ {
		if int64(i) != localId {
			client := ws.Manager.GetClient(executor.SecretMsg.AddressMap[int64(i)])
			if client == nil {
				executor.Sender<-2
				return
			}
			clients = append(clients, client)
		}
	}

	//发送x-r数据
	executor.SecretMsg.X[localId][0] =-x+executor.SecretMsg.X[localId][3]
	executor.SecretMsg.X[localId][1] +=x-executor.SecretMsg.X[localId][3]
	ws.Manager.Broadcast(clients,createTask(executor.TaskId, localId,0, []byte(strconv.FormatInt(executor.SecretMsg.X[localId][0],10))))  //发送消息

	//更新rou
	count := 0
	for true {
		if count == n - 1 {
			break
		}
		res := executor.PopMessage(0)
		for _,v := range res {
			executor.SecretMsg.X[v.PId][0], _ = strconv.ParseInt(string(v.Msg), 10, 64)
		}
		count += len(res)
	}

	e := make([]message.SecretSlice, n)

	for i := 0;i < n;i++ {
		e[i] = message.SecretSlice{0,0,0}
		e[i][0] = executor.SecretMsg.X[i][0] - executor.SecretMsg.A[i][0]  //更新rou
		e[i][1] = executor.SecretMsg.X[i][1] - executor.SecretMsg.A[i][1]  //更新r
		e[i][2] = executor.SecretMsg.X[i][2] - executor.SecretMsg.A[i][2]  //更新mac(r)
	}

	eBytes,_ := json.Marshal(e)
	ws.Manager.Broadcast(clients, createTask(executor.TaskId, localId,1, eBytes))

	//更新e
	count = 0
	for true {
		if count == n -1 {
			break
		}
		res := executor.PopMessage(1)
		for _,v := range res {
			e1 := make([]message.SecretSlice, n)
			err := json.Unmarshal(v.Msg, &e1)
			if  err != nil {
				log.Println(v.Seq)
				log.Println(err)
				return
			}
			for i := 0; i < n; i++ {
				e[i][1] += e1[i][1]
				e[i][2] += e1[i][2]
			}
		}
		count += len(res)
	}

	//计算乘积的分片
	for i := 0;i < len(executor.SecretMsg.C);i++ {
		flag := 1
		for j := 0; j < n;j++ {
			if flag&i == 0 {
				executor.SecretMsg.C[i][0] *= e[j][1]
				executor.SecretMsg.C[i][1] *= e[j][1]
				executor.SecretMsg.C[i][2] *= e[j][1]
			}
			flag <<= 1
		}
	}

	xMul := [3]int64{}

	for i := 0;i < len(executor.SecretMsg.C); i++ {
		xMul[0] += executor.SecretMsg.C[i][0]
		xMul[1] += executor.SecretMsg.C[i][1]
		xMul[2] += executor.SecretMsg.C[i][2]
	}

	//发送x乘积的切片
	xMulBytes,_ := json.Marshal(xMul)
	ws.Manager.Broadcast(clients, createTask(executor.TaskId, localId,2,xMulBytes))

	//计算乘积
	count = 0
	for true {
		if count == n - 1 {
			break
		}
		res := executor.PopMessage(2)
		for _,v := range res {
			xMul1 := [3]int64{}
			json.Unmarshal(v.Msg, &xMul1)
			xMul[0] += xMul1[0]
			xMul[1] += xMul1[1]
			xMul[2] += xMul1[2]
		}
		count += len(res)
	}

	log.Printf("TASK_ID:%d:xMul=%v",executor.TaskId,xMul)
	executor.Sender<-byte(1)
}

func createTask(taskId int64, pid int64,seq int64, msg []byte) ([]byte){
	task := message.Task{
		TaskId: taskId,
		PId: pid,
		Msg: msg,
		Seq: seq,
	}
	b ,_ := json.Marshal(task)
	return b
}
