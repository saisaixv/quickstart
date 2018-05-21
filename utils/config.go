package utils

import(
	"fmt"

	"github.com/go-ini/ini"
)

const(
	CFG_SEP=","
)

var(
	config *Config
)

type Config struct {
	profile string
}

func NewConfig(profile string) *Config  {
	o:=&Config{
		profile,
	}
	return o
}

func (self *Config)Load() (*ini.File,error) {
	
	if FileIsExist(self.profile)==false{
		return nil,fmt.Errorf("%s:%s","file not exist",self.profile)
	}
	return ini.LooseLoad(self.profile)
}

func MakeDefaultConfig(cfg *Config)  {
	if cfg!=nil{
		config=cfg
	}
}

func DefaultConfig() *Config {
	return config
}

