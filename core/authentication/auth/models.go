package auth

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type StatusResponse struct {
	Status bool `json:"status"`
}

type LoginResponse struct {
	UserData *UserData `json:"userData"`
	Token    string    `json:"token"`
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

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token              string `json:"-"`
	NewPassword        string `json:"newPassword"`
	ConfirmNewPassword string `json:"confirmNewPassword"`
}

type OAuthLoginRequest struct {
	Provider string `json:"provider"`
	Token    string `json:"token"`
}
