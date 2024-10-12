package users

import "gorm.io/gorm"

type userApi struct {
	db *gorm.DB
}

type UserAPI interface{
	GetUsers(req *FindRequest) (*[]User, error)
}

func NewUserAPI(db *gorm.DB) UserAPI {
	return &userApi{db: db}
}


// @Summary      	GetUsers
// @Description	
// @Tags			Users
// @Produce			json
// @Success			200								{object}	FindResponse
// @Router			/users		[GET]
func (api *userApi) GetUsers(req *FindRequest) (*[]User, error) {
	var users []User

	if err := api.db.Where("id = ?", req.UserID).Find(&users).Error; err != nil {
		return nil, err
	}

	return &users, nil
}

