package dto

// UserDTO is the public shape of a user sent to the frontend.
type UserDTO struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Verified  bool   `json:"verified"`
	CreatedAt string `json:"createdAt"`
}

// UserIndexDTO — props for Pages/Users/Index.vue.
type UserIndexDTO struct {
	Users      []UserDTO     `json:"users"`
	Pagination PaginationDTO `json:"pagination"`
	Search     string        `json:"search"`
}

// UserCreateDTO — props for Pages/Users/Create.vue.
type UserCreateDTO struct{}

// UserShowDTO — props for Pages/Users/Show.vue.
type UserShowDTO struct {
	User UserDTO `json:"user"`
}

// UserEditDTO — props for Pages/Users/Edit.vue.
type UserEditDTO struct {
	User UserDTO `json:"user"`
}

// CreateUserDTO is the validated input for creating a user.
type CreateUserDTO struct {
	Name     string `json:"name" validate:"required,min=2,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

// UpdateUserDTO is the validated input for updating a user. Password is
// optional — leave it empty to keep the current one.
type UpdateUserDTO struct {
	Name     string `json:"name" validate:"required,min=2,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"omitempty,min=8,max=72"`
}
