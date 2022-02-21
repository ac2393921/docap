package config

func GetPlatformDefaultConfig() OSConfig {
	return OSConfig{
		OpenCommand:     "open {{filename}}",
		OpenLickCommand: "open {{link}}",
	}
}
