package task

import "spdzgo/common/node"

//task是所有计算任务的抽象

type Task interface {
	GetId() string
	GetParticipators() []node.DataNode
}