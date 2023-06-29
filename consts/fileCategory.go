package consts

type category int

const (
	Video category = iota
	Music
	Image
	Doc
	Others
)

var CategoryStr = []string{"video", "music", "image", "doc", "others"}

func (c category) String() string {
	return CategoryStr[c]
}

func (c category) Index() int {
	return int(c)
}
