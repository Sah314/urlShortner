package repository

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"urlShortner/db"
	"urlShortner/models"
)

func SetupUserRepository() (*db.PrismaClient, error) {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		log.Error(fmt.Sprintf("Failed to connect to database: %v", err))
		return nil, err
	}
	return client, nil
}

func GetUser(email string) *models.User {
	ctx := context.Background()
	client, err := SetupUserRepository()
	defer func() {
		if err = client.Prisma.Disconnect(); err != nil {
			log.Error(fmt.Sprintf("Failed to disconnect from database: %v", err))

		}
	}()
	if err != nil {
		log.Error(fmt.Sprintf("Failed to setup user repository: %v", err))
		return nil
	}
	user, err := client.User.FindUnique(
		db.User.Email.Equals(email)).Exec(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to get user: %v", err))
		return nil
	}
	return &models.User{Id: user.ID, Email: user.Email, FirstName: user.FirstName, LastName: user.LastName, Password: user.Password}

}

//func getUsers(limit int, offset int) []*User {
//
//}

func CreateUser(user *models.User) (*models.User, error) {
	client, err := SetupUserRepository()
	if err != nil {
		log.Error(fmt.Sprintf("Failed to setup user repository: %v", err))
		return nil, err
	}
	ctx := context.Background()
	usr, err := client.User.CreateOne(
		db.User.FirstName.Set(user.FirstName),
		db.User.LastName.Set(user.LastName), db.User.Email.Set(user.Email), db.User.Password.Set(user.Password)).Exec(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to create user: %v", err))
		return nil, err
	}
	return &models.User{
		Id:        usr.ID,
		Email:     usr.Email,
		FirstName: usr.FirstName,
		LastName:  usr.LastName,
		Password:  usr.Password,
	}, nil
}

func UpdateUser(user *models.User) (*models.User, error) {
	client, err := SetupUserRepository()
	if err != nil {
		log.Error(fmt.Sprintf("Failed to setup user repository: %v", err))
		return nil, err
	}
	ctx := context.Background()
	usr, err := client.User.UpsertOne(
		db.User.Email.Equals(user.Email),
	).Update(
		db.User.FirstName.Set(user.FirstName),
		db.User.LastName.Set(user.LastName),
		db.User.Password.Set(user.Password),
	).Exec(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to update user: %v", err))
		return nil, err
	}
	return &models.User{
		Id:        usr.ID,
		Email:     usr.Email,
		FirstName: usr.FirstName,
		LastName:  usr.LastName,
		Password:  usr.Password,
	}, nil
}

//func CreateGoogleUser() (*models.User, error) {
//
//}
