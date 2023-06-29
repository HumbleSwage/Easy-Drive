package test

import (
	"easy-drive/consts"
	"fmt"
	"testing"
)

func TestFileCategory(t *testing.T) {
	fmt.Println(consts.Video)
	fmt.Println(consts.Video.Index())
	fmt.Println(consts.Video.String())
}

func TestFileStatus(t *testing.T) {
	fmt.Println(consts.DeletedFile)
	fmt.Println(consts.DeletedFile.Index())
	fmt.Println(consts.DeletedFile.String())
}
