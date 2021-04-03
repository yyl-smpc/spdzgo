package message

type Task struct {
	TaskId int64 //任务号
	Seq int64  //消息序列号
	PId  int64 // 任务负责人消息发送者id
	Msg []byte  //传输的数据
}

