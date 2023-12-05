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
// @Param		mission_id path string true "id миссии"
// @Success		200 {object} schemes.MissionResponse
// @Router		/api/missions/{mission_id} [get]
func (app *Application) GetMission(c *gin.Context) {
	var request schemes.MissionRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userId := getUserId(c)
	mission, err := app.repo.GetMissionById(request.MissionId, userId)
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
	Name             string     `json:"name"`
	DateStartMission *time.Time `json:"date_start_mission" time_format:"2006-01-02"`
	Description      string     `json:"description"`
}

// @Summary		Указать название, дату старта и описание миссии
// @Tags		Миссии
// @Description	Позволяет изменить название, дату старта и описание миссии и возвращает обновлённые данные
// @Access		json
// @Produce		json
// @Param		mission_id path string true "id миссии"
// @Param		name body SwaggerUpdateMissionRequest true "Название"
// @Param		date_start_mission body SwaggerUpdateMissionRequest true "Дата старта"
// @Param		description body SwaggerUpdateMissionRequest true "Описание"
// @Success		200 {object} schemes.UpdateMissionResponse
// @Router		/api/missions/{mission_id} [put]
func (app *Application) UpdateMission(c *gin.Context) {
	var request schemes.UpdateMissionRequest
	if err := c.ShouldBindUri(&request.URI); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userId := getUserId(c)
	mission, err := app.repo.GetMissionById(request.URI.MissionId, userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if mission == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("миссия не найдена"))
		return
	}
	mission.Name = request.Name
	mission.DateStartMission = *request.DateStartMission
	mission.Description = request.Description
	if app.repo.SaveMission(mission); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.UpdateMissionResponse{Mission: schemes.ConvertMission(mission)})
}

// @Summary		Удалить миссию
// @Tags		Миссии
// @Description	Удаляет миссию по id
// @Param		mission_id path string true "id миссии"
// @Success		200
// @Router		/api/missions/{mission_id} [delete]
func (app *Application) DeleteMission(c *gin.Context) {
	var request schemes.MissionRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userId := getUserId(c)
	mission, err := app.repo.GetMissionById(request.MissionId, userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if mission == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("миссия не найдена"))
		return
	}
	mission.Status = ds.DELETED

	if err := app.repo.SaveMission(mission); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

// @Summary		Удалить модуль из миссии
// @Tags		Миссии
// @Description	Удалить модуль из миссии
// @Produce		json
// @Param		mission_id path string true "id миссии"
// @Param		module_id path string true "id модуля"
// @Success		200 {object} schemes.AllModulesResponse
// @Router		/api/missions/{mission_id}/delete_module/{module_id} [delete]
func (app *Application) DeleteFromMission(c *gin.Context) {
	var request schemes.DeleteFromMissionRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userId := getUserId(c)
	mission, err := app.repo.GetMissionById(request.MissionId, userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if mission == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("миссия не найдена"))
		return
	}
	if mission.Status != ds.DRAFT {
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("нельзя редактировать миссию со статусом: %s", mission.Status))
		return
	}

	if err := app.repo.DeleteFromMission(request.MissionId, request.ModuleId); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	modules, err := app.repo.GetFlight(request.MissionId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.AllModulesResponse{Modules: modules})
}

// @Summary		Сформировать миссию
// @Tags		Миссии
// @Description	Сформировать или удалить миссию пользователем
// @Produce		json
// @Param		confirm body boolean true "подтвердить"
// @Success		200
// @Router		/api/missions/user_confirm [put]
func (app *Application) UserConfirm(c *gin.Context) {
	var request schemes.UserConfirmRequest
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

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
	if mission.Status != ds.DRAFT {
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("нельзя сформировать миссию со статусом %s", mission.Status))
		return
	}

	if request.Confirm {
		mission.Status = ds.FORMED
		now := time.Now()
		mission.DateApprove = &now
	} else {
		mission.Status = ds.DELETED
	}

	if err := app.repo.SaveMission(mission); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

// @Summary		Подтвердить миссию
// @Tags		Миссии
// @Description	Подтвердить или отменить миссию модератором
// @Produce		json
// @Param		mission_id path string true "id миссии"
// @Param		confirm body boolean true "подтвердить"
// @Success		200
// @Router		/api/missions/{mission_id}/moderator_confirm [put]
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
	mission, err := app.repo.GetMissionById(request.URI.MissionId, userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if mission == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("миссия не найдена"))
		return
	}
	if mission.Status != ds.FORMED {
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("нельзя изменить статус миссии с \"%s\" на \"%s\"", mission.Status, ds.FORMED))
		return
	}

	if request.Confirm {
		mission.Status = ds.COMPELTED
		now := time.Now()
		mission.DateEnd = &now
	} else {
		mission.Status = ds.REJECTED
	}
	mission.ModeratorId = &userId

	if err := app.repo.SaveMission(mission); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}
