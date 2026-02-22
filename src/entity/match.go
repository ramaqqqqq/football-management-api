package entity

import "time"

type MatchStatus string

const (
	MatchStatusScheduled MatchStatus = "scheduled"
	MatchStatusCompleted MatchStatus = "completed"
)

type Match struct {
	ModelID
	ModelLogTime
	HomeTeamID int64       `db:"home_team_id"`
	AwayTeamID int64       `db:"away_team_id"`
	MatchDate  time.Time   `db:"match_date"`
	MatchTime  string      `db:"match_time"`
	HomeScore  *int        `db:"home_score"`
	AwayScore  *int        `db:"away_score"`
	Status     MatchStatus `db:"status"`
}

type MatchWinStat struct {
	ID         int64  `db:"id"`
	HomeTeamID int64  `db:"home_team_id"`
	AwayTeamID int64  `db:"away_team_id"`
	HomeScore  *int   `db:"home_score"`
	AwayScore  *int   `db:"away_score"`
	Status     string `db:"status"`
}
