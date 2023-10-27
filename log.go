package logs

import (
	"fmt"
	"github.com/chainreactors/files"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var Log *Logger = NewLogger(Warn)

var defaultColor = func(s string) string { return s }
var DefaultColorMap = map[Level]func(string) string{
	Debug:     Yellow,
	Error:     RedBold,
	Info:      Cyan,
	Warn:      YellowBold,
	Important: PurpleBold,
}

var DefaultFormatterMap = map[Level]string{
	Debug:     "[debug] %s ",
	Warn:      "[warn] %s ",
	Info:      "[+] %s {{suffix}}",
	Error:     "[-] %s {{suffix}}",
	Important: "[*] %s {{suffix}}",
}

var LogNameMap = map[Level]string{
	Debug:     "debug",
	Info:      "info",
	Error:     "error",
	Warn:      "warn",
	Important: "important",
}

func NewLogger(level Level) *Logger {
	log := &Logger{
		Level:     level,
		Color:     false,
		Writer:    os.Stdout,
		Formatter: DefaultFormatterMap,
		ColorMap:  DefaultColorMap,
		SuffixFunc: func() string {
			return ", " + getCurtime()
		},
		PrefixFunc: func() string {
			return ""
		},
	}

	return log
}

const (
	Debug     Level = 10
	Warn      Level = 20
	Info      Level = 30
	Error     Level = 40
	Important Level = 50
)

type Level int

func (l Level) Name() string {
	if name, ok := LogNameMap[l]; ok {
		return name
	} else {
		return strconv.Itoa(int(l))
	}
}

func (l Level) Formatter() string {
	if formatter, ok := DefaultFormatterMap[l]; ok {
		return formatter
	} else {
		return "[" + l.Name() + "] %s"
	}
}

func (l Level) Color() func(string) string {
	if f, ok := DefaultColorMap[l]; ok {
		return f
	} else {
		return defaultColor
	}
}

type Logger struct {
	logCh   chan string
	logFile *files.File

	Quiet       bool
	Clean       bool
	Color       bool
	LogFileName string
	Writer      io.Writer
	Level       Level
	Formatter   map[Level]string
	ColorMap    map[Level]func(string) string
	SuffixFunc  func() string
	PrefixFunc  func() string
}

func (log *Logger) SetQuiet(q bool) {
	log.Quiet = q
}

func (log *Logger) SetClean(c bool) {
	log.Clean = c
}

func (log *Logger) SetColor(c bool) {
	log.Color = c
}

func (log *Logger) SetColorMap(cm map[Level]func(string) string) {
	log.ColorMap = cm
}

func (log *Logger) SetLevel(l Level) {
	log.Level = l
}

func (log *Logger) SetOutput(w io.Writer) {
	log.Writer = w
}

func (log *Logger) SetFile(filename string) {
	log.LogFileName = path.Join(files.GetExcPath(), filename)
}

func (log *Logger) SetFormatter(formatter map[Level]string) {
	log.Formatter = formatter
}

func (log *Logger) NewLevel(l int, name string, opt map[string]interface{}) {
	level := Level(l)
	LogNameMap[level] = name
	if opt != nil {
		if f, ok := opt["formatter"]; ok {
			log.Formatter[level] = f.(string)
		} else {
			log.Formatter[level] = "[" + name + "] %s"
		}

		if c, ok := opt["color"]; ok {
			log.ColorMap[level] = c.(func(string) string)
		}
	}
}

func (log *Logger) Init() {
	// 初始化进度文件
	var err error
	log.logFile, err = files.NewFile(Log.LogFileName, false, false, true)
	if err != nil {
		log.Warn("cannot create logfile, err:" + err.Error())
		return
	}
	log.logCh = make(chan string, 100)
}

func (log *Logger) Console(s string) {
	if !log.Clean {
		fmt.Fprint(log.Writer, s)
	}
}

func (log *Logger) Consolef(format string, s ...interface{}) {
	if !log.Clean {
		fmt.Fprintf(log.Writer, format, s...)
	}
}

func (log *Logger) logInterface(level Level, s string) {
	if !log.Quiet && level >= log.Level {
		line := fmt.Sprintf(log.Formatter[level], s)
		line = strings.Replace(line, "{{suffix}}", log.SuffixFunc(), -1)
		line = strings.Replace(line, "{{prefix}}", log.PrefixFunc(), -1)
		line += "\n"
		if log.Color {
			fmt.Fprint(log.Writer, log.ColorMap[level](line))
		} else {
			fmt.Fprint(log.Writer, line)
		}

		if log.logFile != nil {
			log.logFile.SafeWrite(line)
			log.logFile.SafeSync()
		}
	}
}

func (log *Logger) logInterfacef(level Level, format string, s ...interface{}) {
	if !log.Quiet && level >= log.Level {
		line := fmt.Sprintf(fmt.Sprintf(log.Formatter[level], format), s...)
		line = strings.Replace(line, "{{suffix}}", log.SuffixFunc(), -1)
		line = strings.Replace(line, "{{prefix}}", log.PrefixFunc(), -1)
		line += "\n"
		if log.Color {
			fmt.Fprint(log.Writer, log.ColorMap[level](line))
		} else {
			fmt.Fprint(log.Writer, line)
		}

		if log.logFile != nil {
			log.logFile.SafeWrite(line)
			log.logFile.SafeSync()
		}
	}
}

func (log *Logger) Log(level Level, s string) {
	log.logInterface(level, s)
}

func (log *Logger) Logf(level Level, format string, s ...interface{}) {
	log.logInterfacef(level, format, s...)
}

func (log *Logger) Important(s string) {
	log.logInterface(Important, s)
}

func (log *Logger) Importantf(format string, s ...interface{}) {
	log.logInterfacef(Important, format, s...)
}

func (log *Logger) Info(s string) {
	log.logInterface(Info, s)
}

func (log *Logger) Infof(format string, s ...interface{}) {
	log.logInterfacef(Info, format, s...)
}

func (log *Logger) Error(s string) {
	log.logInterface(Error, s)
}

func (log *Logger) Errorf(format string, s ...interface{}) {
	log.logInterfacef(Error, format, s...)
}

func (log *Logger) Warn(s string) {
	log.logInterface(Warn, s)
}

func (log *Logger) Warnf(format string, s ...interface{}) {
	log.logInterfacef(Warn, format, s...)
}

func (log *Logger) Debug(s string) {
	log.logInterface(Debug, s)

}

func (log *Logger) Debugf(format string, s ...interface{}) {
	log.logInterfacef(Debug, format, s...)
}

func (log *Logger) Close(remove bool) {
	if log.logFile != nil && log.logFile.InitSuccess {
		log.logFile.Close()
	}

	if remove {
		err := os.Remove(log.LogFileName)
		if err != nil {
			log.Warn(err.Error())
		}
	}
}

//获取当前时间
func getCurtime() string {
	curtime := time.Now().Format("2006-01-02 15:04.05")
	return curtime
}
