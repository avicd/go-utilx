package logx

import (
	"github.com/avicd/go-utilx/conv"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const dayMilliSecs = 86400000
const Levels = 6

type Split = [Levels]bool

type RollCycle int

const (
	Day RollCycle = iota
	Minute
	Hour
	Week
	Month
	Year
	Second
)

var Layouts = []string{
	Year:   "2006",
	Month:  "2006-01",
	Day:    "2006-01-02",
	Hour:   "2006-01-02 15:00",
	Minute: "2006-01-02 15:04",
	Second: "2006-01-02 15:04:05",
}

type FileAppender struct {
	loggers    [Levels]*log.Logger
	files      [Levels]*os.File
	cycleIds   [Levels]int
	mutex      [Levels]*sync.Mutex
	locker     sync.Mutex
	TimeLayout string      // layout for time formatting
	MaxSize    int64       // max size of log file
	Name       string      // name of log file
	FileFlag   int         // FileFlag to open log file
	FileMode   os.FileMode // FileMode to open log file
	Prefix     string      // prefix of messages
	CycleOff   bool        // do not use the period cycle
	Cycle      RollCycle   // rolling cycle
	OutDir     string      // output directory cycle
	Split      Split       // split different level into different log file
}

func (it *FileAppender) getKey(level Level) Level {
	if it.Split[level] {
		return level
	}
	return ALL
}

func (it *FileAppender) getCycleId() (int, time.Time) {
	timeRef := time.Now()
	var id int
	switch it.Cycle {
	case Second:
		id = timeRef.Second()
	case Minute:
		id = timeRef.Minute()
	case Hour:
		id = timeRef.Hour()
	case Day:
		id = timeRef.Day()
	case Week:
		_, id = timeRef.ISOWeek()
	case Month:
		id = int(timeRef.Month())
	case Year:
		id = timeRef.Year()
	}
	return id, timeRef
}

func (it *FileAppender) init(key Level) {
	if it.mutex[key] == nil {
		it.locker.Lock()
		if it.mutex[key] == nil {
			it.mutex[key] = &sync.Mutex{}
		}
		it.locker.Unlock()
	}
}

func (it *FileAppender) sizeOver(file *os.File) bool {
	if it.MaxSize > 0 {
		info, _ := file.Stat()
		if info.Size() > it.MaxSize {
			return true
		}
	}
	return false
}

func (it *FileAppender) fileName(level Level) string {
	name := it.Name
	if name == "" {
		exeFile, _ := os.Executable()
		name = filepath.Base(exeFile)
		name = strings.TrimSuffix(name, filepath.Ext(name))
		it.Name = name
	}
	if it.Split[level] {
		name = conv.Append(name, strings.ToLower(Labels[level]), "_")
	}
	return name
}

func (it *FileAppender) rollLogFile(level Level) {
	key := it.getKey(level)
	var shortName string
	var cycleId int
	logFile := it.files[key]
	if !it.CycleOff {
		var timeRef time.Time
		cycleId, timeRef = it.getCycleId()
		if logFile != nil {
			if it.cycleIds[key] == cycleId {
				if !it.sizeOver(it.files[key]) {
					return
				}
			}
		}
		var cycleName string
		if it.Cycle == Week {
			year := time.Date(timeRef.Year(), 1, 1, 0, 0, 0, 0, time.Local)
			start := time.UnixMilli(year.UnixMilli() + int64(cycleId-1)*7*dayMilliSecs)
			end := time.UnixMilli(year.UnixMilli() + (int64(cycleId)*7-1)*dayMilliSecs)
			cycleName = start.Format(Layouts[Day]) + "--" + end.Format(Layouts[Day])
		} else {
			cycleName = timeRef.Format(Layouts[it.Cycle])
		}
		shortName = conv.Append(it.fileName(level), cycleName, "_")
	} else if logFile != nil && !it.sizeOver(logFile) {
		return
	} else {
		shortName = it.fileName(level)
	}
	it.mutex[key].Lock()
	if it.files[key] == logFile {
		if logFile != nil {
			logFile.Close()
		}
	} else {
		it.mutex[key].Unlock()
		return
	}
	index := 0
	fileFlag := it.FileFlag
	if fileFlag == 0 {
		fileFlag = os.O_WRONLY | os.O_APPEND | os.O_CREATE
	}
	fileMode := it.FileMode
	if fileMode == 0 {
		fileMode = 0666
	}
	for {
		fileName := filepath.Join(it.OutDir, shortName)
		if index > 0 {
			fileName = conv.Append(fileName, index, "_")
		}
		fileName += ".log"
		index++
		file, err := os.OpenFile(fileName, fileFlag, fileMode)
		if err != nil {
			panic(err)
		}
		if it.sizeOver(file) {
			continue
		}
		it.files[key] = file
		if !it.CycleOff {
			it.cycleIds[key] = cycleId
		}
		break
	}
	it.loggers[key] = log.New(it.files[key], it.Prefix, 0)
	it.mutex[key].Unlock()
}

func (it *FileAppender) Write(level Level, args ...any) {
	key := it.getKey(level)
	it.init(key)
	it.rollLogFile(level)
	var dest []any
	dest = append(dest, timeNow(it.TimeLayout))
	label := Labels[level] + " "
	dest = append(dest, label)
	dest = append(dest, args...)
	it.loggers[key].Print(dest...)
}
