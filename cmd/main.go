package main

import (
	"fmt"
	"github.com/OniGbemiga/simple-bitcoin-wallet/internals"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	router := fiber.New()

	//register routers
	router.Use(logger.New())
	router.Mount("/bitcoin", internals.RegisterHttpHandlers())

	err := router.Listen(fmt.Sprintf(":%v", "9090"))
	if err != nil {
		fmt.Println(err)
		return
	}
}
