package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Port        string
	AppUrl      string
	StoragePath string
	DB          DB
	S3          S3
}

type DB struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type S3 struct {
	Key    string
	Secret string
	Region string
	Bucket string
}

func LoadConfig() Config {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return Config{
		AppUrl:      viper.GetString("app_url"),
		Port:        viper.GetString("port"),
		StoragePath: viper.GetString("storage_path"),
		DB: DB{
			Host:     viper.GetString("db.host"),
			Port:     viper.GetString("db.port"),
			User:     viper.GetString("db.user"),
			Password: viper.GetString("db.password"),
			Database: viper.GetString("db.database"),
		},
		S3: S3{
			Key:    viper.GetString("s3.key"),
			Secret: viper.GetString("s3.secret"),
			Region: viper.GetString("s3.region"),
			Bucket: viper.GetString("s3.bucket"),
		},
	}
}
