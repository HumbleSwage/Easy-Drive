package consts

type fileType int

const (
	VIDEO fileType = iota + 1
	MUSIC
	IMAGE
	PDF
	WORD
	EXCEL
	TXT
	PROGRAM
	ZIP
	OTHERS
)

var CategoryStr = []string{"video", "music", "image", "pdf", "word", "excel", "txt", "program", "zip", "others"}

func (c fileType) String() string {
	return CategoryStr[c-1]
}

func (c fileType) Index() int {
	return int(c)
}

var TypeMapping = map[int][]string{
	VIDEO.Index():   {"mp4", "avi", "rmvb", "mkv", "mov"},
	MUSIC.Index():   {"mp3", "wav", "wma", "mp2", "flac", "midi", "ra", "ape", "aac", "cda"},
	IMAGE.Index():   {"jpeg", "jpg", "png", "gif", "bmp", "dds", "psd", "pdt", "webp", "xmp", "svg", "tiff"},
	PDF.Index():     {"pdf"},
	WORD.Index():    {"docx"},
	EXCEL.Index():   {"xlsx"},
	TXT.Index():     {"txt"},
	PROGRAM.Index(): {"h", "c", "hpp", "hxx", "cpp", "cc", "c++", "m", "o", "s", "dll", "cs", "java", "class", "js", "ts", "css", "scss", "vue", "jsx", "sql", "md", "json", "html", "xml"},
	ZIP.Index():     {"rar", "zip", "7z", "cab", "arj", "lzh", "tar", "gz", "ace", "uue", "bz", "jar", "iso", "mpq"},
	OTHERS.Index():  {},
}
