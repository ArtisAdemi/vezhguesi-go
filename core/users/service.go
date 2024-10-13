package users

import (
	"fmt"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2" // Import gomail
	"gorm.io/gorm"

	helper "vezhguesi/helper"
)

type userApi struct {
	db *gorm.DB
	secretKey string
	mailDialer *gomail.Dialer // Use gomail Dialer
	uiAppUrl string
	logger log.AllLogger
}

type UserAPI interface{
	GetUsers(req *FindRequest) (*[]User, error)
	Signup(req *SignupRequest) (*SignupResponse, error)
	VerifySignup(req *SignupVerifyRequest) (*StatusResponse, error)
}

func NewUserAPI(db *gorm.DB, secretKey string, dialer *gomail.Dialer, uiAppUrl string, logger log.AllLogger) UserAPI {
	return &userApi{db: db, secretKey: secretKey, mailDialer: dialer, uiAppUrl: uiAppUrl, logger: logger}
}

// @Summary      	GetUsers
// @Description	
// @Tags			Users
// @Produce			json
// @Success			200								{object}	FindResponse
// @Router			/api/users		[GET]
func (s *userApi) GetUsers(req *FindRequest) (*[]User, error) {
	var users []User

	if err := s.db.Where("id = ?", req.UserID).Find(&users).Error; err != nil {
		return nil, err
	}

	return &users, nil
}

// @Summary      	Signup
// @Description	Validates email, username, first name, last name, password checks if email exists, if not creates new user and sends email with verification link.
// @Tags			Users
// @Accept			json
// @Produce			json
// @Param			SignupRequest	body		SignupRequest	true	"SignupRequest"
// @Success			200					{object}	SignupResponse
// @Router			/api/users/	[POST]
func (s *userApi) Signup(req *SignupRequest) (*SignupResponse, error) {
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

	var user User 
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

	s.db.Model(User{Email: req.Email}).First(&user)


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
// @Tags			Users
// @Accept			json
// @Produce			json
// @Param			token				path		string			true	"Token"
// @Success			200					{object}	StatusResponse
// @Router			/api/users/verify-signup/{token}	[GET]
func (s *userApi) VerifySignup(req *SignupVerifyRequest) (res *StatusResponse, err error) {
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
	var user User
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
