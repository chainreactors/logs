package logs

import (
	"fmt"
	. "github.com/chainreactors/files"
	"os"
	"path"
	"time"
)

func NewLogger(quiet, debug bool) *Logger {
	return &Logger{
		Quiet:   quiet,
		IsDebug: debug,
	}
}

type Logger struct {
	Quiet       bool
	Clean       bool
	IsDebug     bool
	logCh       chan string
	LogFileName string
	logFile     *File
}

func (log *Logger) Init() {
	if log.LogFileName != "" {
		log.initFile()
		go func() {
			for res := range log.logCh {
				log.logFile.SyncWrite(res)
			}
			log.logFile.Close()
		}()
	}
}

func (log *Logger) initFile() {
	// 初始化进度文件
	var err error
	log.LogFileName = path.Join(GetExcPath(), log.LogFileName)
	log.logFile, err = NewFile(log.LogFileName, false, false, true)
	if err != nil {
		log.Warn("cannot create logfile, err:" + err.Error())
		return
	}
	log.logCh = make(chan string, 100)
}

func (log *Logger) Important(s string) {
	s = fmt.Sprintf("[*] %s , %s\n", s, getCurtime())
	if !log.Quiet {
		fmt.Print(s)
	}
	if log.logFile != nil {
		log.logCh <- s
	}
}

func (log *Logger) Importantf(format string, s ...interface{}) {
	line := fmt.Sprintf("[*] "+format+", "+getCurtime()+"\n", s...)
	if !log.Quiet {
		fmt.Print(line)
	}
	if log.logFile != nil {
		log.logCh <- line
	}
}

func (log *Logger) Default(s string) {
	if !log.Clean {
		fmt.Print(s)
	}
}

func (log *Logger) Defaultf(format string, s ...interface{}) {
	if !log.Clean {
		fmt.Printf(format, s...)
	}
}

func (log *Logger) Error(s string) {
	if !log.Quiet {
		fmt.Println("[-] " + s)
	}
}

func (log *Logger) Errorf(format string, s ...interface{}) {
	if !log.Quiet {
		fmt.Printf("[-] "+format+"\n", s...)
	}
}

func (log *Logger) Warn(s string) {
	if !log.Quiet {
		fmt.Println("[warn] " + s)
	}
}

func (log *Logger) Warnf(format string, s ...interface{}) {
	if !log.Quiet {
		fmt.Printf("[warn] "+format+"\n", s...)
	}
}

func (log *Logger) Debug(s string) {
	if log.IsDebug {
		fmt.Println("[debug] " + s)
	}
}

func (log *Logger) Debugf(format string, s ...interface{}) {
	if log.IsDebug {
		fmt.Printf("[debug] "+format+"\n", s...)
	}
}

func (log *Logger) Close(remove bool) {
	if log.logCh != nil {
		close(log.logCh)
		time.Sleep(time.Microsecond * 200)
	}

	if remove {
		err := os.Remove(log.LogFileName)
		if err != nil {
			log.Warn(err.Error())
		}
	}
}

func IsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

//获取当前时间
func getCurtime() string {
	curtime := time.Now().Format("2006-01-02 15:04.05")
	return curtime
}
