package logs

import (
	"testing"
)

func TestLogger_Console(t *testing.T) {
	Log.SetFile("1.txt")
	Log.Init()
	defer Log.Close(false)
	Log.SetLevel(1)
	Log.Important("test")
	Log.SetColor(true)
	Log.Importantf("%stest", "aaa")
	AddLevel(1, "test")
	Log.Log(1, "test")
}
