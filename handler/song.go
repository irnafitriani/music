package handler

import (
	"fmt"
	"io"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/irnafitriani/music/config"
	"github.com/irnafitriani/music/entity"
	"github.com/irnafitriani/music/helper"
	"gorm.io/gorm"
)

type SongHandler struct {
	db   *gorm.DB
	conf config.Config
}

func NewSongHandler(db *gorm.DB, conf config.Config) *SongHandler {
	return &SongHandler{
		db:   db.Debug(),
		conf: conf,
	}
}

func (h SongHandler) Add(c *fiber.Ctx) error {
	var song entity.Song
	c.BodyParser(&song)

	e := helper.Validate(&song)

	if len(e) > 0 {

		return c.Status(400).JSON(e)
	}

	if _, err := os.Stat(h.conf.StoragePath + "/" + song.File); err != nil {
		return c.Status(400).JSON(map[string]string{"message": "file not found"})
	}

	h.db.Create(&song)
	return c.JSON(song)
}

func (h SongHandler) Upload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(err)
	}

	err = c.SaveFile(file, fmt.Sprintf("%s/%s", h.conf.StoragePath, file.Filename))

	if err != nil {
		return c.Status(400).JSON(err)
	}

	return c.JSON(file.Filename)
}

func (h SongHandler) List(c *fiber.Ctx) error {
	var songs []entity.Song
	h.db.Find(&songs)
	for i, song := range songs {
		song.StreamUrl = fmt.Sprintf("%s/stream/%d", h.conf.AppUrl, song.ID)
		songs[i] = song
	}
	return c.JSON(songs)
}

func (h SongHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var song entity.Song
	h.db.First(&song, id)
	if song.ID == 0 {
		return c.Status(404).JSON(map[string]string{"message": "song not found"})
	}
	var payload entity.Song
	c.BodyParser(&payload)
	e := helper.Validate(&payload)
	if len(e) > 0 {

		return c.Status(400).JSON(e)
	}

	if _, err := os.Stat(h.conf.StoragePath + "/" + song.File); err != nil {
		return c.Status(400).JSON(map[string]string{"message": "file not found"})
	}

	song.Artist = payload.Artist
	song.Title = payload.Title
	song.Cover = payload.Cover
	song.File = payload.File
	h.db.Updates(&song)
	return c.JSON(song)
}

func (h SongHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	var song entity.Song
	h.db.First(&song, id)
	if song.ID == 0 {
		return c.Status(404).JSON(map[string]string{"message": "song not found"})
	}

	h.db.Delete(&song)
	return c.JSON(map[string]string{"message": "song deleted"})
}

func (h SongHandler) Play(c *fiber.Ctx) error {
	id := c.Params("id")
	var song entity.Song
	h.db.First(&song, id)
	if song.ID == 0 {
		return c.Status(404).JSON(map[string]string{"message": "song not found"})
	}

	filePath := fmt.Sprintf("%s/%s", h.conf.StoragePath, song.File)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()
	c.Set("Content-Type", "audio/mpeg")
	_, err = io.Copy(c, file)
	if err != nil {
		return err
	}
	return nil
}
