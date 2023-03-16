package logs

import (
	"testing"
)

func TestLogger_Console(t *testing.T) {
	Log.InitFile("1.txt")
	defer Log.Close(false)
	Log.Level = 1
	Log.Important("test")
	Log.Color = true
	Log.Important("test")
	DefaultNameMap[1] = "detail"
	Log.Log(1, "test")
}
