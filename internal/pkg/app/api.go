package app

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"lab1/internal/app/ds"
	"lab1/internal/app/schemes"

	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

func (app *Application) uploadImage(c *gin.Context, image *multipart.FileHeader, UUID string) (*string, error) {
	src, err := image.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	extension := filepath.Ext(image.Filename)
	if extension != ".jpg" && extension != ".jpeg" {
		return nil, fmt.Errorf("разрешены только jpeg изображения")
	}
	imageName := UUID + extension

	_, err = app.minioClient.PutObject(c, app.config.BucketName, imageName, src, image.Size, minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
	if err != nil {
		return nil, err
	}
	imageURL := fmt.Sprintf("%s/%s/%s", app.config.MinioEndpoint, app.config.BucketName, imageName)
	return &imageURL, nil
}

// вопрос
func (app *Application) getCustomer() string {
	return "a4e0d78c-e12d-4e20-8826-13887e54b424"
}

func (app *Application) getModerator() *string {
	moderaorId := "a6911e14-7f64-40d7-bc19-7f9a171a9a2a"
	return &moderaorId
}

func (app *Application) GetAllModules(c *gin.Context) {
	var request schemes.GetAllModulesRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	modules, err := app.repo.GetModuleByName(request.Name)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	draftMission, err := app.repo.GetDraftMission(app.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	response := schemes.GetAllModulesResponse{DraftMission: nil, Modules: modules}
	if draftMission != nil {
		response.DraftMission = &schemes.MissionShort{UUID: draftMission.UUID}
		containers, err := app.repo.GetFlight(draftMission.UUID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		response.DraftMission.ModuleCount = len(containers)
	}
	c.JSON(http.StatusOK, response)
}

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
	module.IsDeleted = true
	if err := app.repo.SaveModule(module); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (app *Application) AddModule(c *gin.Context) {
	var request schemes.AddModuleRequest
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	log.Println(request)

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

	c.Status(http.StatusOK)
}

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
	if request.Image != nil {
		imageURL, err := app.uploadImage(c, request.Image, module.UUID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		module.ImageURL = imageURL
	}
	if request.Description != nil {
		module.Description = *request.Description
	}
	if request.Mass != nil {
		module.Mass = *request.Mass
	}
	if request.Diameter != nil {
		module.Diameter = *request.Diameter
	}
	if request.Length != nil {
		module.Length = *request.Length
	}

	if err := app.repo.SaveModule(module); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, module)
}

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
	mission, err = app.repo.GetDraftMission(app.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if mission == nil {
		mission, err = app.repo.CreateDraftMission(app.getCustomer())
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	if err = app.repo.AddToMission(mission.UUID, request.ModuleId); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var modules []ds.Module
	modules, err = app.repo.GetFlight(mission.UUID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.AllModulesResponse{Modules: modules})
}

func (app *Application) GetAllMissions(c *gin.Context) {
	var request schemes.GetAllMissionsRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	missions, err := app.repo.GetAllMissions(request.DateApproveStart, request.DateApproveEnd, request.Status)
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

func (app *Application) GetMission(c *gin.Context) {
	var request schemes.MissionRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	mission, err := app.repo.GetMissionById(request.MissionId, app.getCustomer())
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
	mission, err := app.repo.GetMissionById(request.URI.MissionId, app.getCustomer())
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

func (app *Application) DeleteMission(c *gin.Context) {
	var request schemes.MissionRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	mission, err := app.repo.GetMissionById(request.MissionId, app.getCustomer())
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

func (app *Application) DeleteFromMission(c *gin.Context) {
	var request schemes.DeleteFromMissionRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	mission, err := app.repo.GetMissionById(request.MissionId, app.getCustomer())
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

func (app *Application) UserConfirm(c *gin.Context) {
	var request schemes.UserConfirmRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	mission, err := app.repo.GetMissionById(request.MissionId, app.getCustomer())
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
	mission.Status = ds.FORMED
	now := time.Now()
	mission.DateApprove = &now

	if err := app.repo.SaveMission(mission); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

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

	if request.Status != ds.COMPELTED && request.Status != ds.REJECTED {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("status %s not allowed", request.Status))
		return
	}

	mission, err := app.repo.GetMissionById(request.URI.MissionId, app.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if mission == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("миссия не найдена"))
		return
	}
	if mission.Status != ds.FORMED {
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("нельзя изменить статус с \"%s\" на \"%s\"", mission.Status, request.Status))
		return
	}
	mission.Status = request.Status
	mission.ModeratorId = app.getModerator()
	if request.Status == ds.COMPELTED {
		now := time.Now()
		mission.DateEnd = &now
	}

	if err := app.repo.SaveMission(mission); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}
