package url

import "time"

type URL struct {
	ID          uint
	ShortCode   string    `gorm:"uniqueIndex;not null;size:10"`
	OriginalURL string    `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}
