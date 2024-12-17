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

var tomlData = `[lang.C]
listCommand = "make test TFLAGS=--list"
runCommand = "make test TFLAGS="
dirSeperator = "/"
testSeperator = "/"
icon = ""
skiplines = 1
`

// TODO include some sensible defaults
func GetConfig(dir string) UserConfig {
	f, err := os.ReadFile(dir)
	if err != nil {
		log.Fatal(err)
	}

	// f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	// if err != nil {
	// 	log.Fatalf("error opening file: %v", err)
	// }
	// defer f.Close()

	// log.SetOutput(f)

	var conf UserConfig
	_, err = toml.Decode(string(f), &conf)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(conf)
	return conf
}
