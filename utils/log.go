package utils

import (
	"os"
	"path/filepath"
	"time"
)

// 滚动文件日志记录器。
type RollingFileWriter struct {
	Base   string // 日志文件夹。
	Prefix string // 日志文件前缀。
	f      *os.File
	today  time.Time
}

var (
	// 默认的滚动文件日志记录器。
	DefaultRollingFileWriter *RollingFileWriter = &RollingFileWriter{
		Base:   "./log",
		Prefix: "app",
	}
)

func (w *RollingFileWriter) Write(p []byte) (n int, err error) {
	now := time.Now()
	if today := TruncateToDay(now); today != w.today {
		// 新的一天开始，需要创建新文件。
		if w.f != nil {
			w.f.Sync()
			w.f.Close()
			w.f = nil
		}
		w.today = today
	}

	if w.f == nil {
		var bp string
		if v, err := filepath.Abs(w.Base); err != nil {
			panic(err)
		} else {
			bp = v
		}
		if err := os.MkdirAll(bp, 0666); err != nil {
			panic(err)
		}

		if nf, err := os.OpenFile(filepath.Join(w.Base, w.Prefix+w.today.Format("20060102")+".log"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666); err != nil {
			panic(err)
		} else {
			w.f = nf
		}
	}

	return w.f.Write(p)
}
