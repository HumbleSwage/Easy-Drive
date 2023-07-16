package consts

type shareType int

const (
	OneDay shareType = iota
	SevenDay
	ThirtyDay
	NoExpire
)

var shareTypeNum = []string{"1", "7", "30", "0"}

func (st shareType) String() string {
	return shareTypeNum[st]
}

func (st shareType) Index() int {
	return int(st)
}
