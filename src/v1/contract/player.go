package contract

type CreatePlayerRequest struct {
	TeamID       int64   `json:"team_id" binding:"required"`
	Name         string  `json:"name" binding:"required"`
	Height       float64 `json:"height" binding:"required,gt=0"`
	Weight       float64 `json:"weight" binding:"required,gt=0"`
	Position     string  `json:"position" binding:"required,oneof=penyerang gelandang bertahan penjaga_gawang"`
	JerseyNumber int     `json:"jersey_number" binding:"required,min=1,max=99"`
}

type UpdatePlayerRequest struct {
	TeamID       int64   `json:"team_id"`
	Name         string  `json:"name"`
	Height       float64 `json:"height" binding:"omitempty,gt=0"`
	Weight       float64 `json:"weight" binding:"omitempty,gt=0"`
	Position     string  `json:"position" binding:"omitempty,oneof=penyerang gelandang bertahan penjaga_gawang"`
	JerseyNumber int     `json:"jersey_number" binding:"omitempty,min=1,max=99"`
}

type PlayerResponse struct {
	ID           int64   `json:"id"`
	TeamID       int64   `json:"team_id"`
	Name         string  `json:"name"`
	Height       float64 `json:"height"`
	Weight       float64 `json:"weight"`
	Position     string  `json:"position"`
	JerseyNumber int     `json:"jersey_number"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}
