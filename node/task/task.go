package task

import (
	"encoding/json"
	"github.com/yyl-smpc/spdzgo/common/message"
	"github.com/yyl-smpc/spdzgo/common/ws"
	"log"
	"time"
)

type Task struct {

}

func (task *Task)Handler(msg []byte)  {
	taskMsg := &message.Task{}
	json.Unmarshal(msg, taskMsg)
	secretMsg := &message.SliceMessage{

	}
	json.Unmarshal(taskMsg.Msg, secretMsg)
	if ws.Manager.GetHandler(taskMsg.TaskId) != nil {
		return
	}
	sender := make(chan byte)
	reader := make(chan []byte,1024)
	ws.Manager.AddHandler(taskMsg.TaskId, reader)
	defer func() {
		ws.Manager.DeleteHandler(taskMsg.TaskId)
	}()

	executor := &Executor{
		TaskId: taskMsg.TaskId,
		Reader: reader,
		Sender: sender,
		SecretMsg: secretMsg,
		messageMap: map[int64][]*message.Task{},
	}
	go executor.Mul()
	ticker := time.Tick(time.Second*1000)
	for {
		select {
		case <-sender:
			return
		case <-ticker:
			log.Println("计算超时")
			return
		}
	}
}
