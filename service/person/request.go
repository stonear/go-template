package person

type StoreRequest struct {
	Name string `form:"name" binding:"required,notEvil"`
}

type UpdateRequest struct {
	Name string `form:"name" binding:"required,notEvil"`
}
