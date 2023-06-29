package fileUtil

import (
	"easy-drive/consts"
	"easy-drive/pkg/utils"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func GenerateThumbnail(inputFilePath, outputFileName string) error {
	// 获取工作目录
	workingDir, err := os.Getwd()
	if err != nil {
		utils.LogrusObj.Error("获取工作目录失败:", err)
		return err
	}

	// 获取文件类型
	fileType := GetFileType(inputFilePath)

	// 获取输入文件的目录和文件名
	dir := filepath.Dir(inputFilePath)

	// 拼接输出文件的完整路径
	inputFilePath = filepath.Join(workingDir, inputFilePath)
	thumbnailPath := filepath.Join(workingDir, filepath.Join(dir, outputFileName))

	// 根据文件类型，确定命令
	var cmd *exec.Cmd
	if fileType == consts.Video.Index() {
		cmd = exec.Command("ffmpeg", "-i", inputFilePath, "-vf", "thumbnail="+strconv.Itoa(consts.CompressImageLength)+":-1", "-frames:v", "1", thumbnailPath)
	} else if fileType == consts.Image.Index() {
		cmd = exec.Command("ffmpeg", "-i", inputFilePath, "-vf", "scale="+strconv.Itoa(consts.CompressImageLength)+":"+strconv.Itoa(-1), "-frames:v", "1", thumbnailPath)
	}
	utils.LogrusObj.Infoln("执行缩略图的命令:", cmd)

	// 执行命令
	err = cmd.Run()
	if err != nil {
		utils.LogrusObj.Error("创建视频缩略图时出错:", err)
		return err
	}
	return nil
}

func GenerateThumbnailName(originalName string) string {
	extension := filepath.Ext(originalName)
	baseName := strings.TrimSuffix(originalName, extension)
	thumbnailName := fmt.Sprintf("%s_thumbnail.png", baseName)

	// 检查新名称是否与原始文件名相同，如果相同则添加额外的标识
	if thumbnailName == originalName {
		thumbnailName = fmt.Sprintf("%s_thumbnail_new%s", baseName, extension)
	}

	return thumbnailName
}

func SplitVideo(videoFilePath string) error {
	// 获取视频文件名和目录路径
	videoFileName := GetFileNameWithoutExtension(videoFilePath)
	videoDirPath := filepath.Dir(videoFilePath)

	// 创建同名目录
	err := os.Mkdir(videoDirPath+"/"+videoFileName, 0755)
	if err != nil {
		utils.LogrusObj.Error("创建视频同名目录出错:", err)
		return err
	}

	// 将视频转换为index.ts中间文件
	intermediateFilePath := videoDirPath + "/" + videoFileName + "/index.ts"
	err = convertVideoToTS(videoFilePath, intermediateFilePath)
	if err != nil {
		return err
	}

	// 在目录下生成索引文件 index.m3u8 和切片
	outputDir := videoDirPath + "/" + videoFileName
	err = generateHLSFiles(intermediateFilePath, outputDir)
	if err != nil {
		return err
	}

	// 删除 index.ts 中间文件
	err = os.Remove(intermediateFilePath)
	if err != nil {
		return err
	}

	return nil
}

// 使用 FFmpeg 将视频转换为 index.ts 中间文件
func convertVideoToTS(inputFile string, outputFile string) error {
	workDir, _ := os.Getwd()
	inputFilePath := filepath.Join(workDir, inputFile)
	outputFilePath := filepath.Join(workDir, outputFile)
	cmd := exec.Command("ffmpeg", "-i", inputFilePath, "-c", "copy", "-bsf:v", "h264_mp4toannexb", "-f", "mpegts", outputFilePath)
	err := cmd.Run()
	utils.LogrusObj.Infoln("将视频转换为index.ts的命令:", cmd)
	if err != nil {
		utils.LogrusObj.Infoln("将视频转换为index.ts出错:", err)
		return err
	}
	return nil
}

// 使用 FFmpeg 生成索引文件 index.m3u8 和切片
func generateHLSFiles(inputFile string, outputDir string) error {
	workDir, _ := os.Getwd()
	inputFilePath := filepath.Join(workDir, inputFile)
	outputDir = filepath.Join(workDir, outputDir)
	cmd := exec.Command("ffmpeg", "-i", inputFilePath, "-c:v", "copy", "-c:a", "copy", "-map", "0", "-f", "segment", "-segment_time", "30", "-segment_list", outputDir+"/index.m3u8", "-segment_format", "mpegts", outputDir+"/%d.ts")
	err := cmd.Run()
	utils.LogrusObj.Infoln("生成索引文件 index.m3u8 和切片的命令:", err)
	if err != nil {
		return err
	}
	return nil
}
