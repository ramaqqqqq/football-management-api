package entity

type Goal struct {
	ModelID
	ModelLogTime
	MatchID    int64 `db:"match_id"`
	PlayerID   int64 `db:"player_id"`
	GoalMinute int   `db:"goal_minute"`
}
