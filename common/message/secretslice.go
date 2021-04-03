package message

type SecretSlice []int64
type Address string
type SliceMessage struct {
	Id int64  //序号
	AddressMap map[int64]string //序号对应的地址
	X []SecretSlice
	A []SecretSlice //a的分片
	C []SecretSlice //ck的分片值
}
