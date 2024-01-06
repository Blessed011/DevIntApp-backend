package ds

import (
	"lab1/internal/app/role"
	"time"
)

const StatusDraft string = "черновик"
const StatusFormed string = "сформирована"
const StatusCompleted string = "завершена"
const StatusRejected string = "отклонена"
const StatusDeleted string = "удалена"

const FundingApproved string = "финансирование одобрено"
const FundingRejected string = "финансирование отклонено"
const FundingOnConsideration string = "финансирование на рассмотрении"

type Module struct {
	UUID        string  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"  json:"uuid" binding:"-"`
	Name        string  `gorm:"size:50;not null" json:"name"`
	IsDeleted   bool    `gorm:"not null;default:false" json:"-" binding:"-"`
	ImageURL    *string `gorm:"size:100" json:"image_url" binding:"-"`
	Description string  `gorm:"size:100" form:"description" json:"description" binding:"required"`
	Mass        string  `gorm:"size:15" form:"mass" json:"mass" binding:"required"`
	Diameter    string  `gorm:"size:15" form:"diameter" json:"diameter" binding:"required"`
	Length      string  `gorm:"size:15" form:"length" json:"length" binding:"required"`
}

type Mission struct {
	UUID          string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name          *string    `gorm:"size:50"`
	DateCreated   time.Time  `gorm:"not null;type:timestamp"`
	DateApprove   *time.Time `gorm:"type:timestamp"`
	DateEnd       *time.Time `gorm:"type:timestamp"`
	Status        string     `gorm:"size:20;not null"`
	ModeratorId   *string    `json:"-"`
	CustomerId    string     `gorm:not null`
	FundingStatus *string    `gorm:"size:40"`

	Moderator *User
	Customer  User
}

type User struct {
	UUID     string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"-"`
	Role     role.Role
	Login    string `gorm:"size:40;not null" json:"login"`
	Password string `gorm:"size:45;not null" json:"-"`
}

type Flight struct {
	ModuleId  string   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"module_id"`
	MissionId string   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"mission_id"`
	Module    *Module  `gorm:"foreignKey:ModuleId" json:"module"`
	Mission   *Mission `gorm:"foreignKey:MissionId" json:"mission"`
}
