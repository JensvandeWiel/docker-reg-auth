package main

import (
	registry "github.com/JensvandeWiel/docker-reg-auth"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	gen, err := registry.NewDefaultTokenGenerator(".devcerts/RootCA.crt", ".devcerts/RootCA.key")
	if err != nil {
		panic(err)
	}

	h := NewRegistryAuthHandler(registry.NewDummyAuthenticator(), registry.NewDummyAuthorizer(), gen)

	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("v1/registry/auth", h.AuthHandle)

	e.Logger.Fatal(e.Start(":8080"))
}
