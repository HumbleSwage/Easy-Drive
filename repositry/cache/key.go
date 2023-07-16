package cache

import "fmt"

func VerificationCodeCacheKey(kind int, email string) string {
	return fmt.Sprintf("VerificationCodeCacheKey:%s:%d", email, kind)
}

func UserStoreSpaceKey(userId, fileName string) string {
	return fmt.Sprintf("UserStoreSpaceKey:%s:%s", userId, fileName)
}

func DownloadFileKey(userId, downloadCode string) string {
	return fmt.Sprintf("DownloadFileKey:%s:%s", userId, downloadCode)
}

func ShareFileKey(shareId string) string {
	return fmt.Sprintf("ShareFileKey:%s", shareId)
}
