package models

type Post struct {
	Shared `gorm:",inline"`
	UserID string `json:"user_id" gorm:"not null"`
	User   *User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Title  string `json:"title" gorm:"type:varchar(200);not null"`
	Body   string `json:"body" gorm:"type:text;not null"`
}
