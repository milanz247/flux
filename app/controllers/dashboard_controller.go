package controllers

import (
	"flux/app/dto"
	"flux/framework"
)

// DashboardController renders the main dashboard.
type DashboardController struct{}

func NewDashboardController() *DashboardController {
	return &DashboardController{}
}

// Index — GET /dashboard
func (c *DashboardController) Index(req *framework.Request) {
	req.View("Dashboard/Index", dto.DashboardDTO{})
}
