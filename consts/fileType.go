package consts

type fileType int

const (
	// 音频
	MP3 fileType = iota
	WAV
	FLAC
	AAC
	OGG
	WMA

	// 视频
	MP4
	AVI
	MKV
	MOV
	WMV
	FLV

	// 图片
	JPG
	JPEG
	PNG
	GIF
	BMP
	TIFF
	SVG

	// 文档
	PDF
	DOCX
	XLSX
	PPTX
	TXT
	CSV
)

var TypeStr = []string{"mp3", "wav", "flac", "aac", "ogg", "mp4", "avi", "mkv", "mov", "wmv", "flv", "jpg", "jpeg", "png", "gif", "bmp", "tiff", "svg", "pdf", "docx", "xlsx", "pptx", "txt", "csv"}

func (f fileType) String() string {
	return TypeStr[f]
}

func (f fileType) Index() int {
	return int(f)
}
