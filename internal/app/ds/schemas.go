package ds

import (
	"time"
)

type Module struct {
	ModuleID    uint   `gorm:"primarykey"`
	Name        string `gorm:"type:varchar(50);not null"`
	IsDeleted   bool   `gorm:"type:bool;not null"`
	ImageURL    string `gorm:"type:varchar(100)"`
	Description string `gorm:"type:text"`
	Mass        string `gorm:"type:varchar(15)"`
	Diameter    string `gorm:"type:varchar(15)"`
	Length      string `gorm:"type:varchar(15)"`
}

type Mission struct {
	MissionID        uint      `gorm:"primarykey;not null"`
	Name             string    `gorm:"type:varchar(50);not null"`
	DateCreated      time.Time `gorm:"type:date;not null;default:current_date"`
	DateApprove      time.Time `gorm:"type:date"`
	DateEnd          time.Time `gorm:"type:date"`
	DateStartMission time.Time `gorm:"type:date"`
	Status           string    `gorm:"type:varchar(30)"`
	Description      string    `gorm:"type:text"`
}

type User struct {
	UserID    uint   `gorm:"primarykey"`
	Name      string `gorm:"type:varchar(50)"`
	Login     string `gorm:"size:30;not null"`
	Password  string `gorm:"size:30;not null"`
	Moderator bool   `gorm:"type:bool;not null"`
}

type Flight struct {
	ModuleID  uint     `gorm:"primarykey"`
	MissionID uint     `gorm:"primarykey"`
	Module    *Module  `gorm:"foreignkey:ModuleID"`
	Mission   *Mission `gorm:"foreignkey:MissionID"`
}
