package contract

type CreateTeamRequest struct {
	Name        string `json:"name" form:"name" binding:"required"`
	Logo        string `json:"logo" form:"logo"`
	YearFounded int    `json:"year_founded" form:"year_founded" binding:"required,min=1800,max=2100"`
	Address     string `json:"address" form:"address"`
	City        string `json:"city" form:"city" binding:"required"`
}

type UpdateTeamRequest struct {
	Name        string `json:"name" form:"name"`
	Logo        string `json:"logo" form:"logo"`
	YearFounded int    `json:"year_founded" form:"year_founded" binding:"omitempty,min=1800,max=2100"`
	Address     string `json:"address" form:"address"`
	City        string `json:"city" form:"city"`
}

type TeamResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Logo        string `json:"logo"`
	YearFounded int    `json:"year_founded"`
	Address     string `json:"address"`
	City        string `json:"city"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type TeamBrief struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Logo string `json:"logo"`
}
