// Package controllers holds the HTTP layer. Controllers are thin structs:
// they read validated input from the framework Request, call a service, and
// render a page (req.View) or redirect. No business logic, no Gin, no GORM.
package controllers

import (
	"strconv"

	"flux/app/dto"
	"flux/app/services"
	"flux/framework"
)

// UserController is a full RESTful resource controller. Registered with
// Route.Resource("/users", ...) it serves all seven CRUD routes.
type UserController struct {
	service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{service: service}
}

// Index — GET /users
func (c *UserController) Index(req *framework.Request) {
	page, _ := strconv.Atoi(req.QueryDefault("page", "1"))
	search := req.Query("search")

	users, pagination, err := c.service.GetUsers(req.Context(), page, 10, search)
	if err != nil {
		req.Error(err)
		return
	}

	req.View("Users/Index", dto.UserIndexDTO{
		Users:      users,
		Pagination: pagination,
		Search:     search,
	})
}

// Create — GET /users/create
func (c *UserController) Create(req *framework.Request) {
	req.View("Users/Create", dto.UserCreateDTO{})
}

// Store — POST /users
func (c *UserController) Store(req *framework.Request) {
	var input dto.CreateUserDTO
	if !req.Validate(&input) {
		return
	}

	if _, err := c.service.CreateUser(req.Context(), input); err != nil {
		req.Error(err)
		return
	}

	req.Redirect("/users")
}

// Show — GET /users/:id
func (c *UserController) Show(req *framework.Request) {
	user, err := c.service.GetUser(req.Context(), req.ParamUint("id"))
	if err != nil {
		req.Error(err)
		return
	}

	req.View("Users/Show", dto.UserShowDTO{User: user})
}

// Edit — GET /users/:id/edit
func (c *UserController) Edit(req *framework.Request) {
	user, err := c.service.GetUser(req.Context(), req.ParamUint("id"))
	if err != nil {
		req.Error(err)
		return
	}

	req.View("Users/Edit", dto.UserEditDTO{User: user})
}

// Update — PUT /users/:id
func (c *UserController) Update(req *framework.Request) {
	var input dto.UpdateUserDTO
	if !req.Validate(&input) {
		return
	}

	if _, err := c.service.UpdateUser(req.Context(), req.ParamUint("id"), input); err != nil {
		req.Error(err)
		return
	}

	req.Redirect("/users")
}

// Destroy — DELETE /users/:id
func (c *UserController) Destroy(req *framework.Request) {
	if err := c.service.DeleteUser(req.Context(), req.ParamUint("id")); err != nil {
		req.Error(err)
		return
	}

	req.Redirect("/users")
}
