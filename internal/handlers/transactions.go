package handlers

import (
	"github.com/go-pg/pg/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/sanokkk/avito-shop/internal/dto/logic"
	"github.com/sanokkk/avito-shop/internal/models"
)

type TransactionsHandler struct {
	db *pg.DB
}

func NewTransactionsHandler(db *pg.DB) *TransactionsHandler {
	return &TransactionsHandler{
		db: db,
	}
}

func (h *TransactionsHandler) GetInfo(c *fiber.Ctx) error {
	log.Info("Получил запрос на историю транзакций")
	uid := c.Locals("uid").(string)
	userId, err := uuid.Parse(uid)
	if err != nil {
		return RespondWithError(400, "Проблемы валидации", c)
	}

	var coins int
	if err := h.db.Model(&models.User{Id: &userId}).WherePK().Column("coins").Select(&coins); err != nil {
		return RespondWithError(500, "Проблемы с извлечением баланса", c)
	}

	inventory, err := h.getInventory(userId)
	if err != nil {
		return RespondWithError(500, "Проблемы с извлечением инвентаря", c)
	}

	transactions, err := h.getTransactions(userId)
	if err != nil {
		return RespondWithError(500, "Проблемы с извлечением истории", c)
	}

	response := logic.InfoDto{
		Coins:        coins,
		Inventory:    inventory,
		CoinsHistory: transactions,
	}

	return c.Status(200).JSON(&response)
}

func (h *TransactionsHandler) getInventory(uid uuid.UUID) ([]*logic.InventoryItem, error) {
	var results []*logic.InventoryItem

	if err := h.db.
		Model(&models.UserItem{}).
		Where("user_id = ?", uid).
		Column("title", "quantity").Select(&results); err != nil {
		return nil, err
	}

	if results == nil {
		results = make([]*logic.InventoryItem, 0)
	}

	return results, nil
}

func (h *TransactionsHandler) getTransactions(uid uuid.UUID) (*logic.TransactionDto, error) {
	var received []*logic.ReceiveTransaction
	if err := h.db.
		Model(&models.History{}).
		Where("history.to_user_id = ?", &uid).
		Join("JOIN users u ON history.to_user_id=u.id").
		Column("u.username", "history.amount").
		Select(&received); err != nil {
		log.Warn("Ошибка при получении транзакций на переводы мне")
		return nil, err
	}
	if received == nil {
		received = make([]*logic.ReceiveTransaction, 0)
	}

	var sent []*logic.SentTransaction
	if err := h.db.
		Model(&models.History{}).
		Where("history.from_user_id = ?", &uid).
		Join("JOIN users ON history.to_user_id=users.id").
		Column("users.username", "history.amount").
		Select(&sent); err != nil {
		log.Warn("Ошибка при получении транзакций на мои переводы")
		return nil, err
	}
	if sent == nil {
		sent = make([]*logic.SentTransaction, 0)
	}

	return &logic.TransactionDto{
		Received: received,
		Sent:     sent,
	}, nil
}
