package common

type Location struct {
	Node     Node   //数据存储的主机
	ObjectId string //数据对象id
	Key      string //键值
}

type Pointer struct {
	ty       string   //指针类型
	location Location //数据的地址
}

type NumberPointer struct {
	Pointer
	sign string //"+,-.~",正，负，逆
}