package repository

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"time"
	"urlShortner/db"
	"urlShortner/models"
)

func StoreURL(url *models.URL) (bool, error) {
	ctx := context.Background()
	client, err := SetupUserRepository()
	if err != nil {
		return false, err
	}
	defer func() {
		if err = client.Prisma.Disconnect(); err != nil {
			log.Error(fmt.Sprintf("Failed to disconnect from database: %v", err))

		}
	}()

	_, err = client.URL.CreateOne(
		db.URL.ShortURL.Set(url.Shorturl),
		db.URL.LongURL.Set(url.Longurl),
		db.URL.User.Link(
			db.User.ID.Equals(url.UserId), // Link to the user by ID
		),
		db.URL.UpdatedAt.Set(time.Now()),
		db.URL.CreatedAt.Set(time.Now()),
	).Exec(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to save url: %v", err))
		return false, err
	}
	return true, nil
}

func GetURL(short string) (*models.URL, error) {
	ctx := context.Background()
	client, err := SetupUserRepository()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = client.Prisma.Disconnect(); err != nil {
			log.Error(fmt.Sprintf("Failed to disconnect from database: %v", err))

		}
	}()

	url, err := client.URL.FindFirst(
		db.URL.ShortURL.Equals(short),
	).Exec(ctx)

	if err != nil {
		log.Error(fmt.Sprintf("Failed to get url: %v", err))
		return nil, err
	}

	return &models.URL{
		Shorturl: url.ShortURL,
		Longurl:  url.LongURL,
		UserId:   url.UserID,
	}, nil

}

func GetUserURLs(userID string) ([]*models.URL, error) {
	ctx := context.Background()
	client, err := SetupUserRepository()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = client.Prisma.Disconnect(); err != nil {
			log.Error(fmt.Sprintf("Failed to disconnect from database: %v", err))
		}
	}()

	urls, err := client.URL.FindMany(
		db.URL.UserID.Equals(userID),
	).Exec(ctx)

	if err != nil {
		log.Error(fmt.Sprintf("Failed to get urls: %v", err))
		return nil, err
	}
	var userURLs []*models.URL
	for _, url := range urls {
		userURLs = append(userURLs, &models.URL{
			Shorturl: url.ShortURL,
			Longurl:  url.LongURL,
			UserId:   url.UserID,
		})
	}
	return userURLs, nil
}
