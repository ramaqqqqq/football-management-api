package entity

type Team struct {
	ModelID
	ModelLogTime
	Name        string `db:"name"`
	Logo        string `db:"logo"`
	YearFounded int    `db:"year_founded"`
	Address     string `db:"address"`
	City        string `db:"city"`
}
