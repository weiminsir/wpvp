package conf

import (
	"github.com/BurntSushi/toml"
	"testing"
)

func TestInitConf(t *testing.T) {
	runmodel := "dev"
	if _, err := toml.DecodeFile("./conf_"+runmodel+".toml", &Config); err != nil {
		panic("program argument error")
		return
	}
}
