package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	App struct {
		BaseURL string
		Port    string
		ENV     string
	}
	Database struct {
		Username string
		Password string
		Host     string
		Port     string
		Name     string
	}
	AwsS3 struct {
		Bucket    string
		Region    string
		AccessKey string
		SecretKey string
	}
}

var appConfig *AppConfig

func Get() *AppConfig {
	if appConfig == nil {
		appConfig = initConfig()
	}
	return appConfig
}

func initConfig() *AppConfig {

	config := AppConfig{}

	// Load .env file, set default if fail
	err := godotenv.Load()
	if err != nil {
		config.App.Port = "8000"
		config.App.BaseURL = "localhost:" + config.App.Port
		config.App.ENV = ""
		config.Database.Host = "localhost"
		config.Database.Port = "3306"
		config.Database.Username = "root"
		config.Database.Password = "root"
		config.Database.Name = "bringeee"
		config.AwsS3.Bucket = ""
		config.AwsS3.Region = ""
		config.AwsS3.AccessKey = ""
		config.AwsS3.SecretKey = ""

		return &config
	}

	// set config based on .env
	config.App.Port = os.Getenv("APP_PORT")
	config.App.BaseURL = os.Getenv("APP_BASE_URL")
	config.App.ENV = os.Getenv("")
	config.Database.Host = os.Getenv("DB_HOST")
	config.Database.Port = os.Getenv("DB_PORT")
	config.Database.Username = os.Getenv("DB_USERNAME")
	config.Database.Password = os.Getenv("DB_PASSWORD")
	config.Database.Name = os.Getenv("DB_NAME")
	config.AwsS3.Bucket = os.Getenv("AWS_S3_BUCKET")
	config.AwsS3.Region = os.Getenv("AWS_S3_REGION")
	config.AwsS3.AccessKey = os.Getenv("AWS_S3_ACCESS_KEY")
	config.AwsS3.SecretKey = os.Getenv("AWS_S3_SECRET_KEY")

	return &config
}
