package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
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

// 日志结构对象
var logObj *LogFile
var logLevel = 1
var maxFileSize int64
var maxFileCount int32
var dailyFlag bool
var consoleAppender = false
var serviceName = ""

const (
	//TimeDayFormat 日期格式化到日
	TimeDayFormat = "2006-01-02"
	//TimeFormat 日期格式化到秒
	TimeFormat = "2006-01-02 15:04:05"
	//TimeFormat2 毫秒时间
	TimeFormat2 = "2006-01-02T15:04:05.000"
)

var logFormat = "%s %s:%d %s %s"
var logObjFormat = "%s %s:%d %s %s %s %s"
var consoleFormat = "%s:%d %s %s"

// SetConsole 设置终端是否显示
func SetConsole(isConsole bool) {
	consoleAppender = isConsole
}

// SetLevel 设置日子级别
func SetLevel(_level int) {
	logLevel = _level
}

// SetServiceName 设置服务名称
func SetServiceName(name string) {
	serviceName = name
}

// RollingLogger 生成按文件大小及数量分割日子类
func RollingLogger(fileDir, fileName string, maxNumber int32, maxSize int64, _unit UNIT) {
	rollingLogger(fileDir, fileName, maxNumber, maxSize, _unit)
}

// SetRollingFile 生成按文件大小及数量分割日子类
func SetRollingFile(fileDir, fileName string, maxNumber int32, maxSize int64, _unit UNIT) {
	rollingLogger(fileDir, fileName, maxNumber, maxSize, _unit)
}

func rollingLogger(fileDir, fileName string, maxNumber int32, maxSize int64, _unit UNIT) {
	maxFileCount = maxNumber
	maxFileSize = maxSize * int64(_unit)
	dailyFlag = false
	logObj = &LogFile{dir: fileDir, filename: fileName, mu: new(sync.Mutex)}
	logObj.mu.Lock()
	defer logObj.mu.Unlock()

	logObj.logfile, _ = os.OpenFile(fileDir+"/"+fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	fi, err := logObj.logfile.Stat()
	if err != nil {
		log.Println(err.Error())
		return
	}
	logObj.filesize = fi.Size()
}

// DailyLogger new按日期分割日子类
func DailyLogger(fileDir, filename string) {
	dailyLogger(fileDir, filename)
}

func dailyLogger(fileDir, fileName string) {
	dailyFlag = true
	t, _ := time.Parse(TimeDayFormat, time.Now().Format(TimeDayFormat))
	logObj = &LogFile{dir: fileDir, filename: fileName, _date: &t, mu: new(sync.Mutex)}
	logObj.mu.Lock()
	defer logObj.mu.Unlock()

	if !logObj.isMustRename() {
		var err error
		logObj.logfile, err = os.OpenFile(fileDir+"/"+fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
		if err != nil {
			log.Println(err.Error())
		}
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

func buildJSONMessage(level int, l *LogObj, msg string) string {
	file, line := getTraceFileLine()
	logInfo := map[string]interface{}{
		"timestamp": time.Now().Format(TimeFormat2),
		"service":   serviceName,
		"file":      file + " " + strconv.Itoa(line),
		"level":     getTraceLevelName(level),
		"message":   msg,
	}

	if l != nil {
		logInfo["guid"] = l.Logid()
		logInfo["action"] = l.GetTag()
		if l.GetData() != nil {
			logInfo["data"] = l.GetData()
		}
	}
	resb, err := json.Marshal(logInfo)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	return string(resb) + getOsEol()
}

func buildLogMessage(level int, l *LogObj, msg string) string {
	file, line := getTraceFileLine()
	logInfo := ""
	if l == nil {
		logInfo = fmt.Sprintf(logFormat+getOsEol(), time.Now().Format(TimeFormat), file, line, getTraceLevelName(level), msg)
	} else {
		logInfo = fmt.Sprintf(logObjFormat+getOsEol(), time.Now().Format(TimeFormat), file, line, l.logid, l.tag, getTraceLevelName(level), msg)
	}
	return logInfo
}

func catchError() {
	if err := recover(); err != nil {
		log.Println("err", err)
	}
}

// JSON 设置日志格式
func JSON(js bool) {
	if js {
		logObj.SetJSON()
	} else {
		logObj.UnSetJSON()
	}
}

// Trace write
func Trace(level int, l *LogObj, v ...interface{}) bool {
	defer catchError()
	if logObj != nil {
		logObj.mu.Lock()
		defer logObj.mu.Unlock()
	}
	msg := concat(" ", v...)
	logStr := ""
	if (logObj != nil && logObj.json) || (l != nil && l.json) {
		logStr = buildJSONMessage(level, l, msg)
	} else {
		logStr = buildLogMessage(level, l, msg)
	}
	console(logStr)
	if level >= logLevel && !consoleAppender {
		if logObj != nil {
			_, err := logObj.write([]byte(logStr))
			if err != nil {
				log.Println(err.Error())
				return false
			}
		}
	}
	return true
}

// Log LOG
func Log(v ...interface{}) bool {
	return Trace(LOG, nil, v...)
}

// Debug DEBUG
func Debug(v ...interface{}) bool {
	return Trace(DEBUG, nil, v...)
}

// Info INFO
func Info(v ...interface{}) bool {
	return Trace(INFO, nil, v...)
}

// Warn WARN
func Warn(v ...interface{}) bool {
	return Trace(WARN, nil, v...)
}

// Error ERROR
func Error(v ...interface{}) bool {
	return Trace(ERROR, nil, v...)
}

// Fatal FATAL
func Fatal(v ...interface{}) bool {
	return Trace(FATAL, nil, v...)
}
