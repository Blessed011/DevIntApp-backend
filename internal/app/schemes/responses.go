package schemes

import (
	"lab1/internal/app/ds"
	"time"
)

type AllModulesResponse struct {
	Modules []ds.Module `json:"modules"`
}

type MissionShort struct {
	UUID        string `json:"uuid"`
	ModuleCount int    `json:"module_count"`
}

type GetAllModulesResponse struct {
	DraftMission *MissionShort `json:"draft_mission"`
	Modules      []ds.Module   `json:"modules"`
}

type AllMissionsResponse struct {
	Missions []MissionOutput `json:"missions"`
}

type MissionResponse struct {
	Mission MissionOutput `json:"missions"`
	Modules []ds.Module   `json:"modules"`
}

type UpdateMissionResponse struct {
	Mission MissionOutput `json:"missions"`
}

type MissionOutput struct {
	UUID             string  `json:"uuid"`
	Name             string  `json:"name"`
	Status           string  `json:"status"`
	DateCreated      string  `json:"date_created"`
	DateApprove      *string `json:"date_approve"`
	DateEnd          *string `json:"date_end"`
	DateStartMission string  `json:"date_start_mission"`
	Description      string  `json:"description"`
	Moderator        *string `json:"moderator"`
	Customer         string  `json:"customer"`
}

func ConvertMission(mission *ds.Mission) MissionOutput {
	output := MissionOutput{
		UUID:             mission.UUID,
		Name:             mission.Name,
		Status:           mission.Status,
		DateCreated:      mission.DateCreated.Format("2006-01-02 15:04:05"),
		Description:      mission.Description,
		DateStartMission: mission.DateStartMission.Format("2006-01-02"),
		Customer:         mission.Customer.Login,
	}

	if mission.DateApprove != nil {
		dateApprove := mission.DateApprove.Format("2006-01-02 15:04:05")
		output.DateApprove = &dateApprove
	}

	if mission.DateEnd != nil {
		dateEnd := mission.DateEnd.Format("2006-01-02 15:04:05")
		output.DateEnd = &dateEnd
	}

	if mission.Moderator != nil {
		output.Moderator = &mission.Moderator.Login
	}

	return output
}

type LoginResp struct {
	ExpiresIn   time.Duration `json:"expires_in"`
	AccessToken string        `json:"access_token"`
	TokenType   string        `json:"token_type"`
}

type SwaggerLoginResp struct {
	ExpiresIn   int64  `json:"expires_in"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type RegisterResp struct {
	Ok bool `json:"ok"`
}
