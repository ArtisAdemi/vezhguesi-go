package users

import (
	"time"
)

type SignupRequest struct {
	Email           string `json:"email"`
	Username        string `json:"username"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type SignupResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

type SignupVerifyRequest struct {
	Token string `json:"-"`
}

type StatusResponse struct {
	Status bool `json:"status"`
}

type UserData struct {
	ID           int    `json:"id"`
	ProfileID    int    `json:"profileId"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	Role         string `json:"role"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	AvatarImgUrl string `json:"avatarImgUrl"`
}

type UserRequest struct {
	UserID    int    `json:"-"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
}

type IDResponse struct {
	ID int `json:"id"`
}

type UserResponse struct {
	ID           int            `json:"id"`
	FirstName    string         `json:"firstName"`
	LastName     string         `json:"lastName"`
	Username     string         `json:"username"`
	Email        string         `json:"email"`
	Status       string         `json:"status"`
	AvatarImgKey string         `json:"avatarImgUrl"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
}

type FindUserByID struct {
	UserID int `json:"-"`
}

type PasswordUpdateRequest struct {
	UserID             int    `json:"-"`
	// Mode               string `json:"mode"`
	CurrentPassword    string `json:"currentPassword"`
	NewPassword        string `json:"newPassword"`
	ConfirmNewPassword string `json:"confirmNewPassword"`
}

type EmailUpdateRequest struct {
	UserID   int    `json:"-"`
	NewEmail string `json:"newEmail"`
}

type FindRequest struct {
	UserID int `json:"-"`
}

type FindResponse struct {
	Users []User `json:"users"`
}

type FindByIDResponse struct {
	ID int `json:"id"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Username string `json:"username"`
	Email string `json:"email"`
	Status string `json:"status"`
	AvatarImgKey string `json:"avatarImgUrl"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
