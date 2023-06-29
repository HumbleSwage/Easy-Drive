package fileUtil

import (
	"easy-drive/conf"
	"easy-drive/consts"
	"easy-drive/pkg/utils"
	"easy-drive/pkg/utils/commonUtil"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

func UploadAvatarToLocalStatic(file multipart.File, fileSize int64, userId string) (string, error) {
	// 上传路径拼接
	pConfig := conf.Conf.UploadPath
	basePath := pConfig.AvatarPath + "user_" + userId + "/"
	if !dirIsExist(basePath) {
		CreateDir(basePath)
	}
	avatarPath := basePath + "avatar.png"

	// 检查文件
	_, err := os.Stat(avatarPath)
	if os.IsNotExist(err) {
		if _, err = os.Create(avatarPath); err != nil {
			return "", err
		}
	}

	// 覆盖文件
	f, err := os.Create(avatarPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// 读取内容
	content := make([]byte, fileSize)
	if _, err = file.Read(content); err != nil {
		return "", err
	}

	// 写入文件
	if _, err := f.Write(content); err != nil {
		return "", err
	}

	return "user_" + userId + "/" + "avatar.png", err

}

func ChunkedFileToLocalTemp(userId, fileName string, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	uConfig := conf.Conf.UploadPath

	// 分块文件保存的具体文件夹：分块文件统一保存在临时文件分区
	basePath := uConfig.TempPath + "user_" + userId + "/" + fileName + "/"
	if !dirIsExist(basePath) {
		CreateDir(basePath)
	}

	// 分块文件名
	chunkPath := basePath + commonUtil.GenerateChunkName()

	// 检查文件
	_, err := os.Stat(chunkPath)
	if os.IsNotExist(err) {
		if _, err = os.Create(chunkPath); err != nil {
			return "", err
		}
	}

	// 覆盖
	f, err := os.Create(chunkPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// 读取文件内容
	content := make([]byte, fileHeader.Size)
	if _, err = file.Read(content); err != nil {
		return "", err
	}

	// 写入文件
	if _, err := f.Write(content); err != nil {
		return "", err
	}

	return "user_" + "/" + fileName + "/", err
}

func GetUniqueFileName(fileName, userId string) (string, string) {
	// 文件配置
	uConfig := conf.Conf.UploadPath
	extension := GetFileExtension(fileName)

	// 获取文件夹
	uploadPath := ""
	switch GetFileType(fileName) {
	case consts.Video.Index():
		uploadPath = uConfig.VideoPath
	case consts.Image.Index():
		uploadPath = uConfig.ImagePath
	case consts.Doc.Index():
		uploadPath = uConfig.DocPath
	default:
		uploadPath = uConfig.OthersPath
	}

	// 文件夹是否存在
	basePath := uploadPath + "user_" + userId
	if !dirIsExist(basePath) {
		CreateDir(basePath)
		return basePath + "/" + fileName, fileName
	}

	// 无后缀文件名
	fileNameWithoutExtension := strings.TrimSuffix(fileName, "."+extension)

	// 新文件名
	newFileName := fileName
	counter := 1
	for fileExists(newFileName, basePath) {
		newFileName = fmt.Sprintf("%s(%d).%s", fileNameWithoutExtension, counter, extension)
		counter++
	}

	// 文件地址
	newFilePath := basePath + "/" + newFileName

	return newFilePath, newFileName
}

func MergeChunks(userId, fileName, filePath string) error {
	uConfig := conf.Conf.UploadPath
	// 加上工作目录
	workDir, _ := os.Getwd()
	filePath = filepath.Join(workDir, filePath)
	utils.LogrusObj.Infoln(filePath)

	// 打开文件
	targetFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		utils.LogrusObj.Infoln("打开文件出错：", err)
		return err
	}
	defer targetFile.Close()

	// 分块保存分区
	chunksDir := filepath.Join(workDir, uConfig.TempPath+"user_"+userId+"/"+fileName+"/")
	utils.LogrusObj.Infoln(chunksDir)
	// 获取对应所有分块
	chunkFiles, err := getChunkFiles(chunksDir)
	if err != nil {
		utils.LogrusObj.Infoln("获取所有临时文件出错:", err)
		return err
	}

	// 如果没有temp
	if chunkFiles == nil {
		utils.LogrusObj.Infoln("临时文件为空:", err)
		return err
	}

	// 按照时间戳对分块文件进行排序
	sort.Slice(chunkFiles, func(i, j int) bool {
		timestamp1 := getTimestampFromFilename(chunkFiles[i])
		timestamp2 := getTimestampFromFilename(chunkFiles[j])
		return timestamp1 < timestamp2
	})

	// 依次读取分块文件并写入目标文件
	for _, chunkFileName := range chunkFiles {
		// 分块文件地址
		chunkFilePath := filepath.Join(chunksDir, chunkFileName)

		// 打开文件
		chunkFile, err := os.Open(chunkFilePath)
		if err != nil {
			utils.LogrusObj.Infoln("打开分块出错:", err)
			return err
		}

		// 拷贝文件
		if _, err = io.Copy(targetFile, chunkFile); err != nil {
			chunkFile.Close()
			utils.LogrusObj.Infoln("拷贝文件到指定文件中出错:", err)
			return err
		}

		// 手动关闭分块文件
		if err = chunkFile.Close(); err != nil {
			return err
		}

		// 删除分块文件
		if err = os.Remove(chunkFilePath); err != nil {
			return err
		}
	}

	// 删除临时区
	if err = os.Remove(chunksDir); err != nil {
		utils.LogrusObj.Infoln("删除文件临时分区出错:", err)
		return err
	}

	return err
}

func GetFileTypeNumber(fileName string) (int, error) {
	extension := GetFileExtension(fileName)
	// 将扩展名转换为小写，以便进行不区分大小写的比较
	extension = strings.ToLower(extension)

	// 查找扩展名在 TypeStr 切片中的索引
	for i, fileTypeStr := range consts.TypeStr {
		if fileTypeStr == extension {
			return i, nil
		}
	}

	return -1, fmt.Errorf("未找到文件类型数字：%s", extension)
}

// 检查目录是否存在
func dirIsExist(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// CreateDir 创建目录
func CreateDir(path string) bool {
	err := os.MkdirAll(path, 0777)
	if err != nil {
		return false
	}
	return true
}

// GetFileExtension 获取文件后缀名
func GetFileExtension(fileName string) string {
	parts := strings.Split(fileName, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return ""
}

// GetFileType 根据文件后缀名确定文件类型
func GetFileType(fileName string) int {
	extension := GetFileExtension(fileName)
	switch strings.ToLower(extension) {
	case "mp4", "avi", "mkv":
		return consts.Video.Index()
	case "jpg", "jpeg", "png", "gif":
		return consts.Image.Index()
	case "doc", "docx", "pdf", "txt":
		return consts.Doc.Index()
	default:
		return consts.Others.Index()
	}
}

// 检查文件是否存在
func fileExists(fileName, directory string) bool {
	_, err := os.Stat(filepath.Join(directory, fileName))
	return !os.IsNotExist(err)
}

// 获取符合格式的分块文件列表
func getChunkFiles(directory string) ([]string, error) {
	var chunkFiles []string

	// 给予权限
	err := os.Chmod(directory, 0755) // 设置为 rwxr-xr-x 权限
	if err != nil {
		// 处理错误
		utils.LogrusObj.Infoln("无法更改权限")
		return nil, err
	}
	// 读取目录中的文件
	files, err := os.ReadDir(directory)
	if err != nil {
		utils.LogrusObj.Infoln("获取所有文件出错")
		return nil, err
	}

	// 正则表达式匹配文件名格式
	pattern := `^chunk_\d+_.+$`
	regExp, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	// 筛选符合格式的文件
	for _, file := range files {
		if file.Type().IsRegular() && regExp.MatchString(file.Name()) {
			chunkFiles = append(chunkFiles, file.Name())
		}
	}

	return chunkFiles, nil
}

// 从文件名中提取时间戳
func getTimestampFromFilename(filename string) int64 {
	parts := strings.Split(filename, "_")
	if len(parts) != 3 {
		return 0
	}
	var timestamp int64
	_, err := fmt.Sscanf(parts[1], "%d", &timestamp)
	if err != nil {
		return 0
	}
	return timestamp
}

func GetTempFilePath(fileName, userId string) string {
	uConfig := conf.Conf.UploadPath

	// 分块文件保存的具体文件夹：分块文件统一保存在临时文件分区
	basePath := uConfig.TempPath + "user_" + userId + "/" + fileName + "/"
	if !dirIsExist(basePath) {
		CreateDir(basePath)
	}

	// 分块文件名
	chunkPath := basePath + commonUtil.GenerateChunkName()

	return chunkPath
}

// GetFileNameWithoutExtension 从文件路径中提取不带扩展名的文件名
func GetFileNameWithoutExtension(filePath string) string {
	fileNameWithExtension := filepath.Base(filePath)
	fileName := strings.TrimSuffix(fileNameWithExtension, filepath.Ext(fileNameWithExtension))
	return fileName
}
