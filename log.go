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

var Levels = map[Level]string{
	Debug:     "debug",
	Info:      "info",
	Error:     "error",
	Warn:      "warn",
	Important: "important",
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
		level:     level,
		color:     false,
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
		level:     Warn,
		writer:    file,
		formatter: DefaultFormatterMap,
		levels:    Levels,
	}
	return log, nil
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

	quiet       bool // is enable Print
	clean       bool // is enable Console()
	color       bool
	LogFileName string
	writer      io.Writer
	level       Level
	levels      map[Level]string
	formatter   map[Level]string
	colorMap    map[Level]func(string) string
	SuffixFunc  func() string
	PrefixFunc  func() string
}

func (log *Logger) SetQuiet(q bool) {
	log.quiet = q
}

func (log *Logger) SetClean(c bool) {
	log.clean = c
}

func (log *Logger) SetColor(c bool) {
	log.color = c
}

func (log *Logger) SetColorMap(cm map[Level]func(string) string) {
	log.colorMap = cm
}

func (log *Logger) SetLevel(l Level) {
	log.level = l
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
	if !log.clean {
		fmt.Fprint(log.writer, s)
	}
}

func (log *Logger) Consolef(format string, s ...interface{}) {
	if !log.clean {
		fmt.Fprintf(log.writer, format, s...)
	}
}

func (log *Logger) logInterface(level Level, s interface{}) {
	if !log.quiet && level >= log.level {
		line := log.Format(level, s)
		if log.color {
			fmt.Fprint(log.writer, log.Color(level, line))
		} else {
			fmt.Fprint(log.writer, line)
		}

		if log.logFile != nil {
			log.logFile.SafeWrite(line)
			log.logFile.SafeSync()
		}
	}
}

func (log *Logger) logInterfacef(level Level, format string, s ...interface{}) {
	if !log.quiet && level >= log.level {
		line := log.Format(level, fmt.Sprintf(format, s...))
		if log.color {
			fmt.Fprint(log.writer, log.Color(level, line))
		} else {
			fmt.Fprint(log.writer, line)
		}

		if log.logFile != nil {
			log.logFile.SafeWrite(line)
			log.logFile.SafeSync()
		}
	}
}

func (log *Logger) Log(level Level, s interface{}) {
	log.logInterface(level, s)
}

func (log *Logger) Logf(level Level, format string, s ...interface{}) {
	log.logInterfacef(level, format, s...)
}

func (log *Logger) Important(s interface{}) {
	log.logInterface(Important, s)
}

func (log *Logger) Importantf(format string, s ...interface{}) {
	log.logInterfacef(Important, format, s...)
}

func (log *Logger) Info(s interface{}) {
	log.logInterface(Info, s)
}

func (log *Logger) Infof(format string, s ...interface{}) {
	log.logInterfacef(Info, format, s...)
}

func (log *Logger) Error(s interface{}) {
	log.logInterface(Error, s)
}

func (log *Logger) Errorf(format string, s ...interface{}) {
	log.logInterfacef(Error, format, s...)
}

func (log *Logger) Warn(s interface{}) {
	log.logInterface(Warn, s)
}

func (log *Logger) Warnf(format string, s ...interface{}) {
	log.logInterfacef(Warn, format, s...)
}

func (log *Logger) Debug(s interface{}) {
	log.logInterface(Debug, s)

}

func (log *Logger) Debugf(format string, s ...interface{}) {
	log.logInterfacef(Debug, format, s...)
}

func (log *Logger) Color(level Level, line string) string {
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
	line += "\n"
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

//获取当前时间
func getCurtime() string {
	curtime := time.Now().Format("2006-01-02 15:04.05")
	return curtime
}
