package schemes

import (
	"lab1/internal/app/ds"

	"mime/multipart"
	"time"
)

type ModuleRequest struct {
	ModuleId string `uri:"module_id" binding:"required,uuid"`
}

// вопрос
type GetAllModulesRequest struct {
	Name string `form:"name"`
}

// вопрос
type AddModuleRequest struct {
	ds.Module
	Image *multipart.FileHeader `form:"image" json:"image" binding:"required"`
}

type ChangeModuleRequest struct {
	ModuleId    string                `uri:"module_id" binding:"required,uuid"`
	Name        *string               `form:"name" json:"name" binding:"omitempty,max=100"`
	Description *string               `form:"description" json:"description" binding:"omitempty,max=75"`
	Mass        *string               `form:"mass" json:"mass"`
	Image       *multipart.FileHeader `form:"image" json:"image"`
	Diameter    *string               `form:"diameter" json:"diameter" binding:"omitempty,max=100"`
	Length      *string               `form:"length" json:"length" binding:"omitempty,max=100"`
}

type AddToMissionRequest struct {
	ModuleId string `uri:"module_id" binding:"required,uuid"`
}

type GetAllMissionsRequest struct {
	DateApproveStart *time.Time `form:"date_approve_start" json:"date_approve_start" time_format:"2006-01-02"`
	DateApproveEnd   *time.Time `form:"date_approve_end" json:"date_approve_end" time_format:"2006-01-02"`
	Status           string     `form:"status"`
}

type MissionRequest struct {
	MissionId string `uri:"mission_id" binding:"required,uuid"`
}

type UpdateMissionRequest struct {
	URI struct {
		MissionId string `uri:"mission_id" binding:"required,uuid"`
	}
	Name             string     `form:"name" json:"name" binding:"required,max=50"`
	DateStartMission *time.Time `form:"date_start_mission" json:"date_start_mission" time_format:"2006-01-02"`
	Description      string     `form:"description" json:"description" binding:"required,max=100"`
}

type DeleteFromMissionRequest struct {
	MissionId string `uri:"mission_id" binding:"required,uuid"`
	ModuleId  string `uri:"module_id" binding:"required,uuid"`
}

type UserConfirmRequest struct {
	MissionId string `uri:"mission_id" binding:"required,uuid"`
}

type ModeratorConfirmRequest struct {
	URI struct {
		MissionId string `uri:"mission_id" binding:"required,uuid"`
	}
	Status string `form:"status" json:"status" binding:"required"`
}
