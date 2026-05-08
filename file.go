package logs

import (
	"bufio"
	"os"
	"sync"
)

type logFile struct {
	Filename    string
	InitSuccess bool

	file   *os.File
	writer *bufio.Writer
	mu     sync.Mutex
	closed bool
}

func newLogFile(filename string) (*logFile, error) {
	f := &logFile{Filename: filename}
	if err := f.init(); err != nil {
		return nil, err
	}
	return f, nil
}

func (f *logFile) init() error {
	if f.file != nil {
		return nil
	}

	file, err := os.OpenFile(f.Filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	f.file = file
	f.writer = bufio.NewWriter(file)
	f.InitSuccess = true
	return nil
}

func (f *logFile) SafeWrite(s string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.closed {
		return
	}
	_, _ = f.writer.WriteString(s)
}

func (f *logFile) SafeSync() {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.closed {
		return
	}
	_ = f.writer.Flush()
}

func (f *logFile) Close() {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.closed {
		return
	}
	_ = f.writer.Flush()
	_ = f.file.Close()
	f.closed = true
}
