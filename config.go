package goproxy

import (
	"os"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {

	ServerAddr string `yaml:"Server-address"`
	ServerPort int    `yaml:"Server-port"`
	BindAddr   string `yaml:"bind-address"`
	BindPort   int    `yaml:"bind-port"`

}

func NewConfig() Config  {

	if _, err := os.Stat("./config.yml"); os.IsNotExist(err) {
		o, _ := yaml.Marshal(Config{
			ServerAddr: "0.0.0.0",
			ServerPort: 19132,
			BindAddr: "0.0.0.0",
			BindPort: 19133,

		})
		f, err := os.OpenFile("./config.yml", os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		_, err = f.WriteString(string(o))
		if err != nil {
			panic(err)
		}
		f.Close()
	}

	var config Config
	yaml2, _ := ioutil.ReadFile("./config.yml")
	yaml.Unmarshal(yaml2, &config)

	return config
}
