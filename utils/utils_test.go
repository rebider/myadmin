package utils

import (
	"testing"
)

func TestRandomString(t *testing.T) {
	str := RandomString(10)
	InitLogs()
	//Debug("asdfasdfasdf")
	t.Error(str)
}
