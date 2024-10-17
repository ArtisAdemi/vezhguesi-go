package users

import (
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
	GetUsers(req *FindRequest) (*[]User, error)
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

	// Fetch all users from the database
	if err := s.db.Find(&users).Error; err != nil {
		s.logger.Errorf("Error fetching users: %v", err)
		return nil, err
	}

	s.logger.Infof("Number of users found: %d", len(users))

	return &users, nil
}
