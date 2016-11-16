package logger

import (
	"fmt"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	SetConsole(false)
	//根据配置文件设置日志等级

	SetLevel(DEBUG)
	fmt.Println(INFO)
	//根据配置文件，设置日志路径，日志名，日志切割大小限制
	SetRollingFile("./", "a.log", 10, 50, MB)
	//设置远程告警地址及服务id
	SetRemoteUrl("http://url:port/api/alert/report")
	SetRemoteServerId("1")
	Log("log start")
	Debug("debug log")
	Info("info log")
	Warn("warn log")
	Error("error log")
	Fatal("fatal log")
	time.Sleep(1 * time.Second)
}

func TestFarstWrite(t *testing.T) {
	t.Log("start testing", time.Now().Unix())
	for index := 0; index < 10; index++ {
		Info("aaaaaaaaaaaaaaaaaa :=", index)
	}
	t.Log("end testing", time.Now().Unix())
}
