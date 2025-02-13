package handlers

import (
	"errors"
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/sanokkk/avito-shop/internal/dto/requests"
	"github.com/sanokkk/avito-shop/internal/models"
)

type SendCoinHandler struct {
	db *pg.DB
}

func NewSendCoinHandler(db *pg.DB) *SendCoinHandler {
	return &SendCoinHandler{
		db: db,
	}
}

func (h *SendCoinHandler) SendCoin(c *fiber.Ctx) error {
	log.Info("Получил запрос на перевод")

	var request requests.SendCoinDto
	if err := c.BodyParser(&request); err != nil {
		return RespondWithError(400, "Неверный запрос", c)
	}

	validate = validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(&request); err != nil {
		log.Warn("Ошибка валидации при переводе", err)
		return RespondWithError(400, "Неверный запрос", c)
	}

	uid, _ := uuid.Parse(c.Locals("uid").(string))
	if err := h.checkUserBalanceForOperation(uid, request.Amount); err != nil {
		return RespondWithError(500, err.Error(), c)
	}

	if err := h.sendCoins(uid, request.ToUser, request.Amount); err != nil {
		RespondWithError(500, err.Error(), c)
	}

	return c.SendStatus(200)
}

func (h *SendCoinHandler) checkUserBalanceForOperation(uid uuid.UUID, coinsToSend int) error {
	var userCoins int
	if err := h.db.
		Model(&models.User{Id: &uid}).
		WherePK().
		Column("coins").
		For("UPDATE").Select(&userCoins); err != nil {
		return errors.New("ошибка при вычислении баланса")
	}

	if userCoins < coinsToSend {
		log.Warn(fmt.Sprintf("У пользователя %s недостаточно средств", uid))
		return errors.New("недостаточно средств")
	}

	return nil
}

func (h *SendCoinHandler) sendCoins(fromUser uuid.UUID, toUser string, coins int) error {
	tx, err := h.db.Begin()
	if err != nil {
		log.Warn("Ошибка старта транзакции", err)
		return errors.New("внутренняя ошибка сервера")
	}

	defer tx.Close()

	var toUserId uuid.UUID
	if err := tx.
		Model(&models.User{}).
		Where("username = ?", toUser).Column("id").Select(&toUserId); err != nil {
		log.Warn("Откатываю транзакцию, пользователя нет")
		tx.Rollback()

		return errors.New("внутренняя ошибка сервера")
	}

	if _, err = tx.
		Model(&models.User{}).
		Set("coins = coins - ?", coins).
		Where("id = ?", &fromUser).Update(); err != nil {
		log.Warn("Откатываю транзакцию, не получилось вычесть монеты", err)
		tx.Rollback()

		return errors.New("внутренняя ошибка сервера")
	}

	if _, err = tx.
		Model(&models.User{}).
		Set("coins = coins + ?", coins).
		Where("id = ?", &toUserId).Update(); err != nil {
		log.Warn("Откатываю транзакцию, не получилось прибавить монеты", err)
		tx.Rollback()

		return errors.New("внутренняя ошибка сервера")
	}

	//newId := uuid.New()
	toInsert := models.History{
		//Id:         &newId,
		FromUserId: fromUser,
		ToUserId:   toUserId,
		Amount:     coins,
	}

	if _, err := tx.Model(&toInsert).Insert(); err != nil {
		log.Warn("Ошибка при вставке транзакции", err)
		tx.Rollback()
		return errors.New("внутренняя ошибка сервера")
	}

	if err := tx.Commit(); err != nil {
		log.Warn("Не получилось закоммитить транзакцию")
		tx.Rollback()
	}
	log.Info("Завершил транзакцию")

	return nil
}
