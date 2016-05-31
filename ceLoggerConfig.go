package ceLogger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func init() {
	fmt.Println("Init ceLoggerConfig")
}

// ----------
// EntryConfig
// ----------
const (
	ECTrace = "Trace"
	ECInfo  = "Info"
	ECDebug = "Debug"
	ECWarn  = "Warn"
	ECError = "Error"
	ECPanic = "Panic"
)

type EntryConfig struct {
	Tag            string // e.g. "P" -> "Panic"
	IsEnable       bool
	IsWriteFile    bool // if log to file
	IsWriteConsole bool // if log to console
	IsColorFile    bool // if log color in file
	IsColorConsole bool // if log color in console
	DisplayMode    uint
	ForeColor      uint
	BackColor      uint
}
type EntryConfigMap map[string]*EntryConfig

func (ec *EntryConfig) ValidateConfig() *EntryConfig {
	switch ec.DisplayMode {
	case 0, 1, 4, 5, 7, 8:
	default:
		ec.DisplayMode = 0
	}

	if ec.ForeColor < 30 || ec.ForeColor > 37 {
		ec.ForeColor = 37
	}

	if ec.BackColor < 40 || ec.BackColor > 47 {
		ec.BackColor = 0
	}

	return ec
}

// ----------
// CeLoggerConfig
// ----------

type CeLoggerConfig struct {
	ChanLen             uint           // Chan buffer len
	MaxFileSize         uint           // max file size of one log file
	MaxEntryNum         uint           // max entry num in one log file
	ContentDelimiter    string         // default is " ", set "\n" will print content at next line
	IsSyncWriteFile     bool           // sync write file or async by goroutine()
	IsLogOrderFlag      bool           // if log "X" when out of order, latter will happen when write file async
	IsLogSeqIndex       bool           // if log entry index
	SeqIndexWidth       uint           // width of entry index, e.g. =4 means -> 0001 - 9999
	IsLogEntryTag       bool           // if log type tag, = T/I/D/W/E/P, means Trace/Info/Debug/Warn/Error/Panic
	IsLogFuncEnterExit  bool           // if log func enter/exit
	IsLogCodeFilename   bool           // if log current filename in code
	IsLogCodeLineNumber bool           // if log current line number in code
	IsLogCodeFuncName   bool           // if log current func name in code
	IsLogDate           bool           // if log date, e.g. 2015-03-04
	IsLogTime           bool           // if log time, e.g. 12:34.23456
	TimeMsWidth         uint           // width of ms width, e.g. =4 means time -> 12:34.1231
	IsLogColor          bool           // if colorize content
	IsWriteFile         bool           // if log to file
	IsWriteConsole      bool           // if log to console
	LogFilePath         string         // log filename
	ECMap               EntryConfigMap // store all log type info, e.g. Trace/Info/Debug/Warn/Error/Panic
}

func NewCeLoggerConfig() *CeLoggerConfig {
	c := &CeLoggerConfig{}

	c.ChanLen = 1024                // 1K
	c.MaxFileSize = 1 * 1024 * 1024 // 1MB
	c.MaxEntryNum = 10 * 1024       // 10K
	c.ContentDelimiter = " "        // "\n" -> put content to new line
	c.IsSyncWriteFile = true        // Sync write file is safe and in order
	c.IsLogOrderFlag = false        // Only need when async write file
	c.IsLogSeqIndex = true          // Entry index
	c.SeqIndexWidth = 4             // Entry index just like "0023"
	c.IsLogFuncEnterExit = true
	c.IsLogCodeFilename = false
	c.IsLogCodeLineNumber = false
	c.IsLogCodeFuncName = true
	c.IsLogDate = false
	c.IsLogTime = true
	c.TimeMsWidth = 4
	c.IsLogColor = true
	c.IsWriteFile = true
	c.IsWriteConsole = true
	c.IsLogEntryTag = true
	c.LogFilePath = ""

	c.ECMap = make(EntryConfigMap)
	c.ECMap[""] = &EntryConfig{Tag: "", DisplayMode: 0, ForeColor: 33, BackColor: 0}
	c.ECMap[ECTrace] = &EntryConfig{Tag: "T", DisplayMode: 0, ForeColor: 33, BackColor: 0}
	c.ECMap[ECInfo] = &EntryConfig{Tag: "I", DisplayMode: 1, ForeColor: 37, BackColor: 0}
	c.ECMap[ECDebug] = &EntryConfig{Tag: "D", DisplayMode: 0, ForeColor: 36, BackColor: 0}
	c.ECMap[ECWarn] = &EntryConfig{Tag: "W", DisplayMode: 1, ForeColor: 31, BackColor: 43}
	c.ECMap[ECError] = &EntryConfig{Tag: "E", DisplayMode: 1, ForeColor: 37, BackColor: 41}
	c.ECMap[ECPanic] = &EntryConfig{Tag: "P", DisplayMode: 1, ForeColor: 33, BackColor: 41}

	for _, ec := range c.ECMap {
		ec.IsEnable = true
		ec.IsColorConsole = true
		ec.IsColorFile = false
		ec.IsWriteConsole = true
		ec.IsWriteFile = true
	}
    
	return c
}

// Update config with json string
func (c *CeLoggerConfig) UpdateConfigByJson(js string) error {
	if err := json.Unmarshal([]byte(js), c); err != nil {
		fmt.Printf("Parse log config Json string [%s] failed: %s\n", js, err.Error())
		return err
	}

	c.ValidateConfig()

	return nil
}

func (c *CeLoggerConfig) LoadConfigFile(filePath string) error {
	filename := "logConfig.json"

	// Default config filename is "logConfig.json"
	if filePath != "" {
		filename = filePath
	}

	// Read config file
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Read log config file %s failed: %s\n", filename, err.Error())
		return err
	}

	// Parse json to config
	return c.UpdateConfigByJson(string(dat))
}

// Write config to log config file
func (c *CeLoggerConfig) SaveConfigFile(filename string) error {
	c.ValidateConfig()

	dat, _ := json.Marshal(c)
	err := ioutil.WriteFile(filename, dat, 0644)
	if err != nil {
		fmt.Printf("Write log config file %s failed: %s\n", filename, err.Error())
		return err
	}
	return nil
}

func (c *CeLoggerConfig) SetEntryConfig(name string, ec *EntryConfig) *EntryConfig {
	c.ECMap[name] = ec
	return ec
}

// Validate config data
func (c *CeLoggerConfig) ValidateConfig() *CeLoggerConfig {
	if c.ChanLen == 0 {
		c.ChanLen = 1024
	}

	if c.SeqIndexWidth > 8 {
		c.SeqIndexWidth = 8
	}

	if c.TimeMsWidth > 9 {
		c.TimeMsWidth = 9
	}

	for _, ec := range c.ECMap {
		ec.ValidateConfig()
	}
	return c
}
