package users

import (
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	"gopkg.in/gomail.v2" // Import gomail
	"gorm.io/gorm"
)

type userApi struct {
	db *gorm.DB
	secretKey string
	mailDialer *gomail.Dialer // Use gomail Dialer
	uiAppUrl string
	logger log.AllLogger
}

type UserAPI interface{
	GetUsers(req *FindRequest) (*[]UserResponse, error)
	GetUserByID(req *FindUserByID) (*FindByIDResponse, error)
	GetUserData(req *FindUserByID) (*UserData, error)
}

func NewUserAPI(db *gorm.DB, secretKey string, dialer *gomail.Dialer, uiAppUrl string, logger log.AllLogger) UserAPI {
	return &userApi{db: db, secretKey: secretKey, mailDialer: dialer, uiAppUrl: uiAppUrl, logger: logger}
}


// @Summary      	GetUsers
// @Description
// @Tags			Users
// @Produce			json
// @Success			200								{object}	[]UserResponse
// @Router			/api/users		[GET]
func (s *userApi) GetUsers(req *FindRequest) (*[]UserResponse, error) {
	var users []User

	// Fetch all users from the database
	if err := s.db.Find(&users).Error; err != nil {
		s.logger.Errorf("Error fetching users: %v", err)
		return nil, err
	}

	s.logger.Infof("Number of users found: %d", len(users))

	var userResponses []UserResponse
	for _, user := range users {
		status := "inactive"
		if user.Active {
			status = "active"
		} 

		userResponses = append(userResponses, UserResponse{
			ID: user.ID,
			FirstName: user.FirstName,
			LastName: user.LastName,
			Username: *user.Username,
			Email: user.Email,
			Status: status,
		})
	}

	return &userResponses, nil
}


// @Summary      	GetUserByID
// @Description
// @Tags			Users
// @Produce			json
// @Param			Authorization  header string true "Authorization Key (e.g Bearer key)"
// @Param			userId  path int true "User ID"
// @Success			200								{object}	FindByIDResponse
// @Router			/api/users/{userId}		[GET]
func (s *userApi) GetUserByID(req *FindUserByID) (res *FindByIDResponse, err error) {
	var user User
	if err := s.db.First(&user, req.UserID).Error; err != nil {
		s.logger.Errorf("Error fetching user by ID: %v", err)
		return nil, err
	}

	var response = &FindByIDResponse{
		ID: user.ID,
		FirstName: user.FirstName,
		LastName: user.LastName,
		Username: *user.Username,
		Email: user.Email,
		Status: user.Status,
		AvatarImgKey: user.AvatarImgKey,
	}

	return response, nil
}

// @Summary      	GetUserData
// @Description
// @Tags			Users
// @Produce			json
// @Param			Authorization  header string true "Authorization Key (e.g Bearer key)"
// @Success			200								{object}	UserData
// @Router			/api/users/user-data		[GET]
func (s *userApi) GetUserData(req *FindUserByID) (res *UserData, err error) {
	if req.UserID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}

	var user User
	result := s.db.First(&user, req.UserID)
	if result.Error != nil {
		return nil, result.Error
	}

	var userData = &UserData{
		ID: user.ID,
		FirstName: user.FirstName,
		LastName: user.LastName,
		Username: *user.Username,
		Email: user.Email,
		Role: user.Role,
	}

	return userData, nil
}