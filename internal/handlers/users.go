package handlers

import (
	"errors"

	"github.com/go-pg/pg/v10"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/sanokkk/avito-shop/internal/dto/requests"
	"github.com/sanokkk/avito-shop/internal/models"
	"github.com/sanokkk/avito-shop/pkg/hashing"
	"github.com/sanokkk/avito-shop/pkg/tokens"
)

var validate *validator.Validate

type UserHandler struct {
	db *pg.DB
}

func NewUserHandler(db *pg.DB) *UserHandler {
	return &UserHandler{
		db: db,
	}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	log.Info("Получил запрос на регистрацию")

	var request requests.RegisterDto
	if err := c.BodyParser(&request); err != nil {
		return RespondWithError(400, "Проблемы валидации", c)
	}

	validate = validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(&request); err != nil {
		log.Error(err)

		return RespondWithError(400, "Проблемы валидации", c)
	}

	userToFind := models.User{}

	if err := h.db.Model(&userToFind).Where("username = ?", request.Username).First(); err != nil {
		if !errors.Is(err, pg.ErrNoRows) {
			return RespondWithError(400, "Проблемы с поиском юзера", c)
		}
	}

	if userToFind.Id != nil {
		return RespondWithError(400, "Пользователь с таким юзернеймом уже существует", c)
	}

	salt := hashing.GenerateRandomSalt(8)
	hash := hashing.HashPassword(request.Password, salt)

	newId := uuid.New()
	userToCreate := models.User{
		Id:           &newId,
		Username:     request.Username,
		PasswordHash: hash,
		Salt:         salt,
	}

	if _, err := h.db.Model(&userToCreate).Insert(); err != nil {
		return RespondWithError(500, "Ошибка создания пользователя", c)
	}

	token, err := tokens.CreateToken(request.Username, *userToCreate.Id)
	if err != nil {
		return RespondWithError(500, "Ошибка создания токена", c)
	}

	return c.Status(200).JSON(struct {
		Token string `json:"token"`
	}{Token: token})
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	log.Info("Получил запрос на аутентификацию")

	var request requests.LoginDto
	if err := c.BodyParser(&request); err != nil {
		return RespondWithError(400, "Проблемы валидации", c)
	}

	validate = validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(&request); err != nil {
		log.Error(err)

		return RespondWithError(400, "Проблемы валидации", c)
	}

	userToFind := models.User{}

	if err := h.db.Model(&userToFind).Where("username = ?", request.Username).First(); err != nil {
		return RespondWithError(400, "Неверный запрос", c)
	}

	if userToFind.Id == nil {
		return RespondWithError(500, "Внутренняя ошибка сервера", c)
	}

	salt := userToFind.Salt
	isPasswordCorrect := hashing.DoPasswordsMatch(userToFind.PasswordHash, request.Password, salt)
	if !isPasswordCorrect {
		return RespondWithError(401, "Неавторизован", c)
	}

	token, err := tokens.CreateToken(request.Username, *userToFind.Id)
	if err != nil {
		return RespondWithError(500, "Ошибка создания токена", c)
	}

	return c.Status(200).JSON(struct {
		Token string `json:"token"`
	}{Token: token})
}

func (h *UserHandler) CheckAuth(c *fiber.Ctx) error {
	log.Info("Получил запрос на auth")

	return c.SendStatus(205)
}
