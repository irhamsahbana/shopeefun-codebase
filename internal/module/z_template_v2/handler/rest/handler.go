package handler

import "github.com/gofiber/fiber/v2"

type xxxHandler struct {
}

func NewXXXHandler() *xxxHandler {
	return &xxxHandler{}
}

func (h *xxxHandler) Register(router fiber.Router) {

}
