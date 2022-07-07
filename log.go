package logs

import (
	"fmt"
	. "github.com/chainreactors/files"
	"os"
	"path"
	"time"
)

func NewLogger(quiet, debug bool) *Logger {
	log := &Logger{
		Quiet:   quiet,
		IsDebug: debug,
		Level:   1,
	}

	if debug {
		log.Level = 0
	}

	return log
}

type Logger struct {
	Quiet       bool
	Clean       bool
	IsDebug     bool
	logCh       chan string
	LogFileName string
	logFile     *File
	Level       int
}

const (
	Debug = iota
	Warn
	Info
	Error
	Important
)

var (
	DebugFormatter     = "[debug] %s \n"
	WarnFormatter      = "[warn] %s \n"
	InfoFormatter      = "[+] %s ," + getCurtime() + "\n"
	ErrorFormatter     = "[-] %s ," + getCurtime() + "\n"
	ImportantFormatter = "[*] %s ," + getCurtime() + "\n"
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
	s = fmt.Sprintf(formatter, s)
	if !log.Quiet && level >= log.Level {
		fmt.Print(s)
		if log.logFile != nil {
			log.logFile.SafeWrite(s)
			log.logFile.SafeSync()
		}
	}
}

func (log *Logger) logInterfacef(formatter string, level int, format string, s ...interface{}) {
	line := fmt.Sprintf(fmt.Sprintf(formatter, format), s...)
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
	if log.logFile != nil {
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
