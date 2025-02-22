package logs

import (
	"fmt"
	"github.com/chainreactors/files"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

var Log *Logger = NewLogger(WarnLevel)

var defaultColor = func(s string) string { return s }
var DefaultColorMap = map[Level]func(string) string{
	DebugLevel:     Yellow,
	ErrorLevel:     RedBold,
	InfoLevel:      Cyan,
	WarnLevel:      YellowBold,
	ImportantLevel: PurpleBold,
}

var DefaultFormatterMap = map[Level]string{
	DebugLevel:     "[debug] %s \n",
	WarnLevel:      "[warn] %s \n",
	InfoLevel:      "[+] %s {{suffix}}\n",
	ErrorLevel:     "[-] %s {{suffix}}\n",
	ImportantLevel: "[*] %s {{suffix}}\n",
}

var Levels = map[Level]string{
	DebugLevel:     "debug",
	InfoLevel:      "info",
	ErrorLevel:     "error",
	WarnLevel:      "warn",
	ImportantLevel: "important",
}

func AddLevel(level Level, name string, opts ...interface{}) {
	Levels[level] = name
	for _, opt := range opts {
		switch opt.(type) {
		case string:
			DefaultFormatterMap[level] = opt.(string)
		case func(string) string:
			DefaultColorMap[level] = opt.(func(string) string)
		}
	}
}

func NewLogger(level Level) *Logger {
	log := &Logger{
		Level:     level,
		Color:     false,
		writer:    os.Stdout,
		levels:    Levels,
		formatter: DefaultFormatterMap,
		colorMap:  DefaultColorMap,
		SuffixFunc: func() string {
			return ", " + getCurtime()
		},
		PrefixFunc: func() string {
			return ""
		},
	}

	return log
}

// NewFileLogger create a pure file logger
func NewFileLogger(filename string) (*Logger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	log := &Logger{
		Level:     WarnLevel,
		writer:    file,
		formatter: DefaultFormatterMap,
		levels:    Levels,
	}
	return log, nil
}

const (
	DebugLevel     Level = 10
	WarnLevel      Level = 20
	InfoLevel      Level = 30
	ErrorLevel     Level = 40
	ImportantLevel Level = 50
)

type Level int

func (l Level) Name() string {
	if name, ok := Levels[l]; ok {
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

	Quiet       bool // is enable Print
	Clean       bool // is enable Console()
	Color       bool
	LogFileName string
	writer      io.Writer
	Level       Level
	levels      map[Level]string
	formatter   map[Level]string
	colorMap    map[Level]func(string) string
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
	log.colorMap = cm
}

func (log *Logger) SetLevel(l Level) {
	log.Level = l
}

func (log *Logger) SetOutput(w io.Writer) {
	log.writer = w
}

func (log *Logger) SetFile(filename string) {
	log.LogFileName = filename
}

func (log *Logger) SetFormatter(formatter map[Level]string) {
	log.formatter = formatter
}

func (log *Logger) Init() {
	// 初始化进度文件
	var err error
	log.logFile, err = files.NewFile(log.LogFileName, false, false, true)
	if err != nil {
		log.Warn("cannot create logfile, err:" + err.Error())
		return
	}
	log.logCh = make(chan string, 100)
}

func (log *Logger) Console(s string) {
	if !log.Clean {
		fmt.Fprint(log.writer, s)
	}
}

func (log *Logger) Consolef(format string, s ...interface{}) {
	if !log.Clean {
		fmt.Fprintf(log.writer, format, s...)
	}
}

func (log *Logger) FConsolef(writer io.Writer, format string, s ...interface{}) {
	if !log.Clean {
		fmt.Fprintf(writer, format, s...)
	}
}

func (log *Logger) logInterface(writer io.Writer, level Level, s interface{}) {
	if !log.Quiet && level >= log.Level {
		line := log.Format(level, s)
		if log.Color {
			fmt.Fprint(writer, log.SetLevelColor(level, line))
		} else {
			fmt.Fprint(writer, line)
		}

		if log.logFile != nil {
			log.logFile.SafeWrite(line)
			log.logFile.SafeSync()
		}
	}
}

func (log *Logger) logInterfacef(writer io.Writer, level Level, format string, s ...interface{}) {
	if !log.Quiet && level >= log.Level {
		line := log.Format(level, fmt.Sprintf(format, s...))
		if log.Color {
			fmt.Fprint(writer, log.SetLevelColor(level, line))
		} else {
			fmt.Fprint(writer, line)
		}

		if log.logFile != nil {
			log.logFile.SafeWrite(line)
			log.logFile.SafeSync()
		}
	}
}

func (log *Logger) Log(level Level, s interface{}) {
	log.logInterface(log.writer, level, s)
}

func (log *Logger) Logf(level Level, format string, s ...interface{}) {
	log.logInterfacef(log.writer, level, format, s...)
}

func (log *Logger) FLogf(writer io.Writer, level Level, s ...interface{}) {
	log.logInterface(writer, level, fmt.Sprintln(s...))
}

func (log *Logger) Important(s interface{}) {
	log.logInterface(log.writer, ImportantLevel, s)
}

func (log *Logger) Importantf(format string, s ...interface{}) {
	log.logInterfacef(log.writer, ImportantLevel, format, s...)
}

func (log *Logger) FImportantf(writer io.Writer, format string, s ...interface{}) {
	log.logInterfacef(writer, ImportantLevel, format, s...)
}

func (log *Logger) Info(s interface{}) {
	log.logInterface(log.writer, InfoLevel, s)
}

func (log *Logger) Infof(format string, s ...interface{}) {
	log.logInterfacef(log.writer, InfoLevel, format, s...)
}

func (log *Logger) FInfof(writer io.Writer, format string, s ...interface{}) {
	log.logInterfacef(writer, InfoLevel, format, s...)
}

func (log *Logger) Error(s interface{}) {
	log.logInterface(log.writer, ErrorLevel, s)
}

func (log *Logger) Errorf(format string, s ...interface{}) {
	log.logInterfacef(log.writer, ErrorLevel, format, s...)
}

func (log *Logger) FErrorf(writer io.Writer, format string, s ...interface{}) {
	log.logInterfacef(writer, ErrorLevel, format, s...)
}

func (log *Logger) Warn(s interface{}) {
	log.logInterface(log.writer, WarnLevel, s)
}

func (log *Logger) Warnf(format string, s ...interface{}) {
	log.logInterfacef(log.writer, WarnLevel, format, s...)
}

func (log *Logger) FWarnf(writer io.Writer, format string, s ...interface{}) {
	log.logInterfacef(writer, WarnLevel, format, s...)
}

func (log *Logger) Debug(s interface{}) {
	log.logInterface(log.writer, DebugLevel, s)

}

func (log *Logger) Debugf(format string, s ...interface{}) {
	log.logInterfacef(log.writer, DebugLevel, format, s...)
}

func (log *Logger) FDebugf(writer io.Writer, format string, s ...interface{}) {
	log.logInterfacef(writer, DebugLevel, format, s...)
}

func (log *Logger) SetLevelColor(level Level, line string) string {
	if c, ok := log.colorMap[level]; ok {
		return c(line)
	} else if c, ok := DefaultColorMap[level]; ok {
		return c(line)
	} else {
		return line
	}
}

func (log *Logger) Format(level Level, s ...interface{}) string {
	var line string
	if f, ok := log.formatter[level]; ok {
		line = fmt.Sprintf(f, s...)
	} else if f, ok := DefaultFormatterMap[level]; ok {
		line = fmt.Sprintf(f, s...)
	} else {
		line = fmt.Sprintf("[%s] %s ", append([]interface{}{level.Name()}, s...)...)
	}
	line = strings.Replace(line, "{{suffix}}", log.SuffixFunc(), -1)
	line = strings.Replace(line, "{{prefix}}", log.PrefixFunc(), -1)
	return line
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

func Debug(s interface{}) {
	Log.Debug(s)
}

func Debugf(format string, s ...interface{}) {
	Log.Debugf(format, s...)
}

func Info(s interface{}) {
	Log.Info(s)
}

func Infof(format string, s ...interface{}) {
	Log.Infof(format, s...)
}

func Error(s interface{}) {
	Log.Error(s)
}

func Errorf(format string, s ...interface{}) {
	Log.Errorf(format, s...)
}

func Warn(s interface{}) {
	Log.Warn(s)
}

func Warnf(format string, s ...interface{}) {
	Log.Warnf(format, s...)
}

func Important(s interface{}) {
	Log.Important(s)
}

func Importantf(format string, s ...interface{}) {
	Log.Importantf(format, s...)
}

func Console(s string) {
	Log.Console(s)
}

func Consolef(format string, s ...interface{}) {
	Log.Consolef(format, s...)
}

// 获取当前时间
func getCurtime() string {
	curtime := time.Now().Format("2006-01-02 15:04.05")
	return curtime
}
