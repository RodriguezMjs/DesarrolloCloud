package dtos

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	Token string    `json:"token"`
	User  UserDTO   `json:"user"`
	Roles []RoleDTO `json:"roles"`
}

type UserDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type RoleDTO struct {
	//TODO: Validar si se debe devolver el nombre del rol por temas de seguridad o si se
	//homologa con el front
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
