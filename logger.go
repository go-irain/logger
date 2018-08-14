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

	"github.com/go-xorm/core"
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

type LogHandle struct {
	Handle int
}

//日志结构对象
var logObj *LogFile
var logLevel = 1
var maxFileSize int64
var maxFileCount int32
var dailyFlag bool
var consoleAppender = false
var showType = true
var level core.LogLevel
var handle LogHandle

const (
	//TimeDayFormat 日期格式化到日
	TimeDayFormat = "2006-01-02"
	//TimeFormat 日期格式化到秒
	TimeFormat = "2006-01-02 15:04:05"
	//TimeFormat = "2006-01-02T15:04:05.999999-07:00"
)

var logFormat = "%s %s:%d %s %s"
var logObjFormat = "%s %s:%d %s %s %s %s"
var consoleFormat = "%s:%d %s %s"

//SetConsole 设置终端是否显示
func SetConsole(isConsole bool) {
	consoleAppender = isConsole
}
func GetLogHandle() LogHandle {
	return handle
}

func (handle LogHandle) SetLevel(le core.LogLevel) {
	level = le
}

func (handle LogHandle) Level() core.LogLevel {
	return level
}

func (handle LogHandle) IsShowSQL() bool {
	return IsShowSQL()
}

func (handle LogHandle) ShowSQL(show ...bool) {
	if len(show) != 0 {
		ShowSQL(show[0])
	} else {
		ShowSQL(true)
	}
}

//Debug DEBUG
func (handle LogHandle) Debug(v ...interface{}) {
	Trace(DEBUG, nil, v...)
}

//Debug DEBUG
func (handle LogHandle) Debugf(format string, v ...interface{}) {
	Trace(DEBUG, nil, v...)
}

//Info INFO
func (handle LogHandle) Info(v ...interface{}) {
	Trace(INFO, nil, v...)
}

//Info INFO
func (handle LogHandle) Infof(format string, v ...interface{}) {
	Trace(INFO, nil, v...)
}

//Warn WARN
func (handle LogHandle) Warn(v ...interface{}) {
	Trace(WARN, nil, v...)
}

//Warn WARN
func (handle LogHandle) Warnf(format string, v ...interface{}) {
	Trace(WARN, nil, v...)
}

//Error ERROR
func (handle LogHandle) Error(v ...interface{}) {
	Trace(ERROR, nil, v...)
}

//Error ERROR
func (handle LogHandle) Errorf(format string, v ...interface{}) {
	Trace(ERROR, nil, v...)
}

//SetLevel 设置日子级别
func SetLevel(_level int) {
	logLevel = _level
}

//ShowSQL 设置是否打印SQL
func IsShowSQL() bool {
	return showType
}

//ShowSQL 设置是否打印SQL
func ShowSQL(show bool) {
	showType = show
}

//RollingLogger 生成按文件大小及数量分割日子类
func RollingLogger(fileDir, fileName string, maxNumber int32, maxSize int64, _unit UNIT) {
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

//DailyLogger new按日期分割日子类
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
	logInfo := map[string]interface{}{"atime": time.Now().Format(TimeFormat), "bfile": file + " " + strconv.Itoa(line), "clevel": getTraceLevelName(level)}

	if l != nil {
		logInfo["dlogid"] = l.logid
		logInfo["etag"] = l.tag
	}
	logInfo["msg"] = msg
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

//Trace write
func Trace(level int, l *LogObj, v ...interface{}) bool {
	defer catchError()
	if logObj != nil {
		logObj.mu.Lock()
		defer logObj.mu.Unlock()
	}
	msg := concat(" ", v...)
	logStr := ""
	if l != nil && l.json {
		logStr = buildJSONMessage(level, l, msg)
	} else {
		logStr = buildLogMessage(level, l, msg)
	}
	console(logStr)

	if v[0] != nil {
		if remote, ok := v[0].(string); ok && remote == "remote" {
			remoteMsg := concat(" ", v[1:]...)
			go httpLog(remoteMsg)
		}
	}
	if level >= logLevel {
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

//Log LOG
func Log(v ...interface{}) bool {
	return Trace(LOG, nil, v...)
}

//Debug DEBUG
func Debug(v ...interface{}) bool {
	return Trace(DEBUG, nil, v...)
}

//Info INFO
func Info(v ...interface{}) bool {
	return Trace(INFO, nil, v...)
}

//Warn WARN
func Warn(v ...interface{}) bool {
	return Trace(WARN, nil, v...)
}

//Error ERROR
func Error(v ...interface{}) bool {
	return Trace(ERROR, nil, v...)
}

//Fatal FATAL
func Fatal(v ...interface{}) bool {
	return Trace(FATAL, nil, v...)
}
