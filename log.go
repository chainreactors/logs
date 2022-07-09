package logs

import (
	"fmt"
	. "github.com/chainreactors/files"
	"os"
	"path"
	"strings"
	"time"
)

var Log *Logger = NewLogger(false, false)

func NewLogger(quiet, debug bool) *Logger {
	log := &Logger{
		Quiet: quiet,
		Level: 1,
		SuffixFunc: func() string {
			return ", " + getCurtime()
		},
	}

	if debug {
		log.Level = 0
	}

	return log
}

type Logger struct {
	Quiet       bool
	Clean       bool
	logCh       chan string
	LogFileName string
	logFile     *File
	Level       int
	SuffixFunc  func() string
	PrefixFunc  func() string
}

const (
	Debug = iota
	Warn
	Info
	Error
	Important
)

var (
	DebugFormatter     = "[debug] %s "
	WarnFormatter      = "[warn] %s "
	InfoFormatter      = "[+] %s {{suffix}}"
	ErrorFormatter     = "[-] %s {{suffix}}"
	ImportantFormatter = "[*] %s {{suffix}}"
)

func (log *Logger) Init() {
	log.initFile()
}

func (log *Logger) initFile() {
	// 初始化进度文件
	if log.LogFileName == "" {
		return
	}
	var err error
	log.LogFileName = path.Join(GetExcPath(), log.LogFileName)
	log.logFile, err = NewFile(log.LogFileName, false, false, true)
	if err != nil {
		log.Warn("cannot create logfile, err:" + err.Error())
		return
	}
	log.logCh = make(chan string, 100)
}

func (log *Logger) Console(s string) {
	if !log.Clean {
		fmt.Print(s)
	}
}

func (log *Logger) Consolef(format string, s ...interface{}) {
	if !log.Clean {
		fmt.Printf(format, s...)
	}
}

func (log *Logger) logInterface(formatter string, level int, s string) {
	line := fmt.Sprintf(formatter, s)
	if len(line) >= 9 && strings.HasSuffix(line, "{{suffix}}") {
		line = line[:len(line)-10] + log.SuffixFunc()
	}
	if len(line) >= 9 && strings.HasPrefix(line, "{{prefix}}") {
		line = log.PrefixFunc() + line[10:]
	}
	line += "\n"
	if !log.Quiet && level >= log.Level {
		fmt.Print(line)
		if log.logFile != nil {
			log.logFile.SafeWrite(line)
			log.logFile.SafeSync()
		}
	}
}

func (log *Logger) logInterfacef(formatter string, level int, format string, s ...interface{}) {
	line := fmt.Sprintf(fmt.Sprintf(formatter, format), s...)
	if len(line) >= 9 && strings.HasSuffix(line, "{{suffix}}") {
		line = line[:len(line)-10] + log.SuffixFunc()
	}
	if len(line) >= 9 && strings.HasPrefix(line, "{{prefix}}") {
		line = log.PrefixFunc() + line[10:]
	}
	line += "\n"
	if !log.Quiet && level >= log.Level {
		fmt.Print(line)
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
	if log.logFile != nil && log.logFile.FileHandler != nil {
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
