package consts

type status int

const (
	Transfer status = iota
	TransferFailed
	Using
)

var StatusStr = []string{"TRANSFER", "TRANSFER_FAIL", "USING"}

func (s status) String() string {
	return StatusStr[s]
}

func (s status) Index() int {
	return int(s)
}
