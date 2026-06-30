package models

type MovieStatus string

const (
	MovieStatusPending    MovieStatus = "PENDING"
	MovieStatusProcessing MovieStatus = "PROCESSING"
	MovieStatusReady      MovieStatus = "READY"
	MovieStatusFailed     MovieStatus = "FAILED"
)

type Movie struct {
	ID              int64       `gorm:"column:id"`
	Title           string      `gorm:"column:title"`
	FileName        string      `gorm:"column:filename"`
	Status          MovieStatus `gorm:"column:status"`
	DurationSeconds int         `gorm:"column:duration_seconds"`
	HLSPlaylistPath *string     `gorm:"column:hls_playlist_path"`
}

//  For Gorm package
// func (Movie) TableName() string {
// 	return "movies"
// }
