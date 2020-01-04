package config

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Backend string    `yaml:"backend"`
	Port    string    `yaml:"port"`
	Redis   redisConf `yaml:"redis"`
	Mongo   mongoConf `yaml:"mongo"`
}

type redisConf struct {
	RedisURL string `yaml:"redis_url"`
}

type mongoConf struct {
	MongoURL     string `yaml:"mongo_url"`
	MongoDB      string `yaml:"mongo_db"`
	MongoTimeout int    `yaml:"mongo_timeout"`
	CacheEnabled bool   `yaml:"cache_enabled"`
}

// Default config
var config = Configuration{
	Backend: "redis",
	Port:    ":8000",
	Redis: redisConf{
		RedisURL: "redis://localhost:6379",
	},
	Mongo: mongoConf{
		MongoURL:     "mongodb://localhost/shortener",
		MongoDB:      "shortener",
		MongoTimeout: 5,
		CacheEnabled: true,
	},
}

func SetConfigFromFile() error {
	file, err := ioutil.ReadFile("config/config.yaml")
	if err == nil {
		if err = yaml.Unmarshal(file, &config); err != nil {
			return errors.Wrap(err, "Failed to unmarshal config file")
		}
	} else if !os.IsNotExist(err) {
		return errors.Wrap(err, "Failed to read config file")
	} else {
		log.Println("No config file found. Proceed with default config")
	}
	log.Printf("Loaded configuration: %+v\n", config)
	return nil
}

func GetConfig() Configuration {
	return config
}
