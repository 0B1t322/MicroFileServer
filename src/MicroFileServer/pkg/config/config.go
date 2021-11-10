package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"github.com/joho/godotenv"
)

type Config struct {
	DB   *DBConfig   		`json:"DbOptions"`
	Auth *AuthConfig 		`json:"AuthOptions"`
	App  *AppConfig  		`json:"AppOptions"`
	AMQP *AMQPConfig		`json:"AMQPOptions"`
}

type AMQPConfig	struct {
	AMQPURI			string	`envconfig:"MFS_AMQP_URI" json:"amqpUri"`
}

type DBConfig struct {
	URI           	string 	`envconfig:"MFS_MONGO_URI" json:"uri"`
}

type AuthConfig struct {
	KeyURL   		string 	`envconfig:"MFS_AUTH_KEY_URL" json:"keyUrl"`
	Audience 		string 	`default:"itlab" envconfig:"MFS_AUTH_AUDIENCE" json:"audience"`
	Issuer   		string 	`envconfig:"MFS_AUTH_ISSUER" json:"issuer"`
	*RolesConfig			`json:"roles"`
}

type RolesConfig struct {
	UserRole		string	`default:"user" envconfig:"MFS_AUTH_ROLE_USER" json:"user"`
	AdminRole		string	`default:"mfs.admin" envconfig:"MFS_AUTH_ROLE_ADMIN" json:"admin"`
}

type AppConfig struct {
	AppPort  		string 	`envconfig:"MFS_APP_PORT" json:"appPort"`
	TestMode		bool	`envconfig:"MFS_APP_TEST_MODE" json:"testMode"`
	MaxFileSize		int64	`envconfig:"MFS_APP_MAX_FILE_SIZE" json:"maxFileSize"`
}

func GetConfig() *Config {
	var config Config

	if err := godotenv.Load("./.env"); err != nil {
		log.Warn("Don't find .env file")
	}

	data, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.WithFields(log.Fields{
			"function": "GetConfig.ReadFile",
			"error":    err,
		},
		).Warn("Can't read config.json file")
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.WithFields(log.Fields{
			"function": "GetConfig.Unmarshal",
			"error":    err,
		},
		).Warn("Can't correctly parse json from config.json")
	}

	data, err = ioutil.ReadFile("./auth_config.json")
	if err != nil {
		log.WithFields(log.Fields{
			"function": "GetConfig.ReadFile",
			"error":    err,
		},
		).Warn("Can't read auth_config.json file")
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.WithFields(log.Fields{
			"function": "GetConfig.Unmarshal",
			"error":    err,
		},
		).Warning("Can't correctly parse json from auth_config.json")
	}

	err = envconfig.Process("mfs", &config)
	if err != nil {
		log.WithFields(log.Fields{
			"function": "envconfig.Process",
			"error":    err,
		},
		).Fatal("Can't read env vars, shutting down...")
	}
	return &config
}
