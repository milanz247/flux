package controllers

import (
	"fmt"
	"time"

	"flux/app/dto"
	"flux/app/services"
	"flux/framework"
)

// DashboardController renders the main dashboard.
type DashboardController struct {
	users *services.UserService
}

func NewDashboardController(users *services.UserService) *DashboardController {
	return &DashboardController{users: users}
}

// Index — GET /dashboard
func (c *DashboardController) Index(req *framework.Request) {
	ctx := req.Context()

	total, err := c.users.CountUsers(ctx)
	if err != nil {
		req.Error(err)
		return
	}
	verified, err := c.users.CountVerifiedUsers(ctx)
	if err != nil {
		req.Error(err)
		return
	}
	newThisWeek, err := c.users.CountUsersSince(ctx, time.Now().AddDate(0, 0, -7))
	if err != nil {
		req.Error(err)
		return
	}
	recent, err := c.users.RecentUsers(ctx, 5)
	if err != nil {
		req.Error(err)
		return
	}

	req.View("Dashboard/Index", dto.DashboardDTO{
		Stats: []dto.DashboardStatDTO{
			{Label: "Total Users", Value: fmt.Sprint(total), Hint: "All registered accounts"},
			{Label: "Verified", Value: fmt.Sprint(verified), Hint: "Email address confirmed"},
			{Label: "New This Week", Value: fmt.Sprint(newThisWeek), Hint: "Signed up in the last 7 days"},
		},
		RecentUsers: recent,
	})
}
