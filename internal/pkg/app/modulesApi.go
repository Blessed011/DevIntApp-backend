package app

import (
	"fmt"
	_ "lab1/docs"
	"lab1/internal/app/ds"
	"lab1/internal/app/schemes"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary		Получить все модули
// @Tags		Модули
// @Description	Возвращает все доступные модули с опциональной фильтрацией по названию
// @Produce		json
// @Param		name query string false "название модуля для фильтрации"
// @Success		200 {object} schemes.GetAllModulesResponse
// @Router		/api/modules [get]
func (app *Application) GetAllModules(c *gin.Context) {
	var request schemes.GetAllModulesRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	modules, err := app.repo.GetModulesByName(request.ModuleName)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	response := schemes.GetAllModulesResponse{DraftMission: nil, Modules: modules}
	if userId, exists := c.Get("userId"); exists {
		draftMission, err := app.repo.GetDraftMission(userId.(string))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if draftMission != nil {
			response.DraftMission = &draftMission.UUID
		}
	}
	c.JSON(http.StatusOK, response)
}

// @Summary		Получить один модуль
// @Tags		Модули
// @Description	Возвращает более подробную информацию об одном модуле
// @Produce		json
// @Param		id path string true "id модуля"
// @Success		200 {object} ds.Module
// @Router		/api/modules/{id} [get]
func (app *Application) GetModule(c *gin.Context) {
	var request schemes.ModuleRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	module, err := app.repo.GetModuleByID(request.ModuleId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if module == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("модуль не найден"))
		return
	}
	c.JSON(http.StatusOK, module)
}

// @Summary		Удалить модуль
// @Tags		Модули
// @Description	Удаляет модуль по id
// @Param		id path string true "id модуля"
// @Success		200
// @Router		/api/modules/{id} [delete]
func (app *Application) DeleteModule(c *gin.Context) {
	var request schemes.ModuleRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	module, err := app.repo.GetModuleByID(request.ModuleId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if module == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("модуль не найден"))
		return
	}
	if module.ImageURL != nil {
		if err := app.deleteImage(c, module.UUID); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
	module.ImageURL = nil
	module.IsDeleted = true
	if err := app.repo.SaveModule(module); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

// @Summary		Добавить модуль
// @Tags		Модули
// @Description	Добавить новый модуль
// @Accept		mpfd
// @Param		image formData file false "Изображение модуля"
// @Param		name formData string true "Название" format:"string" maxLength:50
// @Param		description formData string true "Описание" format:"string" maxLength:100
// @Param		mass formData string true "Масса" format:"string" maxLength:10
// @Param		length formData int true "Длина" format:"string" maxLength:10
// @Param		diameter formData int true "Диаметр" format:"string" maxLength:10
// @Success		200
// @Router		/api/modules [post]
func (app *Application) AddModule(c *gin.Context) {
	var request schemes.AddModuleRequest
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	module := ds.Module(request.Module)
	if err := app.repo.AddModule(&module); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if request.Image != nil {
		imageURL, err := app.uploadImage(c, request.Image, module.UUID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		module.ImageURL = imageURL
	}
	if err := app.repo.SaveModule(&module); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusCreated)
}

// @Summary		Изменить модуль
// @Tags		Модули
// @Description	Изменить данные полей о модуле
// @Accept		mpfd
// @Param		id path string true "Идентификатор модуля" format:"uuid"
// @Param		image formData file false "Изображение модуля"
// @Param		name formData string true "Название" format:"string" maxLength:50
// @Param		description formData string true "Описание" format:"string" maxLength:100
// @Param		mass formData string true "Масса" format:"string" maxLength:10
// @Param		length formData int true "Длина" format:"string" maxLength:10
// @Param		diameter formData int true "Диаметр" format:"string" maxLength:10
// @Success		200
// @Router		/api/modules/{id} [put]
func (app *Application) ChangeModule(c *gin.Context) {
	var request schemes.ChangeModuleRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	module, err := app.repo.GetModuleByID(request.ModuleId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if module == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("модуль не найден"))
		return
	}
	if request.Name != nil {
		module.Name = *request.Name
	}
	if request.Description != nil {
		module.Description = *request.Description
	}
	if request.Mass != nil {
		module.Mass = *request.Mass
	}
	if request.Length != nil {
		module.Length = *request.Length
	}
	if request.Diameter != nil {
		module.Diameter = *request.Diameter
	}
	if request.Image != nil {
		if module.ImageURL != nil {
			if err := app.deleteImage(c, module.UUID); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}
		imageURL, err := app.uploadImage(c, request.Image, module.UUID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		module.ImageURL = imageURL
	}

	if err := app.repo.SaveModule(module); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

// @Summary		Добавить в миссию
// @Tags		Модули
// @Description	Добавить выбранный модуль в черновик миссии
// @Param		id path string true "id модуля"
// @Success		200
// @Router		/api/modules/{id}/add_to_mission [post]
func (app *Application) AddToMission(c *gin.Context) {
	var request schemes.AddToMissionRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var err error

	module, err := app.repo.GetModuleByID(request.ModuleId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if module == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("модуль не найден"))
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
		mission, err = app.repo.CreateDraftMission(userId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	if err = app.repo.AddToMission(mission.UUID, request.ModuleId); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}
