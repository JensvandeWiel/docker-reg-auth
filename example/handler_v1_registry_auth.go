package main

import (
	registry "github.com/JensvandeWiel/docker-reg-auth"
	"github.com/labstack/echo/v4"
	"log/slog"
)

type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Comment string `json:"comment,omitempty"`
}

type RegistryAuthHandler struct {
	authenticator  registry.Authenticator
	authorizer     registry.Authorizer
	tokenGenerator registry.TokenGenerator
}

func NewRegistryAuthHandler(authenticator registry.Authenticator, authorizer registry.Authorizer, tokenGenerator registry.TokenGenerator) *RegistryAuthHandler {
	return &RegistryAuthHandler{
		authenticator:  authenticator,
		authorizer:     authorizer,
		tokenGenerator: tokenGenerator,
	}
}

func (h *RegistryAuthHandler) AuthHandle(ctx echo.Context) error {
	usr, passwd, ok := ctx.Request().BasicAuth()
	if !ok {
		return ctx.JSON(401, HttpError{
			Code:    401,
			Message: "Unauthorized",
		})
	}

	if err := h.authenticator.Authenticate(usr, passwd); err != nil {
		return ctx.JSON(401, HttpError{
			Code:    401,
			Message: "Unauthorized",
			Comment: "Invalid username or password",
		})
	}

	slog.Info("Authenticated user", slog.String("user", usr))

	req, err := registry.AuthorizationRequestFromContext(ctx)
	if err != nil {
		return ctx.JSON(400, HttpError{
			Code:    400,
			Message: "Bad Request",
			Comment: err.Error(),
		})
	}

	if req.Account != usr {
		return ctx.JSON(400, HttpError{
			Code:    400,
			Message: "Bad Request",
			Comment: "Account does not match authenticated user",
		})
	}

	slog.Info("Parsed authorization request", "request", req)

	actions, err := h.authorizer.Authorize(req)
	if err != nil {
		return ctx.JSON(401, HttpError{
			Code:    401,
			Message: "Forbidden",
			Comment: err.Error(),
		})
	}

	slog.Info("Authorized actions", slog.Any("actions", actions))

	token, err := h.tokenGenerator.GenerateToken(req, actions, &registry.TokenOptions{
		ExpiresIn: 3600,
		Issuer:    "test",
		Audience:  "test",
	})

	if err != nil {
		return ctx.JSON(500, HttpError{
			Code:    500,
			Message: "Internal Server Error",
			Comment: err.Error(),
		})
	}

	if token == nil {
		return ctx.JSON(500, HttpError{
			Code:    500,
			Message: "Internal Server Error",
			Comment: "Failed to generate token",
		})

	}

	return ctx.JSON(200, token)
}
