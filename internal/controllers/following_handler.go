package controllers

import (
	"phota/internal/entities"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FollowingServiceHandler struct {
	FollowingService entities.FollowingService
}

func NewFollowingServiceHandler(c *fiber.App, f entities.FollowingService) {
	handler := &FollowingServiceHandler{f}
	c.Put("/user/:user_id<guid>/followers/:follower_id<guid>", handler.FollowHandler)
}

func (f *FollowingServiceHandler) FollowHandler(c *fiber.Ctx) error {
	session_id := c.FormValue("session_id")
	user_id, err := uuid.Parse(c.Params("user_id"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	follower_id, err := uuid.Parse(c.Params("follower_id"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	err = f.FollowingService.Follow(follower_id, user_id, session_id, c.Context())
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}
func (f *FollowingServiceHandler) UnfollowHandler(c *fiber.Ctx) error {
	session_id := c.FormValue("session_id")
	user_id, err := uuid.Parse(c.Params("user_id"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	follower_id, err := uuid.Parse(c.Params("follower_id"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	err = f.FollowingService.Unfollow(follower_id, user_id, session_id, c.Context())
	if err == entities.ErrNotAuthorized {
		c.SendStatus(fiber.StatusUnauthorized)
	}
	if err != nil {
		c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}
func (f *FollowingServiceHandler) GetFollowersHandler(c *fiber.Ctx) error {
	user_id, err := uuid.Parse(c.Params("user_id"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	followers, err := f.FollowingService.GetFollowers(user_id, c.Context())
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"followers": *followers})
}
