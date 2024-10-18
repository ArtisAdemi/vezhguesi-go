package entities

type CreateEntityRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type EntityResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type GetEntityRequest struct {
	ID   uint   `json:"-"`
	Name string `json:"name"`
}
