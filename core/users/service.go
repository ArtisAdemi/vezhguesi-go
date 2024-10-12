package users

import (
	"fmt"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2/log"
	"github.com/mattevans/postmark-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	helper "vezhguesi/helper"
)

type userApi struct {
	db *gorm.DB
	secretKey string
	postmarkClient *postmark.Client
	uiAppUrl string
	logger log.AllLogger
}

type UserAPI interface{
	GetUsers(req *FindRequest) (*[]User, error)
	Signup(req *SignupRequest) (*SignupResponse, error)
}

func NewUserAPI(db *gorm.DB, secretKey string, pmc *postmark.Client, uiAppUrl string, logger log.AllLogger) UserAPI {
	return &userApi{db: db, secretKey: secretKey, postmarkClient: pmc, uiAppUrl: uiAppUrl, logger: logger}
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


	verifyLink := s.uiAppUrl + "/verify-signup/" +t

	emailMsg := &postmark.Email{
        From: "info@vezhguesi.com",
        To: req.Email,
        Subject: "Verify your email",
		HTMLBody: fmt.Sprintf("Click on the link to verify your email: <a href=\"%s\">Click here</a>", verifyLink),
	}

	_, _, err = s.postmarkClient.Email.Send(emailMsg)
	if err != nil {
		s.logger.Errorf("func: Signup, operation: s.postmarkClient.Email.Send(emailMsg), err: %s", err.Error())
		fmt.Errorf("failed to send email")
	}

	return &SignupResponse{
       ID: user.ID,
	   Status: "pending",
	}, nil


}

