package models

import "time"

type ShortenRequest struct {
	URL       string        `json:"url"`
	CustomURL string        `json:"customURL"`
	Expiry    time.Duration `json:"expiry"`
}

type Response struct {
	URL             string        `json:"url"`
	CustomURL       string        `json:"customURL"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"xRateRemaining"`
	XRateLimitReset time.Duration `json:"xRateLimitRest"`
}

type User struct {
	Id        string
	Email     string
	FirstName string
	LastName  string
	Password  string
}

type URL struct {
	Id       string
	Longurl  string
	Shorturl string
	UserId   string
	Expiry   time.Duration
}

type SignupRequest struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
