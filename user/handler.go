package user

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) GetUsers(c *fiber.Ctx) error {
	users, err := h.service.GetUsers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": fmt.Errorf("failed to fetch users: %w", err).Error()})
	}
	return c.JSON(users)
}

func (h *Handler) GetUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).SendString("Invalid UUID format")
	}

	user, err := h.service.GetUser(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(user)
}

func (h *Handler) CreateUser(c *fiber.Ctx) error {
	var user User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.service.CreateUser(c.Context(), &user); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": fmt.Errorf("failed to create users: %w", err).Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *Handler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userId := uuid.MustParse(id)

	var user User
	user.ID = userId
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.service.UpdateUser(c.Context(), userId, user.Email); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": fmt.Errorf("failed to update users: %w", err).Error()})
	}
	return c.JSON(user)
}

func (h *Handler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userId := uuid.MustParse(id)

	if err := h.service.DeleteUser(c.Context(), userId); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": fmt.Errorf("failed to delete users: %w", err).Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
