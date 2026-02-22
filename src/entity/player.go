package entity

type PlayerPosition string

const (
	PlayerPositionForward    PlayerPosition = "penyerang"
	PlayerPositionMidfielder PlayerPosition = "gelandang"
	PlayerPositionDefender   PlayerPosition = "bertahan"
	PlayerPositionGoalkeeper PlayerPosition = "penjaga_gawang"
)

type Player struct {
	ModelID
	ModelLogTime
	TeamID       int64          `db:"team_id"`
	Name         string         `db:"name"`
	Height       float64        `db:"height"`
	Weight       float64        `db:"weight"`
	Position     PlayerPosition `db:"position"`
	JerseyNumber int            `db:"jersey_number"`
}
