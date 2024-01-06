package app

import (
	"fmt"
	"lab1/internal/app/ds"
	"lab1/internal/app/role"
	"lab1/internal/app/schemes"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary		Получить все миссии
// @Tags		Миссии
// @Description	Возвращает все миссии с фильтрацией по статусу и дате формирования
// @Produce		json
// @Param		status query string false "статус перевозки"
// @Param		date_approve_start query string false "начальная дата формирования"
// @Param		date_approve_end query string false "конечная дата формирвания"
// @Success		200 {object} schemes.AllMissionsResponse
// @Router		/api/missions [get]
func (app *Application) GetAllMissions(c *gin.Context) {
	var request schemes.GetAllMissionsRequest
	var err error
	if err = c.ShouldBindQuery(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userId := getUserId(c)
	userRole := getUserRole(c)
	fmt.Println(userId, userRole)
	var missions []ds.Mission
	if userRole == role.Customer {
		missions, err = app.repo.GetAllMissions(&userId, request.DateApproveStart, request.DateApproveEnd, request.Status)
	} else {
		missions, err = app.repo.GetAllMissions(nil, request.DateApproveStart, request.DateApproveEnd, request.Status)
	}
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	outputMissions := make([]schemes.MissionOutput, len(missions))
	for i, mission := range missions {
		outputMissions[i] = schemes.ConvertMission(&mission)
	}
	c.JSON(http.StatusOK, schemes.AllMissionsResponse{Missions: outputMissions})
}

// @Summary		Получить одну миссию
// @Tags		Миссии
// @Description	Возвращает подробную информацию о миссии
// @Produce		json
// @Param		id path string true "id миссии"
// @Success		200 {object} schemes.MissionResponse
// @Router		/api/missions/{id} [get]
func (app *Application) GetMission(c *gin.Context) {
	var request schemes.MissionRequest
	var err error
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userId := getUserId(c)
	userRole := getUserRole(c)
	var mission *ds.Mission
	if userRole == role.Moderator {
		mission, err = app.repo.GetMissionById(request.MissionId, nil)
	} else {
		mission, err = app.repo.GetMissionById(request.MissionId, &userId)
	}
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if mission == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("миссия не найдена"))
		return
	}

	modules, err := app.repo.GetFlight(request.MissionId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, schemes.MissionResponse{Mission: schemes.ConvertMission(mission), Modules: modules})
}

type SwaggerUpdateMissionRequest struct {
	Name string `json:"name"`
}

// @Summary		Указать название
// @Tags		Миссии
// @Description	Позволяет изменить название миссии и возвращает обновлённые данные
// @Access		json
// @Produce		json
// @Param		name body SwaggerUpdateMissionRequest true "Название"
// @Success		200
// @Router		/api/missions [put]
func (app *Application) UpdateMission(c *gin.Context) {
	var request schemes.UpdateMissionRequest
	var err error
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var mission *ds.Mission
	userId := getUserId(c)
	mission, err = app.repo.GetDraftMission(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if mission == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("миссия не найдена"))
		return
	}
	mission.Name = &request.Name
	if app.repo.SaveMission(mission); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

// @Summary		Удалить черновую миссию
// @Tags		Миссии
// @Description	Удаляет черновую миссию
// @Success		200
// @Router		/api/missions [delete]
func (app *Application) DeleteMission(c *gin.Context) {
	var err error

	// Получить черновую заявку
	var mission *ds.Mission
	userId := getUserId(c)
	mission, err = app.repo.GetDraftMission(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if mission == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("миссия не найдена"))
		return
	}

	mission.Status = ds.StatusDeleted

	if err := app.repo.SaveMission(mission); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

// @Summary		Удалить модуль из черновой миссии
// @Tags		Миссии
// @Description	Удалить модуль из черновой миссии
// @Produce		json
// @Param		id path string true "id модуля"
// @Success		200
// @Router		/api/missions/delete_module/{id} [delete]
func (app *Application) DeleteFromMission(c *gin.Context) {
	var request schemes.DeleteFromMissionRequest
	var err error
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var mission *ds.Mission
	userId := getUserId(c)
	mission, err = app.repo.GetDraftMission(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if mission == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("миссия не найдена"))
		return
	}

	if err := app.repo.DeleteFromMission(mission.UUID, request.ModuleId); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

// @Summary		Сформировать миссию
// @Tags		Миссии
// @Description	Сформировать миссию пользователем
// @Success		200
// @Router		/api/missions/user_confirm [put]
func (app *Application) UserConfirm(c *gin.Context) {
	userId := getUserId(c)
	mission, err := app.repo.GetDraftMission(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if mission == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("миссия не найдена"))
		return
	}

	if err := fundingRequest(mission.UUID); err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf(`funding is impossible: {%s}`, err))
		return
	}

	fundingStatus := ds.FundingOnConsideration
	mission.FundingStatus = &fundingStatus

	mission.Status = ds.StatusFormed
	now := time.Now()
	mission.DateApprove = &now

	if err := app.repo.SaveMission(mission); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

// @Summary		Подтвердить миссию
// @Tags		Миссии
// @Description	Подтвердить или отменить миссию модератором
// @Param		id path string true "id миссии"
// @Param		confirm body boolean true "подтвердить"
// @Success		200
// @Router		/api/missions/{id}/moderator_confirm [put]
func (app *Application) ModeratorConfirm(c *gin.Context) {
	var request schemes.ModeratorConfirmRequest
	if err := c.ShouldBindUri(&request.URI); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userId := getUserId(c)
	mission, err := app.repo.GetMissionById(request.URI.MissionId, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if mission == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("миссия не найдена"))
		return
	}
	if mission.Status != ds.StatusFormed {
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("нельзя изменить статус миссии с \"%s\" на \"%s\"", mission.Status, ds.StatusFormed))
		return
	}

	if *request.Confirm {
		mission.Status = ds.StatusCompleted
		now := time.Now()
		mission.DateEnd = &now
	} else {
		mission.Status = ds.StatusRejected
	}
	moderator, err := app.repo.GetUserById(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	mission.ModeratorId = &userId
	mission.Moderator = moderator

	if err := app.repo.SaveMission(mission); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func (app *Application) Funding(c *gin.Context) {
	var request schemes.FundingReq
	var Token = "secret_token"

	if err := c.ShouldBindUri(&request.URI); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	fmt.Println(request, app.config.Token)

	if request.Token != Token {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	mission, err := app.repo.GetMissionById(request.URI.MissionId, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if mission == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("миссия не найдена"))
		return
	}
	// if mission.Status != ds.StatusFormed || *mission.FundingStatus != ds.FundingOnConsideration {
	// 	c.AbortWithStatus(http.StatusMethodNotAllowed)
	// 	return
	// }

	var fundingStatus string
	if *request.FundingStatus {
		fundingStatus = ds.FundingApproved
	} else {
		fundingStatus = ds.FundingRejected
	}
	mission.FundingStatus = &fundingStatus

	if err := app.repo.SaveMission(mission); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}
