package models

type User struct {
	Shared   `gorm:",inline"`
	FullName string   `json:"full_name" gorm:"type:varchar(100);not null"`
	Email    string   `json:"email" gorm:"type:varchar(100);uniqueIndex;not null"`
	Address  *Address `json:"address" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Posts    []Post   `json:"posts,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
