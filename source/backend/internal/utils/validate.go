package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// validate is the shared instance of the validator used for request body validation.
var validate = validator.New()

// ValidateBody returns a Fiber middleware that parses and validates the request body against a given struct type T.
//
// If the JSON is malformed or validation fails, it returns a 400 Bad Request with an error message.
// On success, the validated struct is stored in c.Locals("body") for downstream handlers to access.
//
// Usage:
//
//	app.Post("/example", utils.ValidateBody[MyStruct](), func(c *fiber.Ctx) error {
//		body := c.Locals("body").(MyStruct)
//		// Handle request...
//	})
func ValidateBody[T any]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body T
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid JSON",
			})
		}

		if err := validate.Struct(body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		c.Locals("body", body)
		return c.Next()
	}
}
