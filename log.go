package logs

import (
	"fmt"
	. "github.com/chainreactors/files"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

var Log *Logger = NewLogger(1, false)

func NewLogger(level Level, quiet bool) *Logger {
	log := &Logger{
		Quiet:  quiet,
		Level:  level,
		Color:  false,
		Writer: os.Stdout,
		SuffixFunc: func() string {
			return ", " + getCurtime()
		},
		PrefixFunc: func() string {
			return ""
		},
	}

	return log
}

type Logger struct {
	Quiet       bool
	Clean       bool
	Color       bool
	logCh       chan string
	LogFileName string
	logFile     *File
	Writer      io.Writer
	Level       Level
	SuffixFunc  func() string
	PrefixFunc  func() string
}

type Level int

const (
	Debug Level = iota
	Warn
	Info
	Error
	Important
)

var DefaultColorMap = map[Level]func(string) string{
	Debug:     Yellow,
	Error:     Red,
	Warn:      Cyan,
	Important: Green,
}

var (
	DebugFormatter     = "[debug] %s "
	WarnFormatter      = "[warn] %s "
	InfoFormatter      = "[+] %s {{suffix}}"
	ErrorFormatter     = "[-] %s {{suffix}}"
	ImportantFormatter = "[*] %s {{suffix}}"
)

func (log *Logger) Init() {
	log.InitFile(log.LogFileName)
}

func (log *Logger) InitFile(filename string) {
	// 初始化进度文件
	var err error
	Log.LogFileName = path.Join(GetExcPath(), filename)
	log.logFile, err = NewFile(Log.LogFileName, false, false, true)
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

func (log *Logger) logInterface(formatter string, level Level, s string) {
	line := fmt.Sprintf(formatter, s)
	line = strings.Replace(line, "{{suffix}}", log.SuffixFunc(), -1)
	line = strings.Replace(line, "{{prefix}}", log.PrefixFunc(), -1)
	line += "\n"
	if !log.Quiet && level >= log.Level {
		if log.Color {
			fmt.Fprint(log.Writer, DefaultColorMap[level](line))
		} else {
			fmt.Fprint(log.Writer, line)
		}

		if log.logFile != nil {
			log.logFile.SafeWrite(line)
			log.logFile.SafeSync()
		}
	}
}

func (log *Logger) logInterfacef(formatter string, level Level, format string, s ...interface{}) {
	line := fmt.Sprintf(fmt.Sprintf(formatter, format), s...)
	line = strings.Replace(line, "{{suffix}}", log.SuffixFunc(), -1)
	line = strings.Replace(line, "{{prefix}}", log.PrefixFunc(), -1)
	line += "\n"
	if !log.Quiet && level >= log.Level {
		if log.Color {
			fmt.Fprint(log.Writer, DefaultColorMap[level](line))
		} else {
			fmt.Fprint(log.Writer, line)
		}

		if log.logFile != nil {
			log.logFile.SafeWrite(line)
			log.logFile.SafeSync()
		}
	}
}

func (log *Logger) Important(s string) {
	log.logInterface(ImportantFormatter, Important, s)
}

func (log *Logger) Importantf(format string, s ...interface{}) {
	log.logInterfacef(ImportantFormatter, Important, format, s...)
}

func (log *Logger) Info(s string) {
	log.logInterface(InfoFormatter, Info, s)
}

func (log *Logger) Infof(format string, s ...interface{}) {
	log.logInterfacef(InfoFormatter, Info, format, s...)
}

func (log *Logger) Error(s string) {
	log.logInterface(ErrorFormatter, Error, s)
}

func (log *Logger) Errorf(format string, s ...interface{}) {
	log.logInterfacef(ErrorFormatter, Error, format, s...)
}

func (log *Logger) Warn(s string) {
	log.logInterface(WarnFormatter, Warn, s)
}

func (log *Logger) Warnf(format string, s ...interface{}) {
	log.logInterfacef(WarnFormatter, Warn, format, s...)
}

func (log *Logger) Debug(s string) {
	log.logInterface(DebugFormatter, Debug, s)

}

func (log *Logger) Debugf(format string, s ...interface{}) {
	log.logInterfacef(DebugFormatter, Debug, format, s...)
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
