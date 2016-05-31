// Logger for ceyes
//
// + 对各类trace，加一个（W）这样的标记
// + log内容可以选择换行显示
// + 无序的可以前加标记X以标记
// + 文件大小和纪录条数的限制，任意一个满足都OK
// + Log文件自增后缀号码
//
// TODO:
//
// - To DB
// - For service
// - 停止Log的命令，从而退出死循环？
// - 写管道的时候如果已关闭是什么报错？
// - 命令管道只做命令传递
// - 如何量化性能指标？
// - error系统
// - session的log系统
// - 自起一个服务，可以提供运行时修改参数
// - json结构的部分更新
// - json解析的性能指标
// - 发送到指定服务器，或邮箱
// - 自动备份
// - web管理界面，实时更新
// - 有的内容写到文件，有的内容写到屏幕

package ceLogger

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

func init() {
	fmt.Println("Init CeLogger")
}

// ----------
// CeLogger
// ----------

type logEntry struct {
	i   uint
	buf []byte
}

type CeLogger struct {
	*CeLoggerConfig

	IsEnable bool

	mutex        sync.Mutex
	seqIndex     uint // auto increment seq index for all log entry
	isFirstEntry bool // init as true to indicate it is first log entry
	entryIndex   uint // entry index in current log file
	entryNum     uint // entry count in current log file
	fileSize     uint // file size of current log file
	writeCount   uint
	readCount    uint
	filename     string         // current log file name
	maxSeqIndex  uint           // max seq index, based on SeqIndexWidth. e.g. 4 -> 9999
	chLogEntry   chan *logEntry // channel for async write file
	chSeqIndex   chan uint      // channel for generate auto increment seq index
	chSeqInd     chan uint      // channel for signal, not in use now
	chLogInd     chan uint
}

// -- New CeLogger

func NewCeLogger() *CeLogger {
	return NewCeLoggerWithLogPath("")
}

func NewCeLoggerWithLogPath(filePath string) *CeLogger {
	cl := &CeLogger{CeLoggerConfig: NewCeLoggerConfig()}

	cl.LogFilePath = filePath
	if cl.LogFilePath == "" {
		// e.g. log_2006_01_02_15_04_05.log
		cl.LogFilePath = fmt.Sprintf("log_%s.log", time.Now().Format("2006_01_02_15_04_05"))
		cl.filename = cl.LogFilePath
	}

	return cl
}

func NewCeLoggerWithConfig(filepath string) *CeLogger {
	cl := &CeLogger{}
	cl.CeLoggerConfig = NewCeLoggerConfig()

	cl.LoadConfigFile(filepath)

	// Complete empty config value
	if cl.LogFilePath == "" {
		// e.g. log_2006_01_02_15_04_05.log
		cl.LogFilePath = fmt.Sprintf("log_%s.log", time.Now().Format("2006_01_02_15_04_05"))
		cl.filename = cl.LogFilePath
	}

	return cl
}

// -- Log type property: Trace/Info/Debug/Warn/Error/Panic

func (cl *CeLogger) IsLogTrace() bool {
	return cl.ECMap[ECTrace].IsEnable
}

func (cl *CeLogger) IsLogInfo() bool {
	return cl.ECMap[ECInfo].IsEnable
}

func (cl *CeLogger) IsLogDebug() bool {
	return cl.ECMap[ECDebug].IsEnable
}

func (cl *CeLogger) IsLogWarn() bool {
	return cl.ECMap[ECWarn].IsEnable
}

func (cl *CeLogger) IsLogError() bool {
	return cl.ECMap[ECError].IsEnable
}

func (cl *CeLogger) IsLogPanic() bool {
	return cl.ECMap[ECPanic].IsEnable
}

func (cl *CeLogger) SetLogTrace(b bool) *CeLogger {
	cl.ECMap[ECTrace].IsEnable = b
	return cl
}

func (cl *CeLogger) SetLogInfo(b bool) *CeLogger {
	cl.ECMap[ECInfo].IsEnable = b
	return cl
}

func (cl *CeLogger) SetLogDebug(b bool) *CeLogger {
	cl.ECMap[ECDebug].IsEnable = b
	return cl
}

func (cl *CeLogger) SetLogWarn(b bool) *CeLogger {
	cl.ECMap[ECWarn].IsEnable = b
	return cl
}

func (cl *CeLogger) SetLogError(b bool) *CeLogger {
	cl.ECMap[ECError].IsEnable = b
	return cl
}

func (cl *CeLogger) SetLogPanic(b bool) *CeLogger {
	cl.ECMap[ECPanic].IsEnable = b
	return cl
}

// -- log public functions: V/Vf/D/Df/I/If/W/Wf/E/Ef

// Log Trace
func (cl *CeLogger) Trace(tag string, e interface{}) *CeLogger {
	return cl.logWithTagColor(ECTrace, tag, e)
}

func (cl *CeLogger) Tracef(tag string, format string, params ...interface{}) *CeLogger {
	return cl.logfWithTagColor(ECTrace, tag, format, params...)
}

// Log Info
func (cl *CeLogger) Info(tag string, e interface{}) *CeLogger {
	return cl.logWithTagColor(ECInfo, tag, e)
}

func (cl *CeLogger) Infof(tag string, format string, params ...interface{}) *CeLogger {
	return cl.logfWithTagColor(ECInfo, tag, format, params...)
}

// Log Debug
func (cl *CeLogger) Debug(tag string, e interface{}) *CeLogger {
	return cl.logWithTagColor(ECDebug, tag, e)
}

func (cl *CeLogger) Debugf(tag string, format string, params ...interface{}) *CeLogger {
	return cl.logfWithTagColor(ECDebug, tag, format, params...)
}

// Log Warn
func (cl *CeLogger) Warn(tag string, e interface{}) *CeLogger {
	return cl.logWithTagColor(ECWarn, tag, e)
}

func (cl *CeLogger) Warnf(tag string, format string, params ...interface{}) *CeLogger {
	return cl.logfWithTagColor(ECWarn, tag, format, params...)
}

// Log Error
func (cl *CeLogger) Error(tag string, e interface{}) *CeLogger {
	return cl.logWithTagColor(ECError, tag, e)
}

func (cl *CeLogger) Errorf(tag string, format string, params ...interface{}) *CeLogger {
	return cl.logfWithTagColor(ECError, tag, format, params...)
}

// Log Panic
func (cl *CeLogger) Panic(tag string, e interface{}) *CeLogger {
	return cl.logWithTagColor(ECPanic, tag, e)
}

func (cl *CeLogger) Panicf(tag string, format string, params ...interface{}) *CeLogger {
	return cl.logfWithTagColor(ECPanic, tag, format, params...)
}

// -- Enter & Exit Func

func (cl *CeLogger) EnterFunc() (funcName string) {
	if !cl.IsEnable || !cl.IsLogFuncEnterExit {
		return ""
	}

	skip := 1
	pc, _, _, _ := runtime.Caller(skip)
	_, funcName = path.Split(runtime.FuncForPC(pc).Name())
	cl.log("+ " + funcName)

	return funcName
}

func (cl *CeLogger) ExitFunc(funcName string) {
	if !cl.IsEnable || !cl.IsLogFuncEnterExit {
		return
	}

	cl.log("- " + funcName)
}

// -- Get property

func (cl *CeLogger) GetFilename() string {
	return cl.filename
}

// -- Set property

func (cl *CeLogger) SetEnable(b bool) *CeLogger {
	if cl.IsEnable == b {
		return cl
	}
	cl.IsEnable = b

	if b {
		cl.applyConfig()

		// Init channel
		cl.chSeqIndex = make(chan uint, cl.ChanLen)
		cl.chLogEntry = make(chan *logEntry, cl.ChanLen)
		//		if cl.chSeqInd == nil {
		cl.chSeqInd = make(chan uint)
		//		}
		//		if cl.chLogInd == nil {
		cl.chLogInd = make(chan uint)
		//		}

		cl.seqIndex = 1
		cl.entryIndex = 0
		cl.entryNum = 0
		cl.fileSize = 0
		cl.filename = cl.LogFilePath
		cl.isFirstEntry = true
		cl.readCount = 0
		cl.writeCount = 0

		// Start channel routine
		go cl.handleEntryChannel()
		go cl.handleSeqIndexChannel()

		time.Sleep(time.Millisecond)
		fmt.Println("Log started")
	} else {
		fmt.Println("Log stopping ...")

		// Notify to close chSeqIndex
		//		fmt.Println("***** cl.chSeqInd <- 1")
		cl.chSeqInd <- 1

		// Wait for chLogEntry closed
		//		if cl.readCount == cl.writeCount && len(cl.chLogEntry) == 0 {
		//			close(cl.chLogEntry)
		//		} else {
		//			//			fmt.Println("***** <-cl.chLogInd")
		//			<-cl.chLogInd
		//		}
		<-cl.chLogInd
		time.Sleep(time.Millisecond)
		fmt.Println("Log stopped")
	}

	return cl
}

func (cl *CeLogger) SetChanLen(n uint) *CeLogger {
	if cl.ChanLen == n {
		return cl
	}
	cl.ChanLen = n

	// TODO: If logger is working, need to do some process before set chLen
	// TODO: Need to treat panic here
	cl.chLogEntry = make(chan *logEntry, cl.ChanLen)
	cl.chSeqIndex = make(chan uint, cl.ChanLen)

	return cl
}

func (cl *CeLogger) SetMaxFileSize(n uint) *CeLogger {
	cl.MaxFileSize = n
	return cl
}
func (cl *CeLogger) SetMaxEntryNum(n uint) *CeLogger {
	cl.MaxEntryNum = n
	return cl
}

func (cl *CeLogger) SetSyncWriteFile(b bool) *CeLogger {
	cl.IsSyncWriteFile = b
	return cl
}

func (cl *CeLogger) SetLogFuncEnterExit(b bool) *CeLogger {
	cl.IsLogFuncEnterExit = b
	return cl
}

func (cl *CeLogger) SetLogSeqIndex(b bool) *CeLogger {
	cl.IsLogSeqIndex = b
	return cl
}

func (cl *CeLogger) SetSeqIndexWidth(n uint) *CeLogger {
	cl.SeqIndexWidth = n
	cl.maxSeqIndex = 0
	if n > 0 {
		cl.maxSeqIndex = uint(math.Pow10(int(n))) - 1
	}
	return cl
}

func (cl *CeLogger) SetLogEntryTag(b bool) *CeLogger {
	cl.IsLogEntryTag = b
	return cl
}

func (cl *CeLogger) SetLogCodeFilename(b bool) *CeLogger {
	cl.IsLogCodeFilename = b
	return cl
}

func (cl *CeLogger) SetLogCodeLineNumber(b bool) *CeLogger {
	cl.IsLogCodeLineNumber = b
	return cl
}

func (cl *CeLogger) SetLogCodeFuncName(b bool) *CeLogger {
	cl.IsLogCodeFuncName = b
	return cl
}

func (cl *CeLogger) SetLogDate(b bool) *CeLogger {
	cl.IsLogDate = b
	return cl
}

func (cl *CeLogger) SetLogTime(b bool) *CeLogger {
	cl.IsLogTime = b
	return cl
}

func (cl *CeLogger) SetTimeMsWidth(n uint) *CeLogger {
	if n > 9 {
		n = 9
	}
	cl.TimeMsWidth = n
	return cl
}

func (cl *CeLogger) SetLogColor(b bool) *CeLogger {
	cl.IsLogColor = b
	return cl
}

func (cl *CeLogger) SetWriteFile(isWriteToFile bool) *CeLogger {
	cl.IsWriteFile = isWriteToFile
	return cl
}

func (cl *CeLogger) SetWriteConsole(b bool) *CeLogger {
	cl.IsWriteConsole = b
	return cl
}

func (cl *CeLogger) SetLogOrderFlag(b bool) *CeLogger {
	cl.IsLogOrderFlag = b
	return cl
}

func (cl *CeLogger) SetContentDelimiter(str string) *CeLogger {
	cl.ContentDelimiter = str
	return cl
}

func (cl *CeLogger) SetLogFilePath(filePath string) *CeLogger {
	cl.LogFilePath = filePath
	cl.filename = filePath

	return cl
}

func (cl *CeLogger) SetConfigFilePath(filePath string) *CeLogger {
	return NewCeLoggerWithConfig(filePath)
}

// -- private log function

func (cl *CeLogger) log(e interface{}) *CeLogger {
	if !cl.IsEnable {
		return cl
	}

	var buf bytes.Buffer

	if cl.IsLogOrderFlag {
		buf.WriteString(" ")
	}

	// Sequence index
	i := cl.getSeqIndex()
	if cl.IsLogSeqIndex && i > 0 {
		buf.WriteString(cl.getSeqIndexString(i))
	}
	buf.WriteString(cl.getDateTimeString())
	buf.WriteString(cl.getFuncInfoString())

	// Default is " "
	buf.WriteString(cl.ContentDelimiter)

	// Real cl content
	buf.WriteString(cl.getString(e))

	// Write log entry
	cl.writeEntry(&logEntry{i, buf.Bytes()})

	return cl
}

func (cl *CeLogger) logf(format string, params ...interface{}) *CeLogger {
	if !cl.IsEnable {
		return cl
	}

	cl.log(fmt.Sprintf(format, params...))

	return cl
}

func (cl *CeLogger) logWithTagColor(etName, tag string, e interface{}) *CeLogger {
	if !cl.IsEnable {
		return cl
	}

	ec, ok := cl.ECMap[etName]
	if !ok {
		ec = cl.ECMap[""]
	}
	if !ec.IsEnable {
		return cl
	}

	var s string
	if cl.IsLogEntryTag {
		s = cl.getTagString(ec.Tag)
	}

	var buf bytes.Buffer
	if cl.IsLogColor {
		buf.WriteString(cl.GetColorString(s+cl.getTagString(tag)+cl.getString(e), ec))
	} else {
		buf.WriteString(s + cl.getTagString(tag) + cl.getString(e))
	}
	cl.log(buf.String())

	return cl
}

func (cl *CeLogger) logfWithTagColor(etName, tag string, format string, params ...interface{}) *CeLogger {
	ec, ok := cl.ECMap[etName]
	if !ok {
		ec = cl.ECMap[""]
	}
	if !cl.IsEnable || !ec.IsEnable {
		return cl
	}

	var s string
	if cl.IsLogEntryTag {
		s = cl.getTagString(ec.Tag)
	}

	var buf bytes.Buffer
	if cl.IsLogColor {
		buf.WriteString(cl.GetColorString(s+cl.getTagString(tag)+fmt.Sprintf(format, params...), ec))
	} else {
		buf.WriteString(s + cl.getTagString(tag) + fmt.Sprintf(format, params...))
	}
	cl.log(buf.String())

	return cl
}

// Write log entry
func (cl *CeLogger) writeEntry(entry *logEntry) *CeLogger {
	if !cl.IsEnable {
		return cl
	}

	if cl.IsWriteConsole {
		fmt.Println(string(entry.buf))
	}

	if cl.IsWriteFile {
		if cl.IsSyncWriteFile {
			cl.mutex.Lock()
			defer cl.mutex.Unlock()

			// Sync write file
			err := cl.writeEntryToFile(entry)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			// Async write file
			go func(e *logEntry) {
				if !cl.IsEnable {
					return
				}

				cl.writeCount++
				cl.chLogEntry <- e
			}(entry)

			time.Sleep(time.Nanosecond)
		}
	}

	return cl
}

func (cl *CeLogger) writeEntryToFile(entry *logEntry) error {
	if cl.IsLogOrderFlag {
		// Mark "X" in front of log entry if out of order
		if !(entry.i == cl.entryIndex+1 ||
			(entry.i == 1 && cl.entryIndex == uint(math.Pow10(int(cl.SeqIndexWidth)))-1)) {
			entry.buf[0] = 'X'
		}
	}

	// If this is the first entry, check log file if available
	// Refresh log.filename if nesessary
	if cl.isFirstEntry {
		if cl.MaxFileSize > 0 || cl.MaxEntryNum > 0 {
			if _, err := os.Stat(cl.LogFilePath); os.IsNotExist(err) {
				// If log file not exist
				cl.filename = cl.LogFilePath
			} else {
				// If log file exist, try to get next one
				// e.g. "test.log" -> "test_1.log"
				cl.filename = cl.getNextValidFilename(cl.LogFilePath)
			}
		}
		cl.isFirstEntry = false
	}

	// Write to a new log file if necessary
	if (cl.MaxEntryNum > 0 && cl.entryNum >= cl.MaxEntryNum) ||
		(cl.MaxFileSize > 0 && cl.fileSize+uint(len(entry.buf))+1 >= cl.MaxFileSize) {
		cl.filename = cl.getNextValidFilename(cl.LogFilePath)
		// reset log tracking data
		cl.fileSize = 0
		cl.entryNum = 0
	}

	// Write buf to log file
	size, err := cl.writeLogFile(entry.buf)
	if err != nil {
		fmt.Sprintf("writeLogFile failed at %d, %s\n", entry.i, err.Error())
		return err
	}

	// Update log tracking data
	cl.entryIndex = entry.i
	cl.fileSize += size
	cl.entryNum++

	return nil
}

// Write buffer to log file, return size of write bytes
func (cl *CeLogger) writeLogFile(buf []byte) (size uint, err error) {
	file, err := os.OpenFile(cl.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Write log failed.", err.Error())
		return 0, err
	}
	defer file.Close()

	file.Write(buf)
	file.WriteString("\n")

	return uint(len(buf)) + 1, nil
}

func (cl *CeLogger) handleEntryChannel() {
	var timer *time.Timer
Loop:
	for {
		select {
		case entry, ok := <-cl.chLogEntry:
			if !ok {
				//fmt.Println("chLogEntry is closed")
				break Loop
			}

			cl.readCount++
			if err := cl.writeEntryToFile(entry); err != nil {
				fmt.Printf("writeEntryToFile failed: %s\n", err.Error())
			}
		case <-func() <-chan time.Time {
			if timer == nil {
				timer = time.NewTimer(time.Millisecond * 20)
			} else {
				timer.Reset(time.Millisecond * 20)
			}
			return timer.C
		}():
			// If stop logger, close chLogEntry when read = write
			if !cl.IsEnable && (cl.writeCount == cl.readCount) {
				//fmt.Println("close chLogEntry")
				close(cl.chLogEntry)
				cl.chLogInd <- 1
			}
		}
	}

	//	fmt.Println("***** - handleEntryChannel")
}

func (cl *CeLogger) handleSeqIndexChannel() {
Loop:
	for {
		select {
		case cl.chSeqIndex <- cl.seqIndex:
			cl.seqIndex++

			if cl.SeqIndexWidth > 0 {
				if cl.maxSeqIndex == 0 {
					cl.SetSeqIndexWidth(cl.SeqIndexWidth)
				}

				if cl.seqIndex >= cl.maxSeqIndex {
					cl.seqIndex = 1
				}
			}
		case _ = <-cl.chSeqInd:
			//fmt.Println("close chSeqIndex", flag)
			close(cl.chSeqIndex)
			break Loop
		}
	}
	//	fmt.Println("***** - handleSeqIndexChannel")
}

// -- private helper function

// Convert everything to string
func (cl *CeLogger) getString(e interface{}) string {
	type Stringer interface {
		String() string
	}

	switch i := e.(type) {
	case string:
		return string(i)
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
		return fmt.Sprintf("%d", i)
	case float32, float64:
		return fmt.Sprintf("%f", i)
	case Stringer:
		return i.String()
	default:
		return fmt.Sprintf("%v", i)
	}
}

// Index
// e.g. 12
func (cl *CeLogger) getSeqIndex() uint {
	if !cl.IsEnable || !cl.IsLogSeqIndex {
		return 0
	}

	// TODO: Need to take care of blocking at chSeqIndex
	i, ok := <-cl.chSeqIndex
	if !ok {
		fmt.Println("SeqIndex chan closed")
		return 0
	}
	return i
}

// Index string
// [Index], e.g. [0012]
func (cl *CeLogger) getSeqIndexString(i uint) string {
	if cl.SeqIndexWidth == 0 {
		return ""
	}
	s := strconv.Itoa(int(cl.SeqIndexWidth))
	return fmt.Sprintf("[%0"+s+"d]", i)
}

// Date time string
// [date time], e.g. [2015-03-04 15:16:17]
func (cl *CeLogger) getDateTimeString() string {
	if !cl.IsEnable {
		return ""
	}

	if !(cl.IsLogDate || cl.IsLogTime) {
		return ""
	}

	t := time.Now()

	var buf bytes.Buffer

	buf.WriteString("[")

	if cl.IsLogDate {
		buf.WriteString(t.Format("2006-01-02"))
	}

	if cl.IsLogTime {
		if buf.Len() > 1 {
			buf.WriteString(" ")
		}

		if cl.TimeMsWidth == 0 {
			s := t.Format("15:04:05")
			buf.WriteString(s)
		} else {
			s := t.Format("15:04:05.999999999")
			n := len("15:04:05.") + int(cl.TimeMsWidth) - len(s)

			switch {
			case n < 0:
				buf.WriteString(s[:len(s)+n])
			case n == 0:
				buf.WriteString(s)
			case n > 0:
				buf.WriteString(s)
				buf.WriteString(strings.Repeat("0", n))
			}
		}
	}

	buf.WriteString("]")
	return buf.String()
}

// Func info
// (filename:line-package.func), e.g. (abc.go:12-main.test)
func (cl *CeLogger) getFuncInfoString() string {
	if !cl.IsEnable {
		return ""
	}

	if !(cl.IsLogCodeFilename || cl.IsLogCodeFuncName) {
		return ""
	}

	skip := 4
	pc, filename, lineNumber, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}

	var buf bytes.Buffer

	buf.WriteString("(")

	if cl.IsLogCodeFilename {
		_, f := path.Split(filename)
		buf.WriteString(f)

		if cl.IsLogCodeLineNumber {
			buf.WriteString(fmt.Sprintf(":%d", lineNumber))
			//buf.WriteString(":")
			//buf.WriteString(strconv.Itoa(lineNumber))
		}
	}

	if cl.IsLogCodeFuncName {
		if buf.Len() > 1 {
			buf.WriteString("-")
		}
		_, funcName := path.Split(runtime.FuncForPC(pc).Name())

		buf.WriteString(funcName)
	}

	buf.WriteString(")")
	return buf.String()
}

// Tag string
// [tag], e.g. [HTTP]
func (cl *CeLogger) getTagString(tag string) string {
	var buf bytes.Buffer

	buf.WriteString("[")
	buf.WriteString(tag)
	buf.WriteString("]")

	return buf.String()
}

// Colorized string
// e.g. '0x1B'[1;31m'0x1B'[1;41mHappy'0x1B'[0m
func (c *CeLoggerConfig) GetColorString(str string, ec *EntryConfig) string {
	if !ec.IsEnable || str == "" {
		return str
	}

	// e.g. printf("\033[1;40;32m%s\033[0m",” Hello,NSFocus\n”);

	// 前景 背景 颜色
	// --------------------
	// 30  40  黑色
	// 31  41  红色
	// 32  42  绿色
	// 33  43  黄色
	// 34  44  蓝色
	// 35  45  紫红色
	// 36  46  青蓝色
	// 37  47  白色
	//
	// 代码 意义
	// --------------------
	//  0  终端默认设置
	//  1  高亮显示
	//  4  使用下划线
	//  5  闪烁
	//  7  反白显示
	//  8  不可见

	switch ec.DisplayMode {
	case 0, 1, 4, 5, 7, 8:
	default:
		ec.DisplayMode = 0
	}

	if ec.ForeColor < 30 || ec.ForeColor > 37 {
		ec.ForeColor = 37
	}

	var buf bytes.Buffer

	if ec.BackColor < 40 || ec.BackColor > 47 {
		// Use bytes.Buffer instead of fmt.Sprintf() for better performance
		//return fmt.Sprintf("%c[%d;%dm%s%c[0m", 0x1B, ct.DisplayMode, ct.ForeColor, str, 0x1B)

		buf.WriteByte(0x1B)
		buf.WriteString("[")
		buf.WriteString(strconv.Itoa(int(ec.DisplayMode)))
		buf.WriteString(";")
		buf.WriteString(strconv.Itoa(int(ec.ForeColor)))
		buf.WriteString("m")
		buf.WriteString(str)
		buf.WriteByte(0x1B)
		buf.WriteString("[0m")
	} else {
		//return fmt.Sprintf("%c[%d;%d;%dm%s%c[0m", 0x1B, ct.DisplayMode, ct.BackColor, ct.ForeColor, str, 0x1B)
		buf.WriteByte(0x1B)
		buf.WriteString("[")
		buf.WriteString(strconv.Itoa(int(ec.DisplayMode)))
		buf.WriteString(";")
		buf.WriteString(strconv.Itoa(int(ec.BackColor)))
		buf.WriteString(";")
		buf.WriteString(strconv.Itoa(int(ec.ForeColor)))
		buf.WriteString("m")
		buf.WriteString(str)
		buf.WriteByte(0x1B)
		buf.WriteString("[0m")
	}

	return buf.String()
}

func (cl *CeLogger) getFilenameExt(path string) (name, ext string) {
	for i := len(path) - 1; i >= 0 && path[i] != '/'; i-- {
		if path[i] == '.' {
			return path[:i], path[i+1:]
		}
	}
	return path, ""
}

func (cl *CeLogger) getFileSize(path string) (size int64, err error) {
	fi, err := os.Stat(path)
	if !(err == nil || os.IsExist(err)) {
		return 0, err
	}

	return fi.Size(), nil
}

func (cl *CeLogger) getNextValidFilename(path string) string {
	i := 0
	f, e := cl.getFilenameExt(path)
	var format string
	if e == "" {
		format = "%s_%d%s"
	} else {
		format = "%s_%d.%s"
	}

	for {
		i++
		filename := fmt.Sprintf(format, f, i, e)
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			return filename
		}
	}
}

func (cl *CeLogger) applyConfig() *CeLogger {
	cl.SetSeqIndexWidth(cl.SeqIndexWidth)

	return cl
}
