package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/irnafitriani/music/config"
	"github.com/irnafitriani/music/entity"
	"github.com/irnafitriani/music/helper"
	"gorm.io/gorm"
)

type UserHandler struct {
	db   *gorm.DB
	conf config.Config
}

func NewUserHandler(db *gorm.DB, conf config.Config) *UserHandler {
	return &UserHandler{
		db:   db.Debug(),
		conf: conf,
	}
}

func (h UserHandler) Register(c *fiber.Ctx) error {
	var user entity.User
	c.BodyParser(&user)

	e := helper.Validate(&user)

	if len(e) > 0 {
		return c.Status(400).JSON(e)
	}
	err := h.db.Create(&user).Error
	if err != nil {
		return c.Status(500).JSON(map[string]string{"message": err.Error()})
	}
	return c.JSON(user)
}
