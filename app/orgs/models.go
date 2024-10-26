package orgs

type AddOrgRequest struct {
	UserID int    `json:"-"`
	Name   string `json:"name"`
	Size   string `json:"size"`
}

type OrgResponse struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	OrgSlug string `json:"orgSlug"`
}

type FindOrgRequest struct {
	UserID int `json:"-"`
}

type OrgWithRole struct {
	OrgID   int    `json:"orgId"`
	RoleID  int    `json:"roleId"`
	Name    string `json:"name"`
	OrgSlug string `json:"orgSlug"`
	UserID  int    `json:"userId"`
}
