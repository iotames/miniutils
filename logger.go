package miniutils

import (
	"fmt"
	"log"
	"os"
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
	logsDir   string
	loggerMap map[LogLevel]singleLogger
}

func getInitLoggerMap() map[LogLevel]singleLogger {
	return make(map[LogLevel]singleLogger, len(levelMap))
}

// TODO TO BE a global unique instance
func NewLogger(logsDir string) *Logger {
	return &Logger{logsDir: logsDir, loggerMap: getInitLoggerMap()}
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

func (l *Logger) getLogger(level LogLevel) *log.Logger {
	lg, ok := l.loggerMap[level]
	var logFile *os.File
	dateExpired := false
	if ok {
		filename := time.Now().Format(fmt.Sprintf("20060102_%s.log", levelMap[level]))
		logFile = l.loggerMap[level].logFile
		if logFile != nil {
			info, _ := logFile.Stat()
			if info != nil {
				if info.Name() != filename {
					logFile.Close()
					dateExpired = true
				}
			} else {
				log.Println("---getLogger--info, _ := logFile.Stat()--- info == nil-------")
			}
		} else {
			log.Println("---getLogger--logFile = l.loggerMap[level].logFile   --- == nil-------")
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

func (l *Logger) print(content ...interface{}) *Logger {
	log.Println(content...)
	l.getLogger(LOG_LEVEL_DEBUG).Println(content...)
	return l
}
func (l *Logger) Debug(conent ...interface{}) {
	l.print(conent...)
}

func (l *Logger) Info(content ...interface{}) {
	l.Debug(content...)
	l.getLogger(LOG_LEVEL_INFO).Println(content...)
}

func (l *Logger) Warn(content ...interface{}) {
	l.Info(content...)
	l.getLogger(LOG_LEVEL_WARN).Println(content...)
}

func (l *Logger) Error(content ...interface{}) {
	l.Warn(content...)
	l.getLogger(LOG_LEVEL_ERROR).Println(content...)
}

func (l *Logger) CloseLogFile() {
	for _, lg := range l.loggerMap {
		if lg.logFile != nil && lg.isopened {
			lg.logFile.Close()
			lg.isopened = false
		}
	}
}
