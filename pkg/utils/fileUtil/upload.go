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
	"strconv"
	"strings"
)

// UploadAvatarToLocalStatic 上传头像到本地
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

// GetFilePath 获取文件路径
func GetFilePath(fileId, userId, fileName string) string {
	// 文件配置
	uConfig := conf.Conf.UploadPath

	// 获取类别key
	categoryKey, flag := FindCategoryKey(fileName)
	if !flag {
		utils.LogrusObj.Infoln("未收录该类别文件")
	}

	// 获取文件夹
	uploadPath := ""
	switch categoryKey {
	case consts.VIDEO.Index():
		uploadPath = uConfig.VideoPath
	case consts.MUSIC.Index():
		uploadPath = uConfig.MusicPath
	case consts.IMAGE.Index():
		uploadPath = uConfig.ImagePath
	case consts.PDF.Index():
		uploadPath = uConfig.DocPath
	case consts.WORD.Index():
		uploadPath = uConfig.DocPath
	case consts.EXCEL.Index():
		uploadPath = uConfig.DocPath
	case consts.TXT.Index():
		uploadPath = uConfig.DocPath
	case consts.PROGRAM.Index():
		uploadPath = uConfig.ProgramPath
	case consts.ZIP.Index():
		uploadPath = uConfig.ZipPath
	default:
		uploadPath = uConfig.OthersPath
	}

	fileName = fileId + "." + GetFileExtension(fileName)

	// 用户文件夹
	basePath := uploadPath + "user_" + userId + "/"
	if !dirIsExist(basePath) {
		CreateDir(basePath)
	}

	// 确定文件路径
	filePath := filepath.Join(basePath, fileName)

	return filePath
}

// MergeChunks 合并分块文件
func MergeChunks(userId, fileId, filePath string) error {
	// 加载配置
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
	chunksDir := filepath.Join(workDir, uConfig.TempPath+"user_"+userId+"/"+fileId+"/")

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
//func GetFileType(fileName string) int {
//	extension := GetFileExtension(fileName)
//	switch strings.ToLower(extension) {
//	case "mp4", "avi", "mkv":
//		return consts.Video.Index()
//	case "jpg", "jpeg", "png", "gif":
//		return consts.Image.Index()
//	case "doc", "docx", "pdf", "txt":
//		return consts.Doc.Index()
//	default:
//		return consts.Others.Index()
//	}
//}

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

func GetTempFilePath(fileId, userId string) string {
	uConfig := conf.Conf.UploadPath

	// 分块文件保存的具体文件夹：分块文件统一保存在临时文件分区
	basePath := uConfig.TempPath + "user_" + userId + "/" + fileId + "/"
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

func FindTypeKey(fileName string) (int, bool) {
	extension := GetFileExtension(fileName)
	for key, extensions := range consts.TypeMapping {
		for _, ext := range extensions {
			if ext == extension {
				return key, true
			}
		}
	}
	return -1, false
}

func FindCategoryKey(fileName string) (int, bool) {
	extension := GetFileExtension(fileName)
	for key, extensions := range consts.TypeMapping {
		for _, ext := range extensions {
			if ext == extension {
				if key >= consts.PDF.Index() && key <= consts.TXT.Index() {
					return consts.DocCategory.Index(), true
				} else if key >= consts.PROGRAM.Index() && key <= consts.OthersCategory.Index() {
					return consts.OthersCategory.Index(), true
				}
				return key, true
			}
		}
	}
	return -1, false
}

func GetFileNameIndex(fileName string) int {
	re := regexp.MustCompile(`\((\d+)\)[^.]*\.`) // 正则表达式匹配括号内的数字，并排除`.`符号之前的内容
	matches := re.FindStringSubmatch(fileName)
	if len(matches) > 1 {
		indexStr := matches[1]
		index, err := strconv.Atoi(indexStr)
		if err == nil {
			return index
		}
	}
	return 0 // 如果没有匹配到索引号，则默认为0
}
