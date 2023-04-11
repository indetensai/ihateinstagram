package controllers

import (
	"phota/internal/entities"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type followingServiceHandler struct {
	FollowingService entities.FollowingService
}

func NewFollowingServiceHandler(app *fiber.App, f entities.FollowingService) {
	handler := &followingServiceHandler{f}
	app.Post("/user/:user_id<guid>/follow", handler.FollowHandler)
	app.Post("/user/:user_id<guid>/unfollow", handler.UnfollowHandler)
	app.Get("/user/:user_id<guid>/followers", handler.GetFollowersHandler)
}

func (f *followingServiceHandler) FollowHandler(c *fiber.Ctx) error {
	follower_id, ok := c.Locals("user_id").(*uuid.UUID)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	user_id, err := uuid.Parse(c.Params("user_id"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	err = f.FollowingService.Follow(*follower_id, user_id, c.Context())
	if err != nil {
		return error_handling(c, err)
	}
	return c.SendStatus(fiber.StatusOK)
}
func (f *followingServiceHandler) UnfollowHandler(c *fiber.Ctx) error {
	follower_id, ok := c.Locals("user_id").(*uuid.UUID)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	user_id, err := uuid.Parse(c.Params("user_id"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	err = f.FollowingService.Unfollow(*follower_id, user_id, c.Context())
	if err != nil {
		return error_handling(c, err)
	}
	return c.SendStatus(fiber.StatusOK)
}
func (f *followingServiceHandler) GetFollowersHandler(c *fiber.Ctx) error {
	user_id, err := uuid.Parse(c.Params("user_id"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	followers, err := f.FollowingService.GetFollowers(user_id, c.Context())
	if err != nil {
		return error_handling(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"followers": followers})
}
