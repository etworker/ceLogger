package ceLogger

import (
	"encoding/json"
	"testing"
)

const (
	filename = "testConfig.json"
)

func TestLoadAndSaveConfig(t *testing.T) {
	c1 := NewCeLoggerConfig()
	c2 := NewCeLoggerConfig()

	t.Log("Test save/load config")

	c1.SeqIndexWidth = 6
	c1.IsLogCodeFilename = true

	c1.SaveConfigFile(filename)
	c2.LoadConfigFile(filename)

	d1, _ := json.Marshal(c1)
	d2, _ := json.Marshal(c2)
	if string(d1) != string(d2) {
		t.Error("config save and load not same")
	}
}

func TestValidateConfig(t *testing.T) {
	c1 := NewCeLoggerConfig()
	c2 := NewCeLoggerConfig()

	t.Log("Test validate config")

	c1.SeqIndexWidth = 100
	c1.SaveConfigFile(filename)
	c2.LoadConfigFile(filename)

	if c2.SeqIndexWidth == 100 {
		t.Error("config validate failed")
	}

}

func TestUpdateConfig(t *testing.T) {
	c1 := NewCeLoggerConfig()

	t.Log("Test update config")

	if err := c1.UpdateConfigByJson(`{SeqIndexWidth}`); err == nil {
		t.Error("update config with wrong json")
	}
	if err := c1.UpdateConfigByJson(`{"SeqIndexWidth":3,"IsLogCodeFilename":false}`); err != nil {
		t.Error("update config with wrong json")
	}
	if c1.SeqIndexWidth != 3 || c1.IsLogCodeFilename != false {
		t.Error("update config with json failed")
	}
}
