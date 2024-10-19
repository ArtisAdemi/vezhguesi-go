package auth

import (
	"fmt"
	"strings"
	"time"
	"vezhguesi/core/users"
	"vezhguesi/helper"

	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

type authApi struct {
	db *gorm.DB
	secretKey string
	mailDialer *gomail.Dialer // Use gomail Dialer
	uiAppUrl string
	logger log.AllLogger
}

type AuthApi interface{
	Signup(req *SignupRequest) (*SignupResponse, error)
	VerifySignup(req *SignupVerifyRequest) (*StatusResponse, error)
	Login(req *LoginRequest) (*LoginResponse, error)
	UpdateUser(req *UpdateUserRequest) (*UserData, error)
}

func NewAuthApi(db *gorm.DB, secretKey string, dialer *gomail.Dialer, uiAppUrl string, logger log.AllLogger) AuthApi {
	return &authApi{db: db, secretKey: secretKey, mailDialer: dialer, uiAppUrl: uiAppUrl, logger: logger}
}

// @Summary      	Signup
// @Description	Validates email, username, first name, last name, password checks if email exists, if not creates new user and sends email with verification link.
// @Tags			Auth
// @Accept			json
// @Produce			json
// @Param			SignupRequest	body		SignupRequest	true	"SignupRequest"
// @Success			200					{object}	SignupResponse
// @Router			/api/auth/	[POST]
func (s *authApi) Signup(req *SignupRequest) (*SignupResponse, error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Username = strings.TrimSpace(req.Username)
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Password = strings.TrimSpace(req.Password)
	req.ConfirmPassword = strings.TrimSpace(req.ConfirmPassword)

	if req.Email == "" || req.Password == "" || req.ConfirmPassword == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	if !helper.ValidEmail(req.Email) {
		return nil, fmt.Errorf("invalid email")
	}

	if req.Password != req.ConfirmPassword {
		return nil, fmt.Errorf("password and confirm password do not match")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 12).Unix()
	claims["email"] = req.Email

	t, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return nil, fmt.Errorf("failed to generate token")
	}

	var user users.User 
	_ = s.db.Where("email = ?", req.Email).First(&user)
	if user.ID > 0 {
		if !user.VerifiedEmail {
			return nil, fmt.Errorf("verify your email first")
		}

		return nil, fmt.Errorf("email already in use")
	}

	_ = s.db.Where("username = ?", req.Username).First(&user)
	if user.ID > 0 {
		return nil, fmt.Errorf("username already in use")
	}

	user.Email = req.Email
	user.Username = &req.Username
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.VerifiedEmail = false
	user.Active = false

	hashedPw, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Errorf("func: Signup, operation: bcrypt")
		return nil, fmt.Errorf("failed to hash password")
	}

	pwhs := string(hashedPw)
	user.Password = pwhs

	result := s.db.Omit("UpdatedAt").Create(&user)
	if result.Error != nil {
		s.logger.Errorf("func: Signup, operation: s.db.Omit('UpdatedAt').Create(&user), err: %s", err.Error())
		return nil, result.Error
	}

	s.db.Model(users.User{Email: req.Email}).First(&user)


	verifyLink := s.uiAppUrl + "/verify-signup/" + t

	m := gomail.NewMessage()
	m.SetHeader("From", "info@vezhguesi.com")
	m.SetHeader("To", req.Email)
	m.SetHeader("Subject", "Verify your email")
	m.SetBody("text/html", fmt.Sprintf("Click on the link to verify your email: <a href=\"%s\">Click here</a>", verifyLink))

	if err := s.mailDialer.DialAndSend(m); err != nil {
		s.logger.Errorf("func: Signup, operation: s.mailDialer.DialAndSend(m), err: %s", err.Error())
		return nil, fmt.Errorf("failed to send email")
	}

	return &SignupResponse{
		ID: user.ID,
		Status: "pending",
	}, nil
}

// @Summary      	VerifySignup
// @Description	Validates token in param, if token parses valid then user will be verified and be updated in DB.
// @Tags			Auth
// @Accept			json
// @Produce			json
// @Param			token				path		string			true	"Token"
// @Success			200					{object}	StatusResponse
// @Router			/api/auth/verify-signup/{token}	[GET]
func (s *authApi) VerifySignup(req *SignupVerifyRequest) (res *StatusResponse, err error) {
	req.Token = strings.TrimSpace(req.Token)
	
	if req.Token == "" {
		return nil, fmt.Errorf("token is required")
	}

	//  Parse and validate token expiration
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(req.Token, claims, func(token *jwt.Token) (interface{}, error){
		return []byte(s.secretKey), nil
	})
	if err != nil {
		s.logger.Errorf("func: VerifySignup, operation: jwt.ParseWithClaims(req.Token, claims, func(token *jwt.Token) (interface{}, error){return []byte(s.secretKey), nil}), err: %s", err.Error())
		return nil, fmt.Errorf("invalid token")
	}

	email := fmt.Sprintf("%v", claims["email"])

	// find user by email
	var user users.User
	s.db.Where("email = ?", email).First(&user)
	if user.ID == 0 {
		return nil, helper.ErrNotFound
	}

	user.VerifiedEmail = true
	user.Active = true
	result := s.db.Save(&user)
	if result.Error != nil {
		s.logger.Errorf("func: VerifySignup, operation: s.db.Save(&user), err: %s", result.Error.Error())
		return nil, result.Error
	}

	return &StatusResponse{
		Status: true,
	}, nil
}

// @Summary      	Login
// @Description		Validates email and password in request, check if user exists in DB if not throw 404 otherwise compare the request password with hash, then check if user is active, then finds relationships of user with orgs and then generates a JWT token, and returns UserData, Orgs, and Token in response.
// @Tags			Auth
// @Accept			json
// @Produce			json
// @Param			LoginRequest	body		LoginRequest	true	"LoginRequest"
// @Success			200				{object}	LoginResponse
// @Router			/api/auth/login			[POST]
func (s *authApi) Login(req *LoginRequest) (res *LoginResponse, err error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Password = strings.TrimSpace(req.Password)

	if req.Email == "" || req.Password == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	if !helper.ValidEmail(req.Email) {
		return nil, fmt.Errorf("invalid email")
	}

	var user users.User
	s.db.Where("email = ?", req.Email).First(&user)
	if user.ID == 0 {
		fmt.Println("err", "email not found")
		return nil, helper.ErrNotFound
	}
	fmt.Println("user", user)

	if !user.VerifiedEmail {
		return nil, fmt.Errorf("email not verified")
	}

	if !user.Active {
		return nil, fmt.Errorf("user is not active")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, helper.ErrNotFound
	}

	// Generate JWT Token with expiration
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	t, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return nil, fmt.Errorf("failed to generate token")
	}
	
	userData := UserData{
		ID: user.ID,
		Email: user.Email,
		Username: *user.Username,
		FirstName: user.FirstName,
		LastName: user.LastName,
	}
	
	return &LoginResponse{
		UserData: &userData,
		Token: t,
	}, nil
}

// @Summary      	UpdateUser
// @Description		Updates user data in DB.
// @Tags			Auth
// @Accept			json
// @Produce			json
// @Param			Authorization  header string true "Authorization Key (e.g Bearer key)"
// @Param			UpdateUserRequest	body		UpdateUserRequest	true	"UpdateUserRequest"
// @Success			200				{object}	UserData
// @Router			/api/auth/update			[PUT]
func (s *authApi) UpdateUser(req *UpdateUserRequest) (res *UserData, err error) {
	if req.UserID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}
	if req.FirstName == "" || req.LastName == "" || req.Username == "" {
		return nil, fmt.Errorf("first name, last name and username are required")
	}

	var user users.User
	result := s.db.First(&user, req.UserID)
	if result.Error != nil {
		return nil, result.Error
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Username = &req.Username

	result = s.db.Save(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &UserData{
		ID: user.ID,
		FirstName: user.FirstName,
		LastName: user.LastName,
		Username: *user.Username,
		Email: user.Email,
		Role: user.Role,
		AvatarImgUrl: user.AvatarImgKey,
	}, nil
}