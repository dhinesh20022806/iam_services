package env

import (
	"fmt"

	"github.com/joeshaw/envdecode"
)



type config struct {
	Database struct {
		URL string `env:"DATABASE_URL,required"`
		MaxIdleConn int `env:"DATABASE_MAX_IDLE_CONNS,default=2,strict"`
		MaxOpenConn int `env:"DATABASE_MAX_OPEN_CONNS,default=2,strict"`
	}
}

var Config config

func init() {

	Reload()
}

func Reload() {

	Config = config{}

	err := envdecode.Decode(&Config)

	if err != nil {
		fmt.Printf("error %v \n", err)
		panic("env not loaded")
	}
}