package lib

import "github.com/BurntSushi/toml"

type Config struct {
	Main    MainConfig     `toml:"main"`
	Headers []HeaderConfig `toml:"headers"`
	Ledisdb LedisdbConfig  `toml:"ledisdb"`
}

type MainConfig struct {
	Protocol string `toml:"protocol"`
	Framed   bool   `toml:"framed"`
	Bufferd  bool   `toml:"buffered"`
	Addr     string `toml:"addr"`
	Secure   bool   `toml:"false"`
}

type HeaderConfig struct {
	Name  string `toml:"name"`
	Value string `toml:"value"`
}

type LedisdbConfig struct {
	Datadir string `toml:"datadir"`
}

// DecodeConfigToml ...
func DecodeConfigToml(tomlfile string) (Config, error) {
	var config Config
	_, err := toml.DecodeFile(tomlfile, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
