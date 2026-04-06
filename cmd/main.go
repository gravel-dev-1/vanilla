package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"gbfw/internal/bootstrap"
	"gbfw/internal/env"
	"gbfw/internal/handlers"
	"gbfw/internal/vite"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func main() {
	err := bootstrap.Run(
		env.Load,
		vite.Load,
	)

	if err != nil {
		log.Fatalln(err)
		return
	}

	app := fiber.New()
	app.Use(logger.New())

	api := app.Group("/api")
	api.Get("/health", handlers.Health)
	api.Use(func(c fiber.Ctx) error { return c.SendStatus(fiber.StatusNotFound) })

	app.Use(static.New("", static.Config{FS: vite.FS}))
	app.Use(func(c fiber.Ctx) error { return c.SendFile("/index.html", fiber.SendFile{FS: vite.FS}) })

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	go func() {
		if err := app.Listen(os.Getenv("LISTEN_ADDR"), fiber.ListenConfig{GracefulContext: ctx}); err != nil {
			log.Println(err)
		}
	}()

	<-ctx.Done()

	err = bootstrap.Run(
		app.Shutdown,
		func() (err error) { return app.ShutdownWithContext(ctx) },
	)

	if err != nil {
		log.Fatalln(err)
	}
}
