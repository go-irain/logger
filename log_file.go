package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

//LogFile 日子结构
type LogFile struct {
	dir      string
	filename string
	filesize int64
	mu       *sync.Mutex
	logfile  *os.File
	_suffix  int
	_date    *time.Time
}

func (f *LogFile) isMustRename() bool {
	if dailyFlag {
		t, _ := time.Parse(TimeDayFormat, time.Now().Format(TimeDayFormat))
		if t.After(*f._date) {
			return true
		}
	} else {
		if maxFileSize > 0 && f.filesize > maxFileSize {
			return true
		}
	}
	return false
}

func (f *LogFile) rename() {
	if dailyFlag {
		fn := f.dir + "/" + f.filename + "." + f._date.Format(TimeDayFormat)
		if !f.isExist(fn) && f.isMustRename() {
			if f.logfile != nil {
				f.logfile.Close()
			}
			err := os.Rename(f.dir+"/"+f.filename, fn)
			if err != nil {
				fmt.Println("rename err", err.Error())
			}
			t, _ := time.Parse(TimeDayFormat, time.Now().Format(TimeDayFormat))
			f._date = &t
			f.logfile, _ = os.Create(f.dir + "/" + f.filename)
		}
	} else {
		if maxFileSize > 0 && f.filesize > maxFileSize {
			err := f.rotate()
			if err != nil {
				fmt.Println("333", err.Error())
			}
		}
	}
}

func (f *LogFile) write(data []byte) (int, error) {
	n, err := f.logfile.Write(data)
	if err != nil {
		fmt.Println("11111", err.Error())
		return n, err
	}
	f.filesize += int64(n)
	f.rename()
	return n, err
}

// 获取目录下指定前缀的所有日志文件
func (f *LogFile) removeFiles() {
	fs, err := filepath.Glob(fmt.Sprintf("%s/%s.*", f.dir, f.filename))
	if err != nil {
		return
	}
	sort.Strings(fs)
	x := len(fs) - (int(maxFileCount) - 1)
	if maxFileCount > 0 && x > 0 {
		dels := fs[:x]
		for _, v := range dels {
			os.Remove(v)
		}
	}
}

// 分割
func (f *LogFile) rotate() error {
	f.removeFiles()
	if f != nil && f.logfile != nil {
		f.logfile.Sync()
		f.logfile.Close()
		os.Rename(f.dir+"/"+f.filename, f.dir+"/"+f.filename+time.Now().Format(".20060102150405"))
	}
	// 创建最新的日志文件
	fd, err := os.OpenFile(f.dir+"/"+f.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	fi, err := fd.Stat()
	if err != nil {
		return err
	}
	f.logfile = fd
	f.filesize = fi.Size()
	return nil
}

func (f *LogFile) isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
