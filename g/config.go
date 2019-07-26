// Package g Global variables and configurations.
package g

import (
	"encoding/json"
	"io/ioutil"
)

type configs struct {
	Debug      bool   `json:"debug"`
	ListenHTTP string `json:"listenHTTP"`
	Logpath    string `json:"logpath"`
	Prof       bool   `json:"prof"`
	ProfHTTP   string `json:"profHTTP"`

	Pidpath string `json:"pidpath"`
}

var (
	// Cfgs global configuration struct.
	Cfgs *configs
)

// ParseConfig read and parse global configuration from a JSON file.
func ParseConfig(path string) error {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	cfgsNew := configs{}
	err = json.Unmarshal(buf, &cfgsNew)
	if err != nil {
		return err
	}

	Cfgs = &cfgsNew
	return nil
}