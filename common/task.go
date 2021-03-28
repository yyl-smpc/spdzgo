package common

//task是所有计算任务的抽象

type Task struct {
	Id string          //任务id
	UserId string      //发起任务的用户id
	Feedback string    //结果反馈地址
	Pointers []Pointer //数据指针
	operation string   //数据操作 "add, mul, mod"
	timestamp string   //任务时间戳
}