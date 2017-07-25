package logger

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const (
	_VER string = "1.0.0"
)

type UNIT int64

const (
	_       = iota
	KB UNIT = 1 << (iota * 10)
	MB
	GB
	TB
)

const (
	LOG int = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

const (
	OS_LINUX = iota
	OS_X
	OS_WIN
	OS_OTHERS
)

//日志结构对象
var logObj *LogFile
var logLevel = 1
var maxFileSize int64
var maxFileCount int32
var dailyFlag bool
var consoleAppender = false

const (
	//TimeDayFormat 日期格式化到日
	TimeDayFormat = "2006-01-02"
	//TimeFormat 日期格式化到秒
	TimeFormat = "2006-01-02 15:04:05"
)

var logFormat = "%s %s:%d %s %s"
var consoleFormat = "%s:%d %s %s"

//SetConsole 设置终端是否显示
func SetConsole(isConsole bool) {
	consoleAppender = isConsole
}

//SetLevel 设置日子级别
func SetLevel(_level int) {
	logLevel = _level
}

//NewRollingLogger 生成按文件大小及数量分割日子类
func NewRollingLogger(fileDir, fileName string, maxNumber int32, maxSize int64, _unit UNIT) {
	rollingLogger(fileDir, fileName, maxNumber, maxSize, _unit)
}

//SetRollingFile 生成按文件大小及数量分割日子类
func SetRollingFile(fileDir, fileName string, maxNumber int32, maxSize int64, _unit UNIT) {
	rollingLogger(fileDir, fileName, maxNumber, maxSize, _unit)
}

func rollingLogger(fileDir, fileName string, maxNumber int32, maxSize int64, _unit UNIT) {
	maxFileCount = maxNumber
	maxFileSize = maxSize * int64(_unit)
	dailyFlag = false
	logObj = &LogFile{dir: fileDir, filename: fileName, mu: new(sync.RWMutex)}
	logObj.mu.Lock()
	defer logObj.mu.Unlock()

	if !logObj.isMustRename() {
		logObj.logfile, _ = os.OpenFile(fileDir+"/"+fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	} else {
		logObj.rename()
	}
}

//NewDailyLogger new按日期分割日子类
func NewDailyLogger(fileDir, filename string) {
	dailyLogger(fileDir, filename)
}

func dailyLogger(fileDir, fileName string) {
	dailyFlag = true
	t, _ := time.Parse(TimeDayFormat, time.Now().Format(TimeDayFormat))
	logObj = &LogFile{dir: fileDir, filename: fileName, _date: &t, mu: new(sync.RWMutex)}
	logObj.mu.Lock()
	defer logObj.mu.Unlock()

	if !logObj.isMustRename() {
		var err error
		logObj.logfile, err = os.OpenFile(fileDir+"/"+fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
		if err != nil {
			fmt.Println("**** 111 ", err.Error)
		}
		fi, err := logObj.logfile.Stat()
		if err != nil {
			fmt.Println("aaaaaa")
			return
		}
		logObj.filesize = fi.Size()
	} else {
		logObj.rename()
	}
}

func concat(delimiter string, input ...interface{}) string {
	buffer := bytes.Buffer{}
	l := len(input)
	for i := 0; i < l; i++ {
		buffer.WriteString(fmt.Sprint(input[i]))
		if i < l-1 {
			buffer.WriteString(delimiter)
		}
	}
	return buffer.String()
}

func console(msg string) {
	if logObj == nil || logObj.logfile == nil || consoleAppender {
		log.Print(msg)
	}
}

func buildConsoleMessage(level int, msg string) string {
	file, line := getTraceFileLine()
	return fmt.Sprintf(logFormat+getOsEol(), time.Now().Format(TimeFormat), file, line, getTraceLevelName(level), msg)
}

func buildLogMessage(level int, msg string) string {
	file, line := getTraceFileLine()
	return fmt.Sprintf(logFormat+getOsEol(), time.Now().Format(TimeFormat), file, line, getTraceLevelName(level), msg)
}

func catchError() {
	if err := recover(); err != nil {
		log.Println("err", err)
	}
}

//Trace write
func Trace(level int, v ...interface{}) bool {
	defer catchError()
	if logObj != nil {
		logObj.mu.RLock()
		defer logObj.mu.RUnlock()
	}
	msg := concat(" ", v...)
	logStr := buildConsoleMessage(level, msg)
	console(logStr)
	if v[0] != nil && v[0].(string) == "remote" {
		remoteMsg := concat(" ", v[1:]...)
		go httpLog(remoteMsg)
	}
	if level >= logLevel {
		if logObj != nil && logObj.logfile != nil {
			logMsg := buildLogMessage(level, msg)
			_, err := logObj.write([]byte(logMsg))
			if err != nil {
				fmt.Println(err.Error())
				return false
			}
		}
	}
	return true
}

//Log LOG
func Log(v ...interface{}) bool {
	return Trace(LOG, v...)
}

//Debug DEBUG
func Debug(v ...interface{}) bool {
	return Trace(DEBUG, v...)
}

//Info INFO
func Info(v ...interface{}) bool {
	return Trace(INFO, v...)
}

//Warn WARN
func Warn(v ...interface{}) bool {
	return Trace(WARN, v...)
}

//Error ERROR
func Error(v ...interface{}) bool {
	return Trace(ERROR, v...)
}

//Fatal FATAL
func Fatal(v ...interface{}) bool {
	return Trace(FATAL, v...)
}
