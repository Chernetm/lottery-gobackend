package models

type AdminStats struct {
	TotalRevenue    float64 `json:"totalRevenue"`
	ActiveLotteries int64   `json:"activeLotteries"`
	TotalUsers      int64   `json:"totalUsers"`
}
