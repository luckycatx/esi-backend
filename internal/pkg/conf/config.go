package conf

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server *Server
	Mysql  *Mysql
	Redis  *Redis
	Token  *Token
}

type Server struct {
	Host       string
	Port       string
	CtxTimeout int `mapstructure:"ctx_timeout"`
}

type Mysql struct {
	Host   string
	Port   string
	User   string
	Pwd    string
	DBName string `mapstructure:"db_name"`
}

type Redis struct {
	Addr string
	Pwd  string
}

type Token struct {
	AccessSecret  string `mapstructure:"access_secret"`
	AccessExp     int    `mapstructure:"access_exp"`
	RefreshSecret string `mapstructure:"refresh_secret"`
	RefreshExp    int    `mapstructure:"refresh_exp"`
}

func Load() *Config {
	var cfg = &Config{}
	_ = godotenv.Load(".env")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	// Path for dev usage
	viper.AddConfigPath("../config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config file: ", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal("Error unmarshalling config: ", err)
	}

	if os.Getenv("APP_ENV") == "dev" {
		log.Println("App is running in development mode")
	}

	return cfg
}
