package orphanage

type OrphanageUserRequest struct {
	Name      string  `json:"name" form:"name" validate:"required,min=3"`
	Address   *string `json:"address" form:"address"`
	Age       *int    `json:"age" form:"age" `
	Education *string `json:"education" form:"education"`
	Position  *string `json:"position" form:"position"`
	Image     *string `json:"image" form:"image"`
}

type OrphanageUserResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Address   *string `json:"address"`
	Age       *int    `json:"age"`
	Education *string `json:"education"`
	Position  *string `json:"position"`
	Image     *string `json:"image"`
}
