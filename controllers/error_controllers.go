package controllers

import "github.com/gofiber/fiber/v2"

func GetError403(c *fiber.Ctx) error {
	return c.Render("pages/errors/403", fiber.Map{})
}

func GetError401(c *fiber.Ctx) error {
	return c.Render("pages/errors/401", fiber.Map{})
}

func GetError404(c *fiber.Ctx) error {
	return c.Render("pages/errors/404", fiber.Map{})
}
