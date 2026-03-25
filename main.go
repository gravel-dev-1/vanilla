package main

import (
	"context"
	"gbfw/api/bootstrap"
	"gbfw/api/controllers"
	"gbfw/api/env"
	"gbfw/api/vite"
	"io/fs"
	"log"
	"os"
	"os/signal"

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
	api.Get("/health", controllers.Health)
	api.Use(func(c fiber.Ctx) error { return c.SendStatus(fiber.StatusNotFound) })

	app.Use(static.New("", static.Config{FS: viteFS}))
	app.Use(func(c fiber.Ctx) error { return c.SendFile("index.html", fiber.SendFile{FS: viteFS}) })

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	go func() {
		if err := app.Listen(":8080", fiber.ListenConfig{GracefulContext: ctx}); err != nil {
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
