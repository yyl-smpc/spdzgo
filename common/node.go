package node

//节点的抽象

type Node interface {
	GetId() string
	GetDomain() string
	GetPort() string
}