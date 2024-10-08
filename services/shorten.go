package services

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"os"
	"strconv"
	"time"
	"urlShortner/database"
	"urlShortner/helpers"
	"urlShortner/models"
	"urlShortner/repository"
)

func ShortenURL(c *fiber.Ctx) error {
	body := new(models.ShortenRequest)
	err := helpers.ParseRequestBody(c, body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	userId, ok := c.Locals("userId").(string)
	if userId == "" || !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	log.Debug("userId", userId)
	r2 := database.CreateClient(1)
	defer func(r2 *redis.Client) {
		err = r2.Close()
		if err != nil {
			log.Error("Error", err)
		}
	}(r2)

	val, err := r2.Get(database.Ctx, c.IP()).Result()
	if errors.Is(err, redis.Nil) {
		_ = r2.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
		log.Error("Error", err)
	} else {
		val, _ = r2.Get(database.Ctx, c.IP()).Result()
		log.Debug("val", val)
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			limit, _ := r2.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":           "Rate Limit Exceeded",
				"rate_limit_rest": limit / time.Nanosecond / time.Minute,
			})
		}
	}
	//Validate input URL
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid URL",
		})
	}

	// Validate Domain Error
	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Cannot access domain",
		})
	}
	body.URL = helpers.EnforceHTTP(body.URL)

	var id string

	if body.CustomURL == "" {
		id = uuid.New().String()[:7]
	} else {
		id = body.CustomURL
	}
	r := database.CreateClient(0)
	defer func(r *redis.Client) {
		err = r.Close()
		if err != nil {
			log.Error("Error", err)
		}
	}(r)

	val, _ = r.Get(database.Ctx, id).Result()
	if val != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "This shortURL already exists",
		})
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}
	redisErrChan := make(chan error, 1)
	dbErrChan := make(chan error, 1)

	go func() {
		err = r.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
		redisErrChan <- err // send error or nil to the channel
	}()

	go func() {
		savedurl, err := repository.StoreURL(&models.URL{
			Shorturl: id,
			Longurl:  body.URL,
			Expiry:   body.Expiry,
			UserId:   userId,
		})
		if !savedurl {
			dbErrChan <- errors.New("unable to save to database")
		} else {
			dbErrChan <- err
		}
	}()

	redisErr := <-redisErrChan
	dbErr := <-dbErrChan

	resp := models.Response{
		URL:             body.URL,
		CustomURL:       "",
		Expiry:          body.Expiry,
		XRateRemaining:  10,
		XRateLimitReset: 30,
	}
	if redisErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "unable to connect to Redis",
		})
	}
	if dbErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "unable to save to database",
		})
	}

	r2.Decr(database.Ctx, c.IP())
	val, _ = r2.Get(database.Ctx, c.IP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)
	ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()
	resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute
	resp.CustomURL = os.Getenv("DOMAIN") + "/" + id
	return c.Status(fiber.StatusOK).JSON(resp)
}
