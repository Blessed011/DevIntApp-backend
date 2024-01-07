package schemes

import (
	"fmt"
	"lab1/internal/app/ds"
)

type AllModulesResponse struct {
	Modules []ds.Module `json:"modules"`
}

type GetAllModulesResponse struct {
	DraftMission *string     `json:"draft_mission"`
	Modules      []ds.Module `json:"modules"`
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
	UUID           string  `json:"uuid"`
	Name           *string `json:"name"`
	Status         string  `json:"status"`
	CreationDate   string  `json:"creation_date"`
	FormationDate  *string `json:"formation_date"`
	CompletionDate *string `json:"completion_date"`
	Moderator      *string `json:"moderator"`
	Customer       string  `json:"customer"`
}

func ConvertMission(mission *ds.Mission) MissionOutput {
	output := MissionOutput{
		UUID:         mission.UUID,
		Name:         mission.Name,
		Status:       mission.Status,
		CreationDate: mission.CreationDate.Format("2006-01-02T15:04:05Z07:00"),
		Customer:     mission.Customer.Login,
	}

	if mission.FormationDate != nil {
		formationDate := mission.FormationDate.Format("2006-01-02T15:04:05Z07:00")
		output.FormationDate = &formationDate
	}

	if mission.CompletionDate != nil {
		completionDate := mission.CompletionDate.Format("2006-01-02T15:04:05Z07:00")
		output.CompletionDate = &completionDate
	}

	if mission.Moderator != nil {
		fmt.Println(mission.Moderator.Login)
		output.Moderator = &mission.Moderator.Login
		fmt.Println(*output.Moderator)
	}

	return output
}

// type LoginResp struct {
// 	ExpiresIn   time.Duration `json:"expires_in"`
// 	AccessToken string        `json:"access_token"`
// 	TokenType   string        `json:"token_type"`
// }

// type SwaggerLoginResp struct {
// 	ExpiresIn   int64  `json:"expires_in"`
// 	AccessToken string `json:"access_token"`
// 	TokenType   string `json:"token_type"`
// }

// type RegisterResp struct {
// 	Ok bool `json:"ok"`
// }

type AddToMissionResp struct {
	ModulesCount int64 `json:"module_count"`
}
type AuthResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}
