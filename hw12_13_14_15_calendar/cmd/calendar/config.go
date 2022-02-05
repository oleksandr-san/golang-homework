package main

import (
	"strconv"

	"github.com/spf13/viper"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger  LoggerConf  `json:"logger"`
	Server  ServerConf  `json:"server"`
	Storage StorageConf `json:"storage"`
}

type LoggerConf struct {
	Level string `json:"level"`
	File  string `json:"file,omitempty"`
}

type ServerConf struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type StorageConf struct {
	Driver string `json:"driver"`
	Dsn    string `json:"dsn"`
}

func (c *ServerConf) Address() string {
	return c.Host + ":" + strconv.Itoa(c.Port)
}

func ReadConfig(path string) (*Config, error) {
	viper.SetConfigType("json")
	viper.SetConfigFile(path)
	viper.SetEnvPrefix("calendar")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := Config{}
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
