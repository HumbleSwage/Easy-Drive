package consts

type fileCategory int

const (
	VideoCategory fileCategory = iota + 1
	MusicCategory
	ImageCategory
	DocCategory
	OthersCategory
)

var fileCategoryStr = []string{"video", "music", "image", "doc", "others"}

func (c fileCategory) String() string {
	return fileCategoryStr[c-1]
}

func (c fileCategory) Index() int {
	return int(c)
}
