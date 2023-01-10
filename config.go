package main

import (
	"gopkg.in/ini.v1"
)

type Config struct {
	ApiAddr   string
	SqliteDsn string
	Inputs    map[string]string
	Outputs   map[string]string
}

func (c *Config) Load(file string) (err error) {
	cfg, err := ini.InsensitiveLoad(file)
	if err != nil {
		return err
	}

	// main section
	cfg.Section("").MapTo(c)

	// inputs section
	c.Inputs = make(map[string]string)
	for inputID, HostChan := range cfg.Section("inputs").KeysHash() {
		c.Inputs[inputID] = HostChan
	}

	// outputs section
	c.Outputs = make(map[string]string)
	for outputID, Chan := range cfg.Section("outputs").KeysHash() {
		c.Outputs[outputID] = Chan
	}

	log.Infow("config", "config", c)

	return nil
}
