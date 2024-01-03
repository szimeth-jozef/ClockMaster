package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"szimeth-jozef/clockmaster/routes"
	"szimeth-jozef/clockmaster/services/db"
	"szimeth-jozef/clockmaster/services/env"
	"szimeth-jozef/clockmaster/services/workitem"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	env, err := env.New(".env")
	if err != nil {
		log.Fatal(err)
	}

	db, err := db.Connect("db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}

	app := echo.New()
	app.HideBanner = true

	app.Use(middleware.CORS())
	app.Static("/", "client/dist")

	// Register routes
	api := app.Group("/api")
	routes.AddWorkItemRoutes(api, db)

	app.HTTPErrorHandler = func(err error, c echo.Context) {
		c.Logger().Error(err)
		log.Println(err)
		c.Redirect(302, "/")
	}

	go func() {
		if err := app.Start(":" + env.Port); err != nil && err != http.ErrServerClosed {
			app.Logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		app.Logger.Fatal(err)
	}

	log.Println("Gracefully shutting down...")
	workItemService := workitem.WorkItemService{DB: db}
	stoppedWorkItems, err := workItemService.Stop()
	if err != nil {
		log.Println(err)
	}
	log.Println("Stopped work items count:", len(stoppedWorkItems))
}
