package easylogger

import (
	"log"
	"os"
	"sync"
)

var (
	//用来判断文件路径的打开状态,防止重复打开文件
	filePathCount map[string]uint
	filePathMap   map[string]*os.File
	loggerMu      sync.Mutex
)

type Logger struct {
	*log.Logger
	openFile *os.File
	logPath  string
	isInit   bool
}

func init() {
	filePathCount = make(map[string]uint)
	filePathMap = make(map[string]*os.File)
}

func NewLogger() *Logger {
	return &Logger{
		Logger:  log.New(os.Stderr, "", log.LstdFlags),
		logPath: "",
		isInit:  false,
	}
}

func (logger *Logger) Close() {
	loggerMu.Lock()
	defer loggerMu.Unlock()

	filePathCount[logger.logPath]--
	if fileOpenCount, ok := filePathCount[logger.logPath]; ok && 0 == fileOpenCount {
		delete(filePathCount, logger.logPath)
		delete(filePathMap, logger.logPath)
		//close file
		logger.openFile.Close()
	}
	logger.isInit = false
}

func (logger *Logger) Open(newLogPath, prefix string) {
	var err error
	if logger.isInit {
		return
	}
	logger.isInit = true
	if "" == newLogPath {
		panic("log file empty!")
		return
	}

	loggerMu.Lock()
	defer loggerMu.Unlock()
	//文件引用计数
	if _, ok := filePathCount[newLogPath]; ok {
		filePathCount[newLogPath]++
		logger.openFile = filePathMap[newLogPath]
	} else {
		logger.openFile, err = os.OpenFile(newLogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic("error opening file")
		}
		filePathCount[newLogPath] = 1
		filePathMap[newLogPath] = logger.openFile
	}
	logger.logPath = newLogPath

	logger.SetFlags(log.LstdFlags)
	logger.SetPrefix(prefix)
	logger.SetOutput(logger.openFile)
}
