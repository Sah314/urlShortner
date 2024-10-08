package services

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"urlShortner/database"
	"urlShortner/repository"
)

func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")
	r := database.CreateClient(0)
	defer func(r *redis.Client) {
		err := r.Close()
		if err != nil {
			log.Error("Error: ", err)
		}
	}(r)

	value, err := r.Get(database.Ctx, url).Result()
	if errors.Is(err, redis.Nil) {
		v, err := repository.GetURL(url)
		if err != nil {
			log.Error("Error: ", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Cannot connect to DB",
			})
		}
		if v == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "ShortUrl not found",
			})
		}
		value = v.Longurl
		log.Error("Error: ", err)
		//TODO: Also check in the database

	} else if err != nil {
		log.Error("Error connecting to redis DB: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot connect to redis DB",
		})
	}
	rInr := database.CreateClient(1)
	defer func(rInr *redis.Client) {
		err = rInr.Close()
		if err != nil {
			log.Error("Error: ", err)
		}
	}(rInr)

	cmd := rInr.Incr(database.Ctx, "counter")
	err = cmd.Err()
	if err != nil {
		log.Error("Error: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot increment counter",
		})
	}

	return c.Redirect(value, 301)
}
