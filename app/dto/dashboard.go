package dto

// DashboardStatDTO is one KPI tile on the dashboard.
type DashboardStatDTO struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Hint  string `json:"hint"`
}

// DashboardDTO — props for Pages/Dashboard/Index.vue.
type DashboardDTO struct {
	Stats       []DashboardStatDTO `json:"stats"`
	RecentUsers []UserDTO          `json:"recentUsers"`
}
