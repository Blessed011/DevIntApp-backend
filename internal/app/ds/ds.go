package ds

import (
	"time"
)

const DRAFT string = "черновик"
const FORMED string = "сформирован"
const COMPELTED string = "завершён"
const REJECTED string = "отклонён"
const DELETED string = "удалён"

type Module struct {
	UUID        string  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"  json:"uuid" binding:"-"`
	Name        string  `gorm:"type:varchar(50);not null" json:"name"`
	IsDeleted   bool    `gorm:"type:bool;not null;default:false" json:"-" binding:"-"`
	ImageURL    *string `gorm:"type:varchar(100)" json:"image_url" binding:"-"`
	Description string  `gorm:"type:text" form:"description" json:"description" binding:"required"`
	Mass        string  `gorm:"type:varchar(15)" form:"mass" json:"mass" binding:"required"`
	Diameter    string  `gorm:"type:varchar(15)" form:"diameter" json:"diameter" binding:"required"`
	Length      string  `gorm:"type:varchar(15)" form:"length" json:"length" binding:"required"`
}

type Mission struct {
	UUID             string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name             string     `gorm:"type:varchar(50);not null"`
	DateCreated      time.Time  `gorm:"type:date;not null;default:current_date"`
	DateApprove      *time.Time `gorm:"type:date"`
	DateEnd          *time.Time `gorm:"type:date"`
	DateStartMission time.Time  `gorm:"type:date"`
	Status           string     `gorm:"type:varchar(30)"`
	Description      string     `gorm:"type:text"`
	ModeratorId      *string    `json:"-"`
	CustomerId       string     `gorm:not null`

	Moderator *User
	Customer  User
}

type User struct {
	UUID      string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"uuid" binding:"-"`
	Name      string `gorm:"type:varchar(50)" json:"name"`
	Login     string `gorm:"size:30;not null" json:"-"`
	Password  string `gorm:"size:30;not null" json:"-"`
	Moderator bool   `gorm:"type:bool;not null" json:"-"`
}

type Flight struct {
	ModuleId  string   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"module_id"`
	MissionId string   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"mission_id"`
	Module    *Module  `gorm:"foreignKey:ModuleId" json:"module"`
	Mission   *Mission `gorm:"foreignKey:MissionId" json:"mission"`
}
