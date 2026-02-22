package contract

type CreateMatchRequest struct {
	HomeTeamID int64  `json:"home_team_id" binding:"required"`
	AwayTeamID int64  `json:"away_team_id" binding:"required"`
	MatchDate  string `json:"match_date" binding:"required"` // YYYY-MM-DD
	MatchTime  string `json:"match_time" binding:"required"` // HH:MM
}

type UpdateMatchRequest struct {
	HomeTeamID int64  `json:"home_team_id"`
	AwayTeamID int64  `json:"away_team_id"`
	MatchDate  string `json:"match_date"`
	MatchTime  string `json:"match_time"`
}

type GoalInput struct {
	PlayerID   int64 `json:"player_id" binding:"required"`
	GoalMinute int   `json:"goal_minute" binding:"required,min=1,max=120"`
}

type SubmitResultRequest struct {
	HomeScore int         `json:"home_score" binding:"gte=0"`
	AwayScore int         `json:"away_score" binding:"gte=0"`
	Goals     []GoalInput `json:"goals"`
}

type GoalDetail struct {
	ID         int64  `json:"id"`
	PlayerID   int64  `json:"player_id"`
	PlayerName string `json:"player_name"`
	GoalMinute int    `json:"goal_minute"`
}

type MatchResponse struct {
	ID        int64      `json:"id"`
	HomeTeam  TeamBrief  `json:"home_team"`
	AwayTeam  TeamBrief  `json:"away_team"`
	MatchDate string     `json:"match_date"`
	MatchTime string     `json:"match_time"`
	HomeScore *int       `json:"home_score"`
	AwayScore *int       `json:"away_score"`
	Status    string     `json:"status"`
	Goals     []GoalDetail `json:"goals,omitempty"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
}

type TopScorerInfo struct {
	PlayerID int64  `json:"player_id"`
	Name     string `json:"name"`
	Goals    int    `json:"goals"`
}

type MatchReportResponse struct {
	MatchID           int64          `json:"match_id"`
	MatchDate         string         `json:"match_date"`
	MatchTime         string         `json:"match_time"`
	HomeTeam          TeamBrief      `json:"home_team"`
	AwayTeam          TeamBrief      `json:"away_team"`
	HomeScore         int            `json:"home_score"`
	AwayScore         int            `json:"away_score"`
	FinalStatus       string         `json:"final_status"`
	TopScorer         *TopScorerInfo `json:"top_scorer"`
	HomeTeamTotalWins int            `json:"home_team_total_wins"`
	AwayTeamTotalWins int            `json:"away_team_total_wins"`
	Goals             []GoalDetail   `json:"goals"`
}
