package main

import (
	"context"
	"io/fs"
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
	)

	if err != nil {
		log.Fatalln(err)
	}

	app := fiber.New()
	app.Use(logger.New())

	var viteFS fs.FS
	viteFS, err = vite.Load()
	if err != nil {
		log.Fatalln(err)
	}

	api := app.Group("/api")
	api.Get("/health", handlers.Health)
	api.Use(func(c fiber.Ctx) error { return c.SendStatus(fiber.StatusNotFound) })

	app.Use(static.New("", static.Config{FS: viteFS}))
	app.Use(func(c fiber.Ctx) error { return c.SendFile("index.html", fiber.SendFile{FS: viteFS}) })

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
