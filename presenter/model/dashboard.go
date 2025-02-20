package model

type DashboardRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
type DashboardResponse struct {
	TotalRoom        int     `json:"total_room"`
	TotalReservation int     `json:"total_reservation"`
	TotalVisitor     int     `json:"total_visitor"`
	TotalOmset       float64 `json:"total_omset"`
	Rooms            []struct {
		ID         int     `json:"id"`
		Name       string  `json:"name"`
		Percentage float64 `json:"percentage"`
		PriceHour  float64 `json:"price_hour"`
	}
}
