package orgs

import (
	"fmt"
	"regexp"
	"strings"

	subscriptionsvc "vezhguesi/app/subscriptions"
	helper "vezhguesi/helper"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type orgApi struct {
	db *gorm.DB
	logger log.AllLogger
}

type OrgAPI interface{
	Add(req *AddOrgRequest) (res *OrgResponse, err error)
}

func NewOrgAPI(db *gorm.DB, logger log.AllLogger) OrgAPI {
	return &orgApi{
		db: db,
		logger: logger,
	}
}

// @Summary      	Add
// @Description		Validates user id, org name and org size, checks if org exists in DB by name or slug, if not a new organization with trial subscription will be created and then the created ID will be returned.
// @Tags			Orgs
// @Accept			json
// @Produce			json
// @Param			Authorization					header		string			true	"Authorization Key(e.g Bearer key)"
// @Param			AddOrgRequest					body		AddOrgRequest	true	"AddOrgRequest"
// @Success			200								{object}	OrgResponse
// @Router			/api/orgs	[POST]
func (s *orgApi) Add(req *AddOrgRequest) (res *OrgResponse, err error) {
	if req.UserID == 0 {
		return nil, fmt.Errorf("user id is required")
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Size = strings.TrimSpace(req.Size)
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.Size == "" {
		return nil, fmt.Errorf("size is required")
	}

	var user User
	s.db.Where("id = ?", req.UserID).First(&user)
	if user.ID == 0 {
		return nil, helper.ErrNotFound
	}
	var org Org 
	orgSlug := regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(req.Name, "")
	orgSlug = strings.ReplaceAll(strings.TrimSpace(orgSlug), " ", "-")
	s.db.Where("slug = ?", orgSlug).First(&org)
	if org.ID != 0 {
		return nil, fmt.Errorf("org slug already exists")
	}

	var trialSubscription subscriptionsvc.Subscription
	trialSubscription.Name = "Trial"
	trialSubDesc := "Platform trial"
	trialSubscription.Description = trialSubDesc
	trialSubscription.DurationType = "day"
	trialSubscription.DurationTime = 14
	trialSubscription.Features = append(trialSubscription.Features, 
		subscriptionsvc.Feature{
			Key: "OrgCreateLimit",
			Value: "10",
		},
		subscriptionsvc.Feature{
			Key: "AdminCreateLimit",
			Value: "2",
		},
		subscriptionsvc.Feature{
			Key: "UserCreateLimit",
			Value: "20",
		})

	result := s.db.Omit("UpdatedAt").Create(&trialSubscription)
	if result.Error != nil {
		return nil, result.Error
	}

	org.Name = req.Name
	org.Size = req.Size
	org.Slug = orgSlug
	org.SubscriptionID = trialSubscription.ID

	result = s.db.Omit("UpdatedAt").Create(&org)
	if result.Error != nil {
		return nil, result.Error
	}
	// TODO: Create org settings
	// var orgSettings OrgSettings
	// orgSettings.OrgID = org.ID
	// orgSettings.


	// Find owner role
	var ownerRole Role
	result = s.db.Where("name = ?", helper.OwnerRoleName).First(&ownerRole)
	if result.Error != nil {
		return nil, result.Error
	}

	// Save relationship
	var usrOrgRole UserOrgRole
	usrOrgRole.OrgID = org.ID
	usrOrgRole.UserID = int(user.ID)
	usrOrgRole.RoleID = int(ownerRole.ID)
	usrOrgRole.Status = "active"

	result = s.db.Omit("UpdatedAt").Create(&usrOrgRole)
	if result.Error != nil {
		return nil, result.Error
	}

	result = s.db.Save(&org)
	if result.Error != nil {
		return nil, result.Error
	}

	return &OrgResponse{
		ID: org.ID,
		Name: org.Name,
		OrgSlug: org.Slug,
	}, nil
}
