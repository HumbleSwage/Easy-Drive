package test

import (
	"easy-drive/conf"
	"easy-drive/pkg/utils"
	"easy-drive/pkg/utils/fileUtil"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

func TestGenerateEmailCaptcha(t *testing.T) {
	str := utils.GenerateEmailCaptcha()
	fmt.Println(str)
}

func TestUploadFileToLocalStatic(t *testing.T) {
	var a int64
	fmt.Println(a)
}

func TestGenerateThumbnail(t *testing.T) {
	utils.InitLog()
	f := " /Users/zhaodeng/Desktop/ff/头像.jpg"
	name := fileUtil.GenerateThumbnailName("头像.jpg")
	fmt.Println(name)
	err := fileUtil.GenerateThumbnail(f, name)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("图片缩略图测试成功")
	}
}

func TestGenerateThumbnailName(t *testing.T) {
	filePath := "/Users/zhaodeng/Desktop/ff/Springboot vue3 仿百度网盘（后端）项目实战 计算机毕业设计 简历项目 - 024 - 管理端分享.mp4"
	cmd := exec.Command("/Users/zhaodeng/opt/ffmpeg/ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", filePath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(output))
}

func TestSplitVideoByPercentage(t *testing.T) {
	file := "/Users/zhaodeng/Desktop/ff/电竞视频.mp4"
	err := fileUtil.SplitVideo(file)
	if err != nil {
		fmt.Println(err)
	}
}

func TestCreateDir(t *testing.T) {
	fileUtil.CreateDir("Desktop/esay-drive/user1/test")
}

func TestGetUniqueFileName(t *testing.T) {
	conf.InitConfig()
	filePath := fileUtil.GetFilePath("电竞视频1.mp4", "a83809b4cdc6", "电竞视频1.mp4")
	fmt.Println(filePath)
}

func TestGetUniqueFileName123(t *testing.T) {
	file := "chunk_1687958392710243000_1c5ebef0-6b98-45f6-895c-0ca7cdc5d48d"
	pattern := `^chunk_\d+_.+$`
	regExp, _ := regexp.Compile(pattern)
	fmt.Println(regExp.MatchString(file))
}

func TestGetFileExtension(t *testing.T) {
	parts := strings.Split("test.mp4", ".")
	if len(parts) > 1 {
		fmt.Println(parts[len(parts)-1])
	} else {
		fmt.Println("出错")
	}
}
