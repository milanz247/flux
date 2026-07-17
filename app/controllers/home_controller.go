package controllers

import (
	"flux/app/dto"
	"flux/framework"
)

// HomeController serves the public marketing/landing page.
type HomeController struct{}

func NewHomeController() *HomeController {
	return &HomeController{}
}

// Welcome — GET /
func (c *HomeController) Welcome(req *framework.Request) {
	req.View("Welcome", dto.WelcomePageDTO{})
}
