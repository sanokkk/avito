package app

import (
	"log"

	"github.com/go-pg/pg/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sanokkk/avito-shop/internal/config"
	"github.com/sanokkk/avito-shop/internal/handlers"
	"github.com/sanokkk/avito-shop/internal/middleware"
	"github.com/sanokkk/avito-shop/internal/migration"
)

type UserHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}

type TransactionsHandler interface {
	GetInfo(c *fiber.Ctx) error
}

type SendCoinHandler interface {
	SendCoin(c *fiber.Ctx) error
}

type ItemsHandler interface {
	Buy(c *fiber.Ctx) error
}

type Server struct {
	userHandler         UserHandler
	transactionsHandler TransactionsHandler
	sendCoinHandler     SendCoinHandler
	itemsHandler        ItemsHandler
	server              *fiber.App
}

func CreateServer(db *pg.DB) *Server {

	server := fiber.New(fiber.Config{
		Immutable: true,
	})

	return &Server{
		server:              server,
		userHandler:         handlers.NewUserHandler(db),
		transactionsHandler: handlers.NewTransactionsHandler(db),
		sendCoinHandler:     handlers.NewSendCoinHandler(db),
		itemsHandler:        handlers.NewItemsHandler(db),
	}
}

func (s *Server) Start() {
	cfg := config.MustLoad()

	if cfg == nil {
		log.Fatal("Не получилось")
	}

	migration.MustMigrate()

	s.applyRoutes()
	s.server.Listen(":8080")
}

func (s *Server) applyRoutes() {
	apiRouter := fiber.New()

	apiRouter.Post("/register", s.userHandler.Register)
	apiRouter.Post("/auth", s.userHandler.Login)
	apiRouter.Get("/info", middleware.AuthMiddleware, s.transactionsHandler.GetInfo)
	apiRouter.Post("/sendCoin", middleware.AuthMiddleware, s.sendCoinHandler.SendCoin)
	apiRouter.Get("/buy/:item", middleware.AuthMiddleware, s.itemsHandler.Buy)
	s.server.Mount("/api", apiRouter)
}
