package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gofiber/fiber/v2"
	"github.com/irnafitriani/music/config"
	"github.com/irnafitriani/music/entity"
	"github.com/irnafitriani/music/handler"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	conf := config.LoadConfig()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.DB.User,
		conf.DB.Password,
		conf.DB.Host,
		conf.DB.Port,
		conf.DB.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&entity.Song{})
	db.AutoMigrate(&entity.User{})

	app := fiber.New(
		fiber.Config{
			BodyLimit: 50 * 1024 * 1024, // limit MB
		},
	)

	s3session := Creates3Session(conf.S3)
	songHandler := handler.NewSongHandler(db, conf, s3session)
	userHadler := handler.NewUserHandler(db, conf)

	app.Get("/", handler.HelloHandler)

	app.Post("/song", songHandler.Add)
	app.Get("/song", songHandler.List)
	app.Put("/song/:id", songHandler.Update)
	app.Delete("/song/:id", songHandler.Delete)

	app.Post("song/upload", songHandler.Upload)
	app.Get("/stream/:id", songHandler.Play)

	app.Post("/register", userHadler.Register)
	app.Listen(":4000")
}

func Creates3Session(s3Conf config.S3) *session.Session {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(s3Conf.Region),
		Credentials: credentials.NewStaticCredentials(
			s3Conf.Key,
			s3Conf.Secret,
			""),
	})

	if err != nil {
		panic(err)
	}

	return sess
}
