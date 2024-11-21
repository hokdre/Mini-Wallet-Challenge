package config

import (
	"log"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	// REST SERVER
	RestPORT             string        `envconfig:"REST_PORT"`
	RestReadTimeOut      time.Duration `envconfig:"REST_READ_TIMEOUT"`
	RestWriteTimeOut     time.Duration `envconfig:"REST_WRITE_TIMEOUT"`
	RestShoutDownTimeOut time.Duration `envconfig:"REST_SHUTDOWN_TIMEOUT"`

	// POSTGRE
	PostgreHost        string `envconfig:"POSTGRE_HOST"`
	PostgrePort        string `envconfig:"POSTGRE_PORT"`
	PostgreUsername    string `envconfig:"POSTGRE_USERNAME"`
	PostgrePassword    string `envconfig:"POSTGRE_PASSWORD"`
	PostgreDB          string `envconfig:"POSTGRE_DB"`
	PostgreSSLMode     string `envconfig:"POSTGRE_SSL_MODE"`
	PostgreMaxIdleConn int    `envconfig:"POSTGRE_MAX_IDLE_CONN"`
	PostgreMaxOpenConn int    `envconfig:"POSTGRE_MAX_OPEN_CONN"`

	// TOKEN
	AESSecret string `envconfig:"AES_SECRET"`
}

var config Config
var onceConfig sync.Once

func Init() Config {
	config = Config{}
	onceConfig.Do(func() {
		err := envconfig.Process("", &config)
		if err != nil {
			log.Fatalf("failed to read config : %s", err)
		}
	})

	return config
}
