package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/irnafitriani/music/entity"
	"github.com/irnafitriani/music/handler"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	dsn := "root:@tcp(127.0.0.1:3306)/music_app?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&entity.Song{})

	app := fiber.New(
		fiber.Config{
			BodyLimit: 50 * 1024 * 1024, // limit MB
		},
	)

	songHandler := handler.NewSongHandler(db)

	app.Get("/", handler.HelloHandler)

	app.Post("/song", songHandler.Add)
	app.Get("/song", songHandler.List)
	app.Put("/song/:id", songHandler.Update)
	app.Delete("/song/:id", songHandler.Delete)

	app.Post("song/upload", songHandler.Upload)
	app.Listen(":4000")
}
