package conf

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type EngineConfig struct {
	ListCommand   string
	RunCommand    string
	DirSeperator  string
	TestSeperator string
	Icon          string
	SkipLines     int
}

type UserConfig struct {
	ClientInfo map[string]EngineConfig `toml:"lang"`
}

func GetConfig(dir string) UserConfig {
	f, err := os.ReadFile(dir)
	if err != nil {
		log.Fatal(err)
	}

	var conf UserConfig
	_, err = toml.Decode(string(f), &conf)
	if err != nil {
		log.Fatal(err)
	}
	return conf
}
