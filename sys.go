package logger

import (
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func getTraceLevelName(level int) string {
	switch level {
	case LOG:
		return "LOG"
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

func getOsFlag() int {
	switch os := runtime.GOOS; os {
	case "darwin":
		return OS_X
	case "linux":
		return OS_LINUX
	case "windows":
		return OS_WIN
	default:
		return OS_OTHERS
	}
}

func getOsEol() string {
	if getOsFlag() == OS_WIN {
		return "\r\n"
	}
	return "\n"
}

var curStackFlag bool
var curStackPath string
var curStackLine int

func markCurStack() {
	if !curStackFlag {
		curStackFlag = true
		_, curStackPath, curStackLine, _ = runtime.Caller(0)
	}
}

func getCurStackPath() string {
	if !curStackFlag {
		markCurStack()
	}
	return curStackPath
}

func getStack(skip int) (pc uintptr, file string, line int, ok bool) {
	return runtime.Caller(skip)
}

func detectStack() (string, int) {
	curPath := getCurStackPath()
	for skip := 0; ; skip++ {
		_, path, line, ok := runtime.Caller(skip)
		if path != curPath {
			return path, line
		}
		if !ok {
			break
		}
	}
	return "", 0
}

func splitDirFile(path string) (string, string) {
	return filepath.Dir(path), filepath.Base(path)
}

func getTraceFileLine() (string, int) {
	fpath, line := detectStack()
	spath, file := splitDirFile(fpath)

	if getOsFlag() == OS_WIN {
		spath = strings.Replace(spath, "\\", "/", -1)
	}
	return path.Base(spath) + "/" + file, line
}
