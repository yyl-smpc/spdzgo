package handler

import "github.com/yyl-smpc/spdzgo/common"

type NumberHandler struct {
	task common.Task  //负责的任务
	c chan []byte     //通道
}

func NewNumberHandler(task common.Task, c chan []byte) NumberHandler{
	return NumberHandler{task, c}
}

//处理数字计算任务
func (r NumberHandler) Handler()  {

}


