package handlers

import (
	"errors"
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/sanokkk/avito-shop/internal/models"
)

type ItemsHandler struct {
	db *pg.DB
}

func NewItemsHandler(db *pg.DB) *ItemsHandler {
	return &ItemsHandler{
		db: db,
	}
}

func (h *ItemsHandler) Buy(c *fiber.Ctx) error {
	log.Info("Получил запрос на покупку товара")

	title := c.Params("item", "")
	if title == "" {
		return RespondWithError(400, "Неверный запрос", c)
	}

	uid, _ := uuid.Parse(c.Locals("uid").(string))
	userCoins, err := h.getUserBalance(uid)
	if err != nil {
		return RespondWithError(500, "внутренняя ошибка сервера", c)
	}

	item, err := h.getItem(title)
	if err != nil {
		return RespondWithError(500, "внутренняя ошибка сервера", c)
	}

	if item.Cost > userCoins {
		log.Warn(fmt.Sprintf("Пользователю %s не хватает денег на %s", uid, title))
		return RespondWithError(400, "неверный запрос", c)
	}

	if err := h.buyItem(item, uid); err != nil {
		return RespondWithError(500, "внутренняя ошибка сервера", c)
	}

	return c.SendStatus(200)
}

func (h *ItemsHandler) getItem(title string) (*models.Item, error) {
	var item models.Item
	if err := h.db.
		Model(&models.Item{}).
		Where("title = ?", &title).
		For("UPDATE").
		Limit(1).
		Select(&item); err != nil {
		log.Warn("ошибка при получении товара ", err)
		return nil, errors.New("ошибка при получении товара")
	}

	return &item, nil
}

func (h *ItemsHandler) getUserBalance(uid uuid.UUID) (int, error) {
	var userCoins int
	if err := h.db.
		Model(&models.User{Id: &uid}).
		WherePK().
		Column("coins").
		For("UPDATE").Select(&userCoins); err != nil {
		return -1, errors.New("ошибка при вычислении баланса")
	}

	return userCoins, nil
}

func (h *ItemsHandler) buyItem(item *models.Item, uid uuid.UUID) error {
	tx, err := h.db.Begin()
	if err != nil {
		return errors.New("ошибка при создании транзакции")
	}

	defer tx.Close()

	if _, err := tx.
		Model(&models.User{Id: &uid}).
		Set("coins = coins - ?", item.Cost).WherePK().Update(); err != nil {
		log.Warn("Ошибка при обновлении баланса: ", err)
		tx.Rollback()
		return errors.New("внутренняя ошибка сервера")
	}

	var userItem models.UserItem
	if err := tx.
		Model(&models.UserItem{}).
		Where("item_id = ? AND user_id = ?", item.Id, uid).
		Select(&userItem); err != nil {
		if !errors.Is(err, pg.ErrNoRows) {
			log.Warn("Ошибка при поиске товара юзера: ", err)
			tx.Rollback()
			return errors.New("внутренняя ошибка сервера")
		}
	}

	if userItem.Id == uuid.Nil {
		toInsert := models.UserItem{
			ItemId:   item.Id,
			Title:    item.Title,
			Quantity: 1,
			UserId:   uid,
		}
		log.Debug("Товара у юзера нет, создаю")
		if _, err := tx.Model(&toInsert).Insert(); err != nil {
			tx.Rollback()
			log.Warn("Ошибка при вставке товара юзера: ", err)
			return errors.New("внутренняя ошибка сервера")
		}
	} else {
		if _, err := tx.
			Model(&models.UserItem{}).
			Set("quantity = quantity + 1").
			Where("id = ?", &userItem.Id).
			Update(); err != nil {
			tx.Rollback()
			log.Warn("Ошибка при обновлении товара юзера: ", err)
			return errors.New("внутренняя ошибка сервера")
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		log.Warn("Ошибка при сохранении транзакции: ", err)
		return errors.New("внутренняя ошибка сервера")
	}

	return nil
}
