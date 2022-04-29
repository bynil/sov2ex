package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

var C Config

type Config struct {
	Host             string `env:"HOST" envDefault:"127.0.0.1"`
	Port             int    `env:"PORT" envDefault:"8080"`
	MongoHost        string `env:"MONGO_HOST" envDefault:"127.0.0.1"`
	MongoPort        int    `env:"MONGO_PORT" envDefault:"27017"`
	MongoUser        string `env:"MONGO_USER" envDefault:"v2exuser"`
	MongoPass        string `env:"MONGO_PASS" envDefault:"readwrite"`
	MongoDBName      string `env:"MONGO_DB_NAME" envDefault:"v2ex"`
	Debug            bool   `env:"DEBUG" envDefault:"true"`
	LogStdout        bool   `env:"LOG_STDOUT" envDefault:"true"`
	EnableCORS       bool   `env:"ENABLE_CORS" envDefault:"true"`
	LogDir           string `env:"LOG_DIR" envDefault:"/var/log/sov2ex"`
	ESURL            string `env:"ES_URL" envDefault:"http://127.0.0.1:9200"`
	DisableUserCheck bool   `env:"DISABLE_USER_CHECK" envDefault:"true"`
}

func init() {
	if err := env.Parse(&C); err != nil {
		panic(err)
	}
	if C.Debug {
		C.Host = ""
	}
	fmt.Printf("%+v\n", C)
}
