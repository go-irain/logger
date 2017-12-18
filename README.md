# logger
golang 基本日志类库，支持自定义logid,tag，普通格式及json格式

# 获取
    go get github.com/go-irain/logger

# 日志格式如下:

```json
    {"atime":"2017-12-18 10:45:50","bfile":"logger/logger_test.go 41","clevel":"LOG","dlogid":"bbbbbbbbbbbbb","etag":"getUserInfo","msg":"rolling  fatal log"}
```

     2017-12-18 10:51:29 logger/logger_test.go:53 aaaaaaaaaaa login ERROR rolling  error log


* 日期 2017-07-06
* 时间 15:45:41
* 等级 INFO
* 编号 14e0e960525d46c7b1223cf09160cb38 (默认-)
* 标签 login (默认-) 刻录本条日志的事件名称或分类等。
* 代码 main:32 main.go文件的32行
* 内容 used time 0.0005s

LogOjb对象JSON方法设置了则以json格式输出，正常模式下以文本格式输出

# 切割

### 按文件大小切割

```golang
    RollingLogger("log", "a.log", 10, 1, MB)
```

### 按日期切割

```golang
    DailyLogger("log", "a.log")
```

# 使用示例

### 对象使用
```go
    package main

    import (
        "github.com/go-irain/logger"
    )

    func main(){
        logger.RollingLogger("log", "a.log", 10, 1, MB)
        logobj := new(LogObj).ID("bbbbbbbbbbbbb").Tag("getUserInfo").JSON()
        logobj.Log("rolling ", "log start")
        logobj.Debug("rolling ", "debug log")
        logobj.Info("rolling ", "info log")
        logobj.Warn("rolling ", "warn log")
        logobj.Error("rolling ", "error log")
        logobj.Fatal("rolling ", "fatal log")
        logobj.Log("rolling ", "fatal log")
    }
```

### 非对象使用

```go
    package main

    import (
        "github.com/go-irain/logger"
    )

    func main(){
        logger.RollingLogger("log", "a.log", 10, 1, MB)
        logger.Log("rolling ", "log start")
        logger.Debug("rolling ", "debug log")
        logger.Info("rolling ", "info log")
        logger.Warn("rolling ", "warn log")
        logger.Error("rolling ", "error log")
        logger.Fatal("rolling ", "fatal log")
        logger.Log("rolling ", "fatal log")
    }
    
```

# 远程日志

需要发短信预警的日志使用方式如下:

```go
    //需要到ums服务上注册告警服务id
    logger.SetRemoteUrl(alertURL) //ums告警地址
    logger.SetRemoteServerId(serverID)//ums告警服务id

    logger.Debug("remote","warning will send by sms")
```