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

func (h *UserHandler) Auth(c *fiber.Ctx) error {
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
		return h.login(c, &userToFind, request.Password)
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

func (h *UserHandler) login(c *fiber.Ctx, userToFind *models.User, password string) error {
	log.Info("начинаю аутентификацию")

	salt := userToFind.Salt
	isPasswordCorrect := hashing.DoPasswordsMatch(userToFind.PasswordHash, password, salt)
	if !isPasswordCorrect {
		return RespondWithError(401, "Неавторизован", c)
	}

	token, err := tokens.CreateToken(userToFind.Username, *userToFind.Id)
	if err != nil {
		return RespondWithError(500, "Ошибка создания токена", c)
	}

	return c.Status(200).JSON(struct {
		Token string `json:"token"`
	}{Token: token})
}
