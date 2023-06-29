package commonUtil

import (
	"easy-drive/consts"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"math/rand"
	"strconv"
	"time"
)

var node *snowflake.Node

func InitNodeId(startTime string, machineID int64) {
	var st time.Time
	// 格式化 1月2号下午3时4分5秒  2006年
	st, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}

	snowflake.Epoch = st.UnixNano() / 1e6
	node, err = snowflake.NewNode(machineID)
	if err != nil {
		return
	}
	return
}

func GenerateUserID() string {
	id := node.Generate().Int64()
	idStr := strconv.FormatInt(int64(id), consts.UserIdLength)
	return idStr
}

func GenerateFileID() string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, consts.FileIdLength)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GenerateRandomNumber() string {
	length := consts.RandomNumberLength
	rand.Seed(time.Now().UnixNano())
	numberRunes := []rune("0123456789")
	result := make([]rune, length)
	for i := 0; i < length; i++ {
		result[i] = numberRunes[rand.Intn(len(numberRunes))]
	}
	return string(result)
}

func GenerateChunkName() string {
	// 生成 UUID
	uuid1 := uuid.New()

	// 获取当前时间戳
	timestamp := time.Now().UnixNano()

	// 将时间戳转换为字符串
	timestampStr := fmt.Sprintf("%d", timestamp)

	// 构建带有排序信息的文件名
	filename := fmt.Sprintf("chunk_%s_%s", timestampStr, uuid1.String())

	return filename
}
