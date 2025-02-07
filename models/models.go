package models

import (
	"github.com/google/uuid"
	"time"
)

type (
	Models interface {
		Initialize(id uuid.UUID, now time.Time)
		GetId() string
		SetID(id uuid.UUID)
		SetUpdatedAt()
		GetVersion() uint
		SetVersion(v uint)
	}

	AccountInfo struct {
		Id       string `json:"id"`
		FullName string `json:"full_name"`
		Email    string `json:"email"`
	}
)

type Shared struct {
	ID        uuid.UUID  `json:"id" gorm:"_id"`
	CreatedAt *time.Time `json:"created_at" gorm:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"deleted_at"`
	Version   uint       `json:"version" gorm:"_version"`
}

var _ Models = &Shared{}

func (m Shared) GetId() string {
	return m.ID.String()
}

func (m Shared) SetID(id uuid.UUID) {
	m.ID = id
}

func (m Shared) GetVersion() uint {
	return m.Version
}

func (m Shared) SetVersion(v uint) {
	m.Version = v
}

func (m Shared) SetUpdatedAt() {
	t := time.Now().UTC()
	m.UpdatedAt = &t
}

func (m Shared) Initialize(id uuid.UUID, now time.Time) {
	m.ID = id
	t := now.UTC()
	m.CreatedAt = &t
}
