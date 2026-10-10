package env



type config struct {
	Database struct {
		URL string `env:"DATABASE_URL,required"`
		MaxIdleConn int `env:"DATABASE_MAX_IDLE_CONNS,default=2,strict"`
		MaxOpenConn int `env:"DATABASE_MAX_OPEN_CONNS,default=2,strict"`
	}
}

var Config config