package schemes

import (
	"lab1/internal/app/ds"
	"mime/multipart"
	"time"
)

type ModuleRequest struct {
	ModuleId string `uri:"id" binding:"required,uuid"`
}

type GetAllModulesRequest struct {
	ModuleName string `form:"name"`
}

type AddModuleRequest struct {
	ds.Module
	Image *multipart.FileHeader `form:"image" json:"image"`
}

type ChangeModuleRequest struct {
	ModuleId    string                `uri:"id" binding:"required,uuid"`
	Name        *string               `form:"name" json:"name" binding:"omitempty,max=100"`
	Description *string               `form:"description" json:"description" binding:"omitempty,max=75"`
	Mass        *string               `form:"mass" json:"mass"`
	Image       *multipart.FileHeader `form:"image" json:"image"`
	Diameter    *string               `form:"diameter" json:"diameter" binding:"omitempty,max=100"`
	Length      *string               `form:"length" json:"length" binding:"omitempty,max=100"`
}

type AddToMissionRequest struct {
	ModuleId string `uri:"id" binding:"required,uuid"`
}

type GetAllMissionsRequest struct {
	FormationDateStart *time.Time `form:"formation_date_start" json:"formation_date_start" time_format:"2006-01-02 15:04"`
	FormationDateEnd   *time.Time `form:"formation_date_end" json:"formation_date_end" time_format:"2006-01-02 15:04"`
	Status             string     `form:"status" json:"status"`
}

type MissionRequest struct {
	MissionId string `uri:"id" binding:"required,uuid"`
}

type UpdateMissionRequest struct {
	Name string `form:"name" json:"name" binding:"required,max=50"`
}

type DeleteFromMissionRequest struct {
	ModuleId string `uri:"id" binding:"required,uuid"`
}

type ModeratorConfirmRequest struct {
	URI struct {
		MissionId string `uri:"id" binding:"required,uuid"`
	}
	Confirm *bool `form:"confirm" binding:"required"`
}

type LoginReq struct {
	Login    string `form:"login" binding:"required,max=30"`
	Password string `form:"password" binding:"required,max=30"`
}

type RegisterReq struct {
	Login    string `form:"login" binding:"required,max=30"`
	Password string `form:"password" binding:"required,max=30"`
}
