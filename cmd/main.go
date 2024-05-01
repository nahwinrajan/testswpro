package main

import (
	"os"

	"github.com/nahwinrajan/testswpro/generated"
	"github.com/nahwinrajan/testswpro/handler"
	"github.com/nahwinrajan/testswpro/repository"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	server := newServer()

	generated.RegisterHandlers(e, server)
	e.Use(middleware.Logger())
	// TODO: ideally we want to add configuration for
	// cors (unless we are only accessible from within cluster)
	// metrics gathering
	// http knobs (read, write, idle timeout, etc)

	// TODO GRACEFUL SHUTDOWN
	e.Logger.Fatal(e.Start(":1323"))
}

func newServer() *handler.Server {
	dbDsn := os.Getenv("DATABASE_URL")
	repo := repository.New(dbDsn)

	return handler.New(repo)
}
