package seeds

import (
	"fmt"
	rolesvc "vezhguesi/core/authorization/role"

	"gorm.io/gorm"
)

var defaultRoleMappingPerms map[string][]rolesvc.Permission = map[string][]rolesvc.Permission{
	Owner: {
		orgSettingsPerms[read],
		orgSettingsPerms[update],
		orgSettingsPerms[delete],
		orgSettingsPerms[uploadLogoHeroImgs],
		orgSettingsPerms[orgImgs],

		analyticsPerms[create],
		analyticsPerms[read],
		analyticsPerms[update],
		analyticsPerms[delete],

		userPerms[create],
		userPerms[read],
		userPerms[update],
		userPerms[delete],
		userPerms[invite],
		userPerms[updateProfile],
		userPerms[uploadAvatar],
		userPerms[downloadAvatar],
		userPerms[roleCount],

		contentPerms[create],
		contentPerms[read],
		contentPerms[update],
		contentPerms[delete],

		contentLibPerms[create],
		contentLibPerms[read],
		contentLibPerms[update],
		contentLibPerms[delete],

		projectPerms[read],
		bountyPerms[read],
		productServicePerms[read],

		messagingPerms[create],
		messagingPerms[read],
		messagingPerms[update],
		messagingPerms[delete],

		meetingPerms[create],
		meetingPerms[read],
		meetingPerms[update],
		meetingPerms[delete],

		activityPerms[activityDashboardStats],
	},
	Admin: {
		analyticsPerms[create],
		analyticsPerms[read],
		analyticsPerms[update],
		analyticsPerms[delete],

		userPerms[create],
		userPerms[read],
		userPerms[update],
		userPerms[delete],
		userPerms[invite],
		userPerms[updateProfile],
		userPerms[uploadAvatar],
		userPerms[downloadAvatar],
		userPerms[roleCount],

		contentPerms[create],
		contentPerms[read],
		contentPerms[update],
		contentPerms[delete],

		contentLibPerms[create],
		contentLibPerms[read],
		contentLibPerms[update],
		contentLibPerms[delete],

		projectPerms[read],
		bountyPerms[read],
		productServicePerms[read],

		messagingPerms[create],
		messagingPerms[read],
		messagingPerms[update],
		messagingPerms[delete],

		meetingPerms[create],
		meetingPerms[read],
		meetingPerms[update],
		meetingPerms[delete],

		activityPerms[activityDashboardStats],
		orgSettingsPerms[read],
		orgSettingsPerms[update],
		orgSettingsPerms[uploadLogoHeroImgs],
		orgSettingsPerms[orgImgs],
	},
	Coach: {
		userPerms[invite],
		userPerms[updateProfile],
		userPerms[uploadAvatar],
		userPerms[downloadAvatar],
		userPerms[roleCount],
		userPerms[read],

		contentPerms[create],
		contentPerms[read],
		contentPerms[update],
		contentPerms[delete],

		contentLibPerms[create],
		contentLibPerms[read],
		contentLibPerms[update],
		contentLibPerms[delete],

		projectPerms[read],
		bountyPerms[read],
		productServicePerms[read],

		messagingPerms[create],
		messagingPerms[read],
		messagingPerms[update],
		messagingPerms[delete],

		meetingPerms[create],
		meetingPerms[read],
		meetingPerms[update],
		meetingPerms[delete],

		activityPerms[activityDashboardStats],
		orgSettingsPerms[read],
		orgSettingsPerms[orgImgs],
	},
	SME: {
		userPerms[invite],
		userPerms[updateProfile],
		userPerms[uploadAvatar],
		userPerms[downloadAvatar],
		userPerms[roleCount],
		userPerms[read],

		contentPerms[read],
		projectPerms[read],
		bountyPerms[read],

		productServicePerms[create],
		productServicePerms[read],
		productServicePerms[update],
		productServicePerms[delete],

		messagingPerms[create],
		messagingPerms[read],
		messagingPerms[update],
		messagingPerms[delete],

		meetingPerms[create],
		meetingPerms[read],
		meetingPerms[update],
		meetingPerms[delete],

		activityPerms[activityDashboardStats],
		orgSettingsPerms[read],
		orgSettingsPerms[orgImgs],
	},
	ClientAlumn: {
		userPerms[invite],
		userPerms[updateProfile],
		userPerms[uploadAvatar],
		userPerms[downloadAvatar],
		userPerms[roleCount],
		userPerms[read],

		contentPerms[read],

		projectPerms[create],
		projectPerms[read],
		projectPerms[update],
		projectPerms[delete],

		bountyPerms[read],
		productServicePerms[read],

		messagingPerms[create],
		messagingPerms[read],
		messagingPerms[update],
		messagingPerms[delete],

		meetingPerms[create],
		meetingPerms[read],
		meetingPerms[update],
		meetingPerms[delete],

		activityPerms[activityDashboardStats],
		orgSettingsPerms[read],
		orgSettingsPerms[orgImgs],
	},
	ClientCurrent: {
		userPerms[invite],
		userPerms[updateProfile],
		userPerms[uploadAvatar],
		userPerms[downloadAvatar],
		userPerms[roleCount],
		userPerms[read],

		contentPerms[read],

		projectPerms[create],
		projectPerms[read],
		projectPerms[update],
		projectPerms[delete],

		bountyPerms[read],
		productServicePerms[read],

		messagingPerms[create],
		messagingPerms[read],
		messagingPerms[update],
		messagingPerms[delete],

		meetingPerms[create],
		meetingPerms[read],
		meetingPerms[update],
		meetingPerms[delete],

		activityPerms[activityDashboardStats],
		orgSettingsPerms[read],
		orgSettingsPerms[orgImgs],
	},
	ClientFuture: {
		userPerms[invite],
		userPerms[updateProfile],
		userPerms[uploadAvatar],
		userPerms[downloadAvatar],
		userPerms[roleCount],
		userPerms[read],

		contentPerms[read],
		projectPerms[read],
		bountyPerms[read],
		productServicePerms[read],

		messagingPerms[create],
		messagingPerms[read],
		messagingPerms[update],
		messagingPerms[delete],

		meetingPerms[create],
		meetingPerms[read],
		meetingPerms[update],
		meetingPerms[delete],

		activityPerms[activityDashboardStats],
		orgSettingsPerms[read],
		orgSettingsPerms[orgImgs],
	},
	Partner: {
		userPerms[invite],
		userPerms[updateProfile],
		userPerms[uploadAvatar],
		userPerms[downloadAvatar],
		userPerms[roleCount],
		userPerms[read],

		projectPerms[read],

		bountyPerms[create],
		bountyPerms[read],
		bountyPerms[update],
		bountyPerms[delete],

		productServicePerms[create],
		productServicePerms[read],
		productServicePerms[update],
		productServicePerms[delete],

		messagingPerms[create],
		messagingPerms[read],
		messagingPerms[update],
		messagingPerms[delete],

		meetingPerms[create],
		meetingPerms[read],
		meetingPerms[update],
		meetingPerms[delete],

		activityPerms[activityDashboardStats],
		orgSettingsPerms[read],
		orgSettingsPerms[orgImgs],
	},
	Guest: {
		userPerms[updateProfile],
		userPerms[uploadAvatar],
		userPerms[downloadAvatar],

		contentPerms[read],
		projectPerms[read],
		bountyPerms[read],
		productServicePerms[read],

		messagingPerms[create],
		messagingPerms[read],
		messagingPerms[update],
		messagingPerms[delete],

		meetingPerms[create],
		meetingPerms[read],
		meetingPerms[update],
		meetingPerms[delete],
		orgSettingsPerms[read],
		orgSettingsPerms[orgImgs],
	},
}

// db *gorm.DB
func SeedDefaultRolesAndPermissions(db *gorm.DB) {
	var roles []string = []string{Owner, Admin, Coach, SME, ClientAlumn, ClientCurrent, ClientFuture, Partner, Guest}
	for _, role := range roles {
		// Handle role
		var rl rolesvc.Role
		result := db.Where("name = ?", role).First(&rl)
		if result.Error != nil && result.Error.Error() == "record not found" {
			// Create new role
			globalOrgID := 0
			rl.OrgID = &globalOrgID // global
			rl.Name = role
			result = db.Omit("UpdatedAt").Create(&rl)
			if result.Error != nil {
				fmt.Println("Create(&rl).err", result.Error.Error())
			}
		}

		// Handle permissions
		rolePermSet := defaultRoleMappingPerms[role]
		for _, perm := range rolePermSet {
			var prm rolesvc.Permission
			result := db.Where("name = ?", perm.Name).First(&prm)
			if result.Error != nil && result.Error.Error() == "record not found" {
				// Create new permission
				prm.Name = perm.Name
				prm.HTTPMethods = perm.HTTPMethods // Correct field name
				prm.Path = perm.Path               // Correct field name

				// Log the JSON data before insertion
				fmt.Printf("Inserting Permission: Name=%s, HTTPMethods=%v, Path=%v\n", prm.Name, prm.HTTPMethods, prm.Path)

				result = db.Omit("UpdatedAt").Create(&prm)
				if result.Error != nil {
					fmt.Println("Create(&prm).err", result.Error.Error())
				}
			}

			// Save relationship
			rl.Permissions = append(rl.Permissions, prm)
			result = db.Save(&rl)
			if result.Error != nil {
				fmt.Println("db.Save(&rl)", result.Error.Error())
			}
		}
	}
}
