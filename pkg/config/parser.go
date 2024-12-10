package conf

type EngineConfig struct {
	ListCommand   string
	RunCommand    string
	DirSeperator  string
	TestSeperator string
	Icon          string
	SkipLines     int
}

func GetConfig(dir string) []EngineConfig {
	return []EngineConfig{
		{
			"make test TFLAGS=--list",
			"make test TFLAGS=",
			"/",
			"/",
			"banana",
			1,
		},
	}
}
