package miniutils

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type LogLevel int

const (
	LOG_LEVEL_DEBUG = 0
	LOG_LEVEL_INFO  = 1
	LOG_LEVEL_WARN  = 2
	LOG_LEVEL_ERROR = 3
)

var levelMap = map[LogLevel]string{
	LOG_LEVEL_DEBUG: "debug",
	LOG_LEVEL_INFO:  "info",
	LOG_LEVEL_WARN:  "warn",
	LOG_LEVEL_ERROR: "error",
}

type singleLogger struct {
	logFile  *os.File
	logger   *log.Logger
	isopened bool
}

type Logger struct {
	lock      *sync.Mutex
	logLevel  LogLevel
	logsDir   string
	loggerMap map[LogLevel]singleLogger
}

func getInitLoggerMap() map[LogLevel]singleLogger {
	return make(map[LogLevel]singleLogger, len(levelMap))
}

// TODO TO BE a global unique instance
func NewLogger(logsDir string) *Logger {
	return &Logger{lock: &sync.Mutex{}, logsDir: logsDir, loggerMap: getInitLoggerMap()}
}

var slogger *Logger

// GetLogger 获取日志工具
// logger:= GetLogger(""); logger.Debug("写一条日志111")
func GetLogger(dirpath string) *Logger {
	if slogger == nil {
		slogger = NewLogger(dirpath)
	}
	return slogger
}

func (l *Logger) SetLogLevel(level LogLevel) *Logger {
	l.logLevel = level
	return l
}

func (l *Logger) getLogger(level LogLevel) *log.Logger {
	l.lock.Lock()
	defer l.lock.Lock()
	lg, ok := l.loggerMap[level]
	var logFile *os.File
	dateExpired := false
	if ok {
		filename := time.Now().Format(fmt.Sprintf("20060102_%s.log", levelMap[level]))
		logFile = l.loggerMap[level].logFile
		if logFile != nil {
			info, err := logFile.Stat()
			if info != nil {
				if info.Name() != filename {
					logFile.Close()
					dateExpired = true
				}
			} else {
				log.Println("---getLogger--[logFile.Stat() == nil]-----err=", err)
			}
		} else {
			log.Println("---getLogger--[l.loggerMap[level].logFile == nil]---")
		}
	}
	if !ok || dateExpired {
		// log.Printf("\n---getLogger--[if !ok || dateExpired]-ok(%t)---dateExpired(%t)--\n", ok, dateExpired)
		logFile = createPath(l.logsDir, level)
		if logFile == nil {
			log.Println("---getLogger--createPath-return-nil--")
		}
		lg = singleLogger{logFile: logFile, logger: log.New(logFile, "\r\n", log.Ldate|log.Ltime), isopened: true}
		l.loggerMap[level] = lg
	}
	return lg.logger
}

func createPath(dirpath string, level LogLevel) *os.File {
	// timeText := time.Now().Format("当前时间: 2006-01-02 15:04:05")
	nowTime := time.Now()
	logsDir := "runtime/logs"
	if dirpath != "" {
		logsDir = dirpath
	}
	dir := logsDir + nowTime.Format("/200601")
	Mkdir(dir)
	filename := nowTime.Format(fmt.Sprintf("20060102_%s.log", levelMap[level]))
	filepath := dir + "/" + filename
	logfile, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return logfile
}

func (l *Logger) Debugf(format string, v ...any) (isPrint bool) {
	l.getLogger(LOG_LEVEL_DEBUG).Printf(format, v...)
	if l.logLevel >= LOG_LEVEL_DEBUG {
		log.Printf(format, v...)
		return true
	}
	return false
}

func (l *Logger) Infof(format string, v ...any) (isPrint bool) {
	p := l.Debugf(format, v...)
	l.getLogger(LOG_LEVEL_INFO).Printf(format, v...)
	if !p && l.logLevel >= LOG_LEVEL_INFO {
		log.Printf(format, v...)
		return true
	}
	return false
}
func (l *Logger) Warnf(format string, v ...any) (isPrint bool) {
	p := l.Infof(format, v...)
	l.getLogger(LOG_LEVEL_WARN).Printf(format, v...)
	if !p && l.logLevel >= LOG_LEVEL_WARN {
		log.Printf(format, v...)
		return true
	}
	return false
}
func (l *Logger) Errorf(format string, v ...any) (isPrint bool) {
	p := l.Warnf(format, v...)
	l.getLogger(LOG_LEVEL_ERROR).Printf(format, v...)
	if !p && l.logLevel >= LOG_LEVEL_ERROR {
		log.Printf(format, v...)
		return true
	}
	return false
}

func (l *Logger) Debug(content ...interface{}) (isPrint bool) {
	l.getLogger(LOG_LEVEL_DEBUG).Println(content...)
	if l.logLevel >= LOG_LEVEL_DEBUG {
		log.Println(content...)
		return true
	}
	return false
}

func (l *Logger) Info(content ...interface{}) (isPrint bool) {
	p := l.Debug(content...)
	l.getLogger(LOG_LEVEL_INFO).Println(content...)
	if !p && l.logLevel >= LOG_LEVEL_INFO {
		log.Println(content...)
		return true
	}
	return false
}

func (l *Logger) Warn(content ...interface{}) (isPrint bool) {
	p := l.Info(content...)
	l.getLogger(LOG_LEVEL_WARN).Println(content...)
	if !p && l.logLevel >= LOG_LEVEL_WARN {
		log.Println(content...)
		return true
	}
	return false
}

func (l *Logger) Error(content ...interface{}) (isPrint bool) {
	p := l.Warn(content...)
	l.getLogger(LOG_LEVEL_ERROR).Println(content...)
	if !p && l.logLevel >= LOG_LEVEL_ERROR {
		log.Println(content...)
		return true
	}
	return false
}

func (l *Logger) CloseLogFile() {
	l.lock.Lock()
	for _, lg := range l.loggerMap {
		if lg.logFile != nil && lg.isopened {
			lg.logFile.Close()
			lg.isopened = false
		}
	}
	l.lock.Unlock()
}
