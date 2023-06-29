package test

import (
	"easy-drive/pkg/e"
	"fmt"
	"testing"
)

func TestGetMsg(t *testing.T) {
	fmt.Println(e.GetMsg(e.UserNotRegisterError))
}
