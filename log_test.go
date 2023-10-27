package logs

import (
	"testing"
)

func TestLogger_Console(t *testing.T) {
	Log.SetFile("1.txt")
	Log.Init()
	defer Log.Close(false)
	Log.Level = 1
	Log.Important("test")
	Log.Color = true
	Log.Important("test")
	LogNameMap[1] = "detail"
	Log.Log(1, "test")
}
