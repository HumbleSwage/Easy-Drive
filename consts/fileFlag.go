package consts

type flag int

const (
	DeletedFile flag = iota
	RestoreFile
	NormalFile
)

var flagStr = []string{"删除", "回收站", "正常"}

func (s flag) String() string {
	return flagStr[s]
}

func (s flag) Index() int {
	return int(s)
}
