package ceLogger

import (
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"
)

func init() {
	runtime.GOMAXPROCS(4)
}

func TestLogDefault(t *testing.T) {
	//t.Skip()
	l := NewCeLogger()
	os.Remove(l.GetFilename())

	t.Log("TestLogDefault")
	l.SetEnable(true)
	logAllType(l)
	l.SetEnable(false)
}

func TestMaxFileSize(t *testing.T) {
	//t.Skip()
	l := NewCeLogger()
	l.SetLogFilePath("TestMaxFileSize-256.log")
	os.Remove(l.LogFilePath)

	t.Log("SetMaxFileSize(256)")
	l.SetMaxFileSize(256)
	l.SetEnable(true)
	logAllType(l)
	l.SetEnable(false)

	t.Log("SetMaxFileSize(1024)")
	l.SetMaxFileSize(1024)
	l.SetLogFilePath("TestMaxFileSize-1024.log")
	os.Remove(l.LogFilePath)
	l.SetEnable(true)
	logAllType(l)
	l.SetEnable(false)
}

func TestMaxEntryNum(t *testing.T) {
	//t.Skip()
	l := NewCeLogger()
	l.SetLogSeqIndex(true)

	l.SetLogFilePath("TestMaxEntryNum-5.log")
	os.Remove(l.LogFilePath)

	t.Log("SetMaxEntryNum(5)")
	l.SetMaxEntryNum(5)
	l.SetEnable(true)
	logAllType(l)
	l.SetEnable(false)

	l.SetEnable(true)
	t.Log("SetMaxEntryNum(10)")
	l.SetMaxEntryNum(10)
	l.SetLogFilePath("TestMaxEntryNum-10.log")
	logAllType(l)

	l.SetEnable(false)
}

func TestSeqIndexWidth(t *testing.T) {
	//t.Skip()
	l := NewCeLogger()
	l.SetLogFilePath("TestSeqIndexWidth.log")
	os.Remove(l.LogFilePath)

	t.Log("SetSeqIndexWidth(0)")
	l.SetSeqIndexWidth(0)
	l.SetEnable(true)
	logAllType(l)
	l.SetEnable(false)

	l.SetEnable(true)
	t.Log("SetSeqIndexWidth(3)")
	l.SetSeqIndexWidth(3)
	logAllType(l)
	l.SetEnable(false)
}

func TestTimeMsWidth(t *testing.T) {
	//t.Skip()
	l := NewCeLogger()
	l.SetLogFilePath("TestTimeMsWidth.log")
	os.Remove(l.LogFilePath)

	t.Log("SetTimeMsWidth(0)")
	l.SetTimeMsWidth(0)
	l.SetEnable(true)
	logAllType(l)
	l.SetEnable(false)

	l.SetEnable(true)
	t.Log("SetTimeMsWidth(3)")
	l.SetTimeMsWidth(3)
	logAllType(l)
	l.SetEnable(false)

	l.SetEnable(true)
	t.Log("SetTimeMsWidth(12)")
	l.SetTimeMsWidth(12)
	logAllType(l)
	l.SetEnable(false)
}

func TestSyncWriteFile(t *testing.T) {
	//t.Skip()
	l := NewCeLogger()
	l.SetLogFilePath("TestSyncWriteFile.log")
	os.Remove(l.LogFilePath)

	t.Log("SetSyncWriteFile(true)")
	l.SetSyncWriteFile(true)
	l.SetEnable(true)
	logAllType(l)
	l.SetEnable(false)

	t.Log("SetSyncWriteFile(false)")
	l.SetSyncWriteFile(false)
	l.SetEnable(true)
	logAllType(l)
	l.SetEnable(false)

}

func TestLogOrderFlag(t *testing.T) {
	//t.Skip()
	l := NewCeLogger()
	l.SetLogFilePath("TestLogOrderFlag.log")
	os.Remove(l.LogFilePath)

	t.Log("SetLogOrderFlag(true)")
	l.SetLogOrderFlag(true)
	l.SetEnable(true)
	logAllType(l)
	l.SetEnable(false)

	t.Log("SetLogOrderFlag(false)")
	l.SetLogOrderFlag(false)
	l.SetEnable(true)
	logAllType(l)
	l.SetEnable(false)
}

func TestLogSeqIndex(t *testing.T) {
	//t.Skip()
	l := NewCeLogger()
	l.SetLogFilePath("TestLogSeqIndex.log")
	os.Remove(l.LogFilePath)

	t.Log("SetLogSeqIndex(true)")
	l.SetLogSeqIndex(true)
	l.SetEnable(true)
	logAllType(l)
	l.SetEnable(false)

	t.Log("SetLogSeqIndex(false)")
	l.SetLogSeqIndex(false)
	l.SetEnable(true)
	logAllType(l)
	l.SetEnable(false)
}

func TestLogFuncEnterExit(t *testing.T) {
	//t.Skip()
	cl := NewCeLogger()
	cl.SetLogFilePath("TestLogFuncEnterExit.log")
	os.Remove(cl.LogFilePath)

	t.Log("SetLogFuncEnterExit(true)")
	cl.SetLogFuncEnterExit(true)
	cl.SetEnable(true)

	func() {
		defer cl.ExitFunc(cl.EnterFunc())
		fmt.Println("I am in func F")
	}()
	cl.SetEnable(false)

	t.Log("SetLogFuncEnterExit(false)")
	cl.SetLogFuncEnterExit(false)
	cl.SetEnable(true)
	func() {
		defer cl.ExitFunc(cl.EnterFunc())
		fmt.Println("I am in func F")
	}()
	cl.SetEnable(false)
}

func TestLogCodeFilename(t *testing.T) {
	//t.Skip()
	l := NewCeLogger()
	l.SetLogFilePath("TestLogCodeFilename.log")
	os.Remove(l.LogFilePath)

	t.Log("SetLogCodeFilename(true)")
	l.SetLogCodeFilename(true)
	l.SetEnable(true)
	logAllType(l)
	l.SetEnable(false)

	t.Log("SetLogCodeFilename(false)")
	l.SetLogCodeFilename(false)
	l.SetEnable(true)
	logAllType(l)
	l.SetEnable(false)
}

func TestLogCodeLineNumber(t *testing.T) {
	//t.Skip()
	l := NewCeLogger()
	l.SetLogFilePath("TestLogCodeLineNumber.log")
	os.Remove(l.LogFilePath)

	t.Log("SetLogCodeLineNumber(true)")
	l.SetLogCodeFilename(true).SetLogCodeLineNumber(true)
	l.SetEnable(true)
	logAllType(l)
	l.SetEnable(false)

	t.Log("SetLogCodeLineNumber(false)")
	l.SetLogCodeFilename(true).SetLogCodeLineNumber(false)
	l.SetEnable(true)
	logAllType(l)
	l.SetEnable(false)
}

func TestLogCodeFuncName(t *testing.T) {
	//t.Skip()
	l := NewCeLogger()
	l.SetLogFilePath("TestLogCodeFuncName.log")
	os.Remove(l.LogFilePath)

	t.Log("SetLogCodeFuncName(true)")
	l.SetLogCodeFuncName(true)
	l.SetEnable(true)
	logAllType(l)
	l.SetEnable(false)

	t.Log("SetLogCodeFuncName(false)")
	l.SetLogCodeFuncName(false)
	l.SetEnable(true)
	logAllType(l)
	l.SetEnable(false)
}

func TestLogDate(t *testing.T) {
	//t.Skip()
	l := NewCeLogger()
	l.SetLogFilePath("TestLogDate.log")
	os.Remove(l.LogFilePath)

	t.Log("SetLogDate(true)")
	l.SetLogDate(true)
	l.SetEnable(true)
	logAllType(l)
	l.SetEnable(false)

	t.Log("SetLogDate(false)")
	l.SetLogDate(false)
	l.SetEnable(true)
	logAllType(l)
	l.SetEnable(false)
}

func TestLogTime(t *testing.T) {
	//t.Skip()
	l := NewCeLogger()
	l.SetLogFilePath("TestLogTime.log")
	os.Remove(l.LogFilePath)

	t.Log("SetLogTime(true)")
	l.SetLogTime(true)
	l.SetEnable(true)
	logAllType(l)
	l.SetEnable(false)

	t.Log("SetLogTime(false)")
	l.SetLogTime(false)
	l.SetEnable(true)
	logAllType(l)
	l.SetEnable(false)
}

func TestLogColor(t *testing.T) {
	//t.Skip()
	l := NewCeLogger()
	l.SetLogFilePath("TestLogColor.log")
	os.Remove(l.LogFilePath)

	t.Log("***** SetLogColor(true)")
	l.SetLogColor(true)
	l.SetEnable(true)
	logAllType(l)
	l.SetEnable(false)

	t.Log("SetLogColor(false)")
	l.SetLogColor(false)
	l.SetEnable(true)
	logAllType(l)
	l.SetEnable(false)
}

func TestWriteFile(t *testing.T) {
	//t.Skip()
	l := NewCeLogger()
	l.SetLogFilePath("TestWriteFile.log")
	os.Remove(l.LogFilePath)

	t.Log("SetWriteFile(true)")
	l.SetWriteFile(true)
	l.SetEnable(true)
	logAllType(l)
	l.SetEnable(false)

	t.Log("SetWriteFile(false)")
	l.SetWriteFile(false)
	l.SetEnable(true)
	logAllType(l)
	l.SetEnable(false)
}

func TestWriteConsole(t *testing.T) {
	//t.Skip()
	l := NewCeLogger()
	l.SetLogFilePath("TestWriteConsole.log")
	os.Remove(l.LogFilePath)

	t.Log("SetWriteConsole(true)")
	l.SetWriteConsole(true)
	l.SetEnable(true)
	logAllType(l)
	l.SetEnable(false)

	t.Log("SetWriteConsole(false)")
	l.SetWriteConsole(false)
	l.SetEnable(true)
	logAllType(l)
	l.SetEnable(false)
}

func TestLogEntryTag(t *testing.T) {
	//t.Skip()
	l := NewCeLogger()
	l.SetLogFilePath("TestLogEntryTag.log")
	os.Remove(l.LogFilePath)

	t.Log("SetLogEntryTag(true)")
	l.SetLogEntryTag(true)
	l.SetEnable(true)
	logAllType(l)
	l.SetEnable(false)

	t.Log("SetLogEntryTag(false)")
	l.SetLogEntryTag(false)
	l.SetEnable(true)
	logAllType(l)
	l.SetEnable(false)
}

func TestLogEnableDisable(t *testing.T) {
	//t.Skip()

	l := NewCeLogger()
	l.SetLogFilePath("TestLogEnableDisable.log")
	os.Remove(l.LogFilePath)

	for i := 0; i < 10; i++ {
		l.SetEnable(true)
		logAllType(l)
		l.SetEnable(false)
	}
}

func TestLogFile(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	l := NewCeLogger()
	l.SetWriteConsole(false)

	for i, num := range []uint{1024 * 10} {
		for j, chanLen := range []uint{1, 4, 1024} {
			m := make(map[bool]time.Duration)
			for k, sync := range []bool{false, true} {
				t.Log("TestLogFile", num, chanLen, sync)

				t0 := time.Now()
				tag := fmt.Sprintf("\nT%d-%d-%d", i, j, k)
				logFileTest(l, tag, num, chanLen, sync)

				t1 := time.Now()
				m[sync] = t1.Sub(t0)
				fmt.Printf("%s cost time = %v\n", tag, m[sync])
				l.SetEnable(false)
			}

			fmt.Printf("\n--- num=%d, chanLen=%d, sync-async = %v\n\n", num, chanLen, m[true]-m[false])
		}
	}
}

func logFileTest(l *CeLogger, tag string, num, chanLen uint, sync bool) {
	l.
		SetChanLen(chanLen).
		SetSyncWriteFile(sync).
		SetLogFilePath(fmt.Sprintf("TestLogFile_%s.log", tag))
	os.Remove(l.LogFilePath)
	l.SetEnable(true)

	str := `
        修改了一下。
        
        package main
        
        import (
        "time"
        "fmt"
        "strings"
        )
        
        const (
        DefaultTimeFormat = "2006-01-02 15:04:05.999999999"
        )
        
        func main() {
        timeStr := time.Now().Format(DefaultTimeFormat)
        timeStr = strings.Replace(timeStr, ".", ",", 1)
        fmt.Printf("TimeNow: %s\n", timeStr)`

	for i := 0; i < int(num); i++ {
		l.Infof(tag, "\n===%4d\n\n,%s\n", i, str)
		//time.Sleep(time.Nanosecond)
	}
}

func logAllType(l *CeLogger) {
	l.Trace("TraceTag", "I am a Trace() test")
	l.Tracef("TraceTag", "I am a Trace(...) test:%s", "Trace something")
	l.Trace("int", 100)
	l.Trace("negative int", -8)
	l.Trace("float", 9.8)
	l.Trace("slice", []int{3, 2, 1})
	l.Trace("map", map[string]int{"a": 1, "b": 2, "c": 3})

	l.Info("InfoTag", "I am a Info() test")
	l.Infof("InfoTag", "I am a Info(...) test:%s", "Info something")

	l.Debug("DebugTag", "I am a Debug() test")
	l.Debugf("DebugTag", "I am a Debug(...) test:%s", "Debug something")

	l.Warn("WarnTag", "I am a Warn() test")
	l.Warnf("WarnTag", "I am a Warn(...) test:%s", "Warn something")

	l.Error("ErrorTag", "I am a Error() test")
	l.Errorf("ErrorTag", "I am a Error(...) test:%s", "Error something")

	l.Panic("PanicTag", "I am a Panic() test")
	l.Panicf("PanicTag", "I am a Panic(...) test:%s", "Panic something")
}
