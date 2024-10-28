package seeds

import (
	"vezhguesi/core/authorization/role"
	rolesvc "vezhguesi/core/authorization/role"
)

const (
	// Default roles
	Owner         = string("owner")
	Admin         = string("admin")
	Coach         = string("coach")
	SME           = string("sme")
	ClientAlumn   = string("client-alumn")
	ClientCurrent = string("client-current")
	ClientFuture  = string("client-future")
	Partner       = string("partner")
	Guest         = string("guest")
	// Default permission names
	create                 = string("create")
	read                   = string("read")
	update                 = string("update")
	delete                 = string("delete")
	invite                 = string("invite")
	updateProfile          = string("update-profile")
	uploadAvatar           = string("upload-avatar")
	downloadAvatar         = string("download-avatar")
	activityDashboardStats = string("activity-dashboard-stats")
	uploadLogoHeroImgs     = string("upload-logo-hero-imgs")
	roleCount              = string("role-count")
	orgImgs                = string("org-imgs")
)

// Map default permissions with their values
var orgSettingsPerms map[string]role.Permission = map[string]role.Permission{
	read: {
		Name:       "org-settings:read",
		HTTPMethods: "GET",
		Path:       "/api/o/:orgId/org-settings",
	},
	update: {
		Name:       "org-settings:update",
		HTTPMethods: "PUT",
		Path:       "/api/o/:orgId/org-settings",
	},
	delete: {
		Name:       "org-settings:delete",
		HTTPMethods: "DELETE",
		Path:       "/api/o/:orgId/org-settings",
	},
	uploadLogoHeroImgs: {
		Name:       "org-settings:upload",
		HTTPMethods: "PUT",
		Path:       "/api/o/:orgId/org-settings/upload",
	},
	orgImgs: {
		Name:       "org-settings:imgs",
		HTTPMethods: "GET",
		Path:       "/api/o/:orgId/org-settings/imgs",
	},
}

var analyticsPerms map[string]role.Permission = map[string]role.Permission{
	create: {
		Name:       "analytics:create",
		HTTPMethods: "POST",
		Path:       "/api/o/:orgId/analytics",
	},
	read: {
		Name:       "analytics:read",
		HTTPMethods: "READ",
		Path:       "/api/o/:orgId/analytics",
	},
	update: {
		Name:       "analytics:update",
		HTTPMethods: "UPDATE",
		Path:       "/api/o/:orgId/analytics/:id",
	},
	delete: {
		Name:       "analytics:delete",
		HTTPMethods: "DELETE",
		Path:       "/api/o/:orgId/analytics/:id",
	},
}

var userPerms map[string]role.Permission = map[string]role.Permission{
	create: {
		Name:       "user:create",
		HTTPMethods: "POST",
		Path:       "/api/o/:orgId/users",
	},
	read: {
		Name:       "user:read",
		HTTPMethods: "GET",
		Path:       "/api/o/:orgId/users/find",
	},
	update: {
		Name:       "user:update",
		HTTPMethods: "PUT",
		Path:       "/api/o/:orgId/users/:id",
	},
	delete: {
		Name:       "user:delete",
		HTTPMethods: "DELETE",
		Path:       "/api/o/:orgId/users/:id",
	},
	invite: {
		Name:       "user:invite",
		HTTPMethods: "GET",
		Path:       "/api/o/:orgId/users/invite/:email/:roleId",
	},
	updateProfile: {
		Name:       "user:update-profile",
		HTTPMethods: "PUT",
		Path:       "/api/o/:orgId/users/update-profile",
	},
	uploadAvatar: {
		Name:       "user:upload-avatar",
		HTTPMethods: "PUT",
		Path:       "/api/o/:orgId/users/upload-avatar",
	},
	downloadAvatar: {
		Name:       "user:download-avatar",
		HTTPMethods: "GET",
		Path:       "/api/o/:orgId/users/download-avatar",
	},
	roleCount: {
		Name:       "user:role-counts",
		HTTPMethods: "GET",
		Path:       "/api/o/:orgId/users/role-counts",
	},
}

var contentPerms map[string]role.Permission = map[string]role.Permission{
	create: {
		Name:       "content:create",
		HTTPMethods: "POST",
		Path:       "/api/o/:orgId/contents",
	},
	read: {
		Name:       "content:read",
		HTTPMethods: "GET",
		Path:       "/api/o/:orgId/contents",
	},
	update: {
		Name:       "content:update",
		HTTPMethods: "PUT",
		Path:       "/api/o/:orgId/contents/:id",
	},
	delete: {
		Name:       "content:delete",
		HTTPMethods: "DELETE",
		Path:       "/api/o/:orgId/contents/:id",
	},
}

var contentLibPerms map[string]role.Permission = map[string]role.Permission{
	create: {
		Name:       "content-library:create",
		HTTPMethods: "POST",
		Path:       "/api/o/:orgId/content-libs",
	},
	read: {
		Name:       "content-library:read",
		HTTPMethods: "GET",
		Path:       "/api/o/:orgId/content-libs",
	},
	update: {
		Name:       "content-library:update",
		HTTPMethods: "PUT",
		Path:       "/api/o/:orgId/content-libs/:id",
	},
	delete: {
		Name:       "content-library:delete",
		HTTPMethods: "DELETE",
		Path:       "/api/o/:orgId/content-libs/:id",
	},
}

var projectPerms map[string]role.Permission = map[string]role.Permission{
	create: {
		Name:       "project:create",
		HTTPMethods: "POST",
		Path:       "/api/o/:orgId/projects",
	},
	read: {
		Name:       "project:read",
		HTTPMethods: "GET",
		Path:       "/api/o/:orgId/projects",
	},
	update: {
		Name:       "project:update",
		HTTPMethods: "PUT",
		Path:       "/api/o/:orgId/projects/:id",
	},
	delete: {
		Name:       "project:delete",
		HTTPMethods: "DELETE",
		Path:       "/api/o/:orgId/projects/:id",
	},
}

var bountyPerms map[string]role.Permission = map[string]role.Permission{
	create: {
		Name:       "bounty:create",
		HTTPMethods: "POST",
		Path:       "/api/o/:orgId/bounties",
	},
	read: {
		Name:       "bounty:read",
		HTTPMethods: "GET",
		Path:       "/api/o/:orgId/bounties",
	},
	update: {
		Name:       "bounty:update",
		HTTPMethods: "PUT",
		Path:       "/api/o/:orgId/bounties/:id",
	},
	delete: {
		Name:       "bounty:delete",
		HTTPMethods: "DELETE",
		Path:       "/api/o/:orgId/bounties/:id",
	},
}

var productServicePerms map[string]role.Permission = map[string]role.Permission{
	create: {
		Name:       "product-service:create",
		HTTPMethods: "POST",
		Path:       "/api/o/:orgId/product-services",
	},
	read: {
		Name:       "product-service:read",
		HTTPMethods: "GET",
		Path:       "/api/o/:orgId/product-services",
	},
	update: {
		Name:       "product-service:update",
		HTTPMethods: "PUT",
		Path:       "/api/o/:orgId/product-services/:id",
	},
	delete: {
		Name:       "product-service:delete",
		HTTPMethods: "DELETE",
		Path:       "/api/o/:orgId/product-services/:id",
	},
}

var messagingPerms map[string]role.Permission = map[string]role.Permission{
	create: {
		Name:       "message:create",
		HTTPMethods: "POST",
		Path:       "/api/o/:orgId/msgs",
	},
	read: {
		Name:       "message:read",
		HTTPMethods: "GET",
		Path:       "/api/o/:orgId/msgs",
	},
	update: {
		Name:       "message:update",
		HTTPMethods: "PUT",
		Path:       "/api/o/:orgId/msgs/:id",
	},
	delete: {
		Name:       "message:delete",
		HTTPMethods: "DELETE",
		Path:       "/api/o/:orgId/msgs/:id",
	},
}


var meetingPerms map[string]role.Permission = map[string]role.Permission{
	create: {
		Name:       "meeting:create",
		HTTPMethods: "POST",
		Path:       "/api/o/:orgId/meetings",
	},
	read: {
		Name:       "meeting:read",
		HTTPMethods: "GET",
		Path:       "/api/o/:orgId/meetings",
	},
	update: {
		Name:       "meeting:update",
		HTTPMethods: "PUT",
		Path:       "/api/o/:orgId/meetings/:id",
	},
	delete: {
		Name:       "meeting:delete",
		HTTPMethods: "DELETE",
		Path:       "/api/o/:orgId/meetings/:id",
	},
}

var activityPerms map[string]rolesvc.Permission = map[string]rolesvc.Permission{
	activityDashboardStats: {
		Name:       "activity:dashboard-stats",
		HTTPMethods: "GET",
		Path:       "/api/o/:orgId/activities/dashboard-stats",
	},
}