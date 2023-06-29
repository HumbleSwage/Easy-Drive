package consts

type uploadMode int

const (
	FastUpload uploadMode = iota
	OnUpload
	Uploaded
)

var modeStr = []string{"upload_seconds", "uploading", "upload_finish"}

func (m uploadMode) String() string {
	return modeStr[m]
}

func (m uploadMode) Index() int {
	return int(m)
}
