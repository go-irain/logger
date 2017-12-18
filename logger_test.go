package logger

import (
	//"fmt"
	"testing"
	"time"
)

func TestLoggerDefault(t *testing.T) {

	//根据配置文件设置日志等级

	for i := 0; i < 1; i++ {
		Log("log start")
		Debug("debug log")
		Info("info log")
		Warn("warn log")
		Error("error log")
		Fatal("fatal log")
		Info("remote", " log")
		Log("log log")
	}
}

func TestLoggerRolling(t *testing.T) {
	SetConsole(true)
	SetLevel(DEBUG)
	//根据配置文件，设置日志路径，日志名，日志切割大小限制
	RollingLogger("log", "a.log", 10, 1, MB)
	go func() {
		logobj := new(LogObj).ID("bbbbbbbbbbbbb").Tag("getUserInfo")
		for i := 0; i < 100000; i++ {
			time.Sleep(1 * time.Millisecond)
			go func() {
				logobj.Log("rolling ", "log start")
				logobj.Debug("rolling ", "debug log")
				logobj.Info("rolling ", "info log")
				logobj.Warn("rolling ", "warn log")
				logobj.Error("rolling ", "error log")
				logobj.Fatal("rolling ", "fatal log")
				logobj.Log("rolling ", "fatal log")
			}()
		}
	}()
	logobj := new(LogObj).ID("aaaaaaaaaaa").Tag("login").JSON()
	for i := 0; i < 100000; i++ {
		time.Sleep(1 * time.Millisecond)
		go func() {
			logobj.Log("rolling ", "log start")
			logobj.Debug("rolling ", "debug log")
			logobj.Info("rolling ", "info log")
			logobj.Warn("rolling ", "warn log")
			logobj.Error("rolling ", "error log")
			logobj.Fatal("rolling ", "fatal log")
			logobj.Log("rolling ", "fatal log")
		}()
	}

	time.Sleep(30 * time.Second)
}

/*
func TestLoggerDaily(t *testing.T) {
	//SetConsole(true)
	DailyLogger("log", "a.log")
	//根据配置文件，设置日志路径，日志名，日志切割大小限制
	for i := 0; i < 100000000; i++ {
		Log("daily ", "log start")
		Debug("daily ", "debug log")
		Info("daily ", "info log")
		Warn("daily ", "warn log")
		Error("daily ", "error log")
		Fatal("daily ", "fatal log")
		Log("daily ", "log log")
		time.Sleep(time.Microsecond * 100)
	}
}
*/
