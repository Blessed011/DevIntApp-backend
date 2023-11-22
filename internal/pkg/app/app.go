package app

import (
	"lab1/internal/app/config"
	"lab1/internal/app/dsn"
	"lab1/internal/app/repository"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"log"
	"time"
)

type Application struct {
	repo        *repository.Repository
	minioClient *minio.Client
	config      *config.Config
}

func (app *Application) Run() {

	r := gin.Default()
	//r.LoadHTMLGlob("../../templates/*")

	r.Use(ErrorHandler())

	// Услуги - модули
	r.GET("/modules", app.GetAllModules)                           // Список с поиском
	r.GET("/modules/:module_id", app.GetModule)                    // Одна услуга
	r.DELETE("/modules/:module_id", app.DeleteModule)              // Удаление
	r.PUT("/modules/:module_id", app.ChangeModule)                 // Изменение
	r.POST("/modules", app.AddModule)                              // Добавление
	r.POST("/modules/:module_id/add_to_mission", app.AddToMission) // Добавление в заявку

	// Заявки - миссии
	r.GET("/missions", app.GetAllMissions)                                            // Список (отфильтровать по дате формирования и статусу)
	r.GET("/missions/:mission_id", app.GetMission)                                    // Одна заявка
	r.PUT("/missions/:mission_id/update", app.UpdateMission)                          // Изменение (добавление)
	r.DELETE("/missions/:mission_id", app.DeleteMission)                              //Удаление
	r.DELETE("/missions/:mission_id/delete_module/:module_id", app.DeleteFromMission) // Изменение (удаление услуг)
	r.PUT("/missions/:mission_id/user_confirm", app.UserConfirm)                      // Сформировать создателем
	r.PUT("missions/:mission_id/moderator_confirm", app.ModeratorConfirm)             // Завершить отклонить модератором

	r.Static("/image", "../../resources/images")
	r.Static("/css", "../../static/css")
	r.Run("localhost:8081")
	log.Println("Server down")
}

func New() (*Application, error) {
	var err error
	loc, _ := time.LoadLocation("UTC")
	time.Local = loc
	app := Application{}
	app.config, err = config.NewConfig()
	if err != nil {
		return nil, err
	}

	app.repo, err = repository.New(dsn.FromEnv())
	if err != nil {
		return nil, err
	}

	app.minioClient, err = minio.New(app.config.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4("", "", ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return &app, nil
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			log.Println(err.Err)
		}
		lastError := c.Errors.Last()
		if lastError != nil {
			switch c.Writer.Status() {
			case http.StatusBadRequest:
				c.JSON(-1, gin.H{"error": "wrong request"})
			case http.StatusNotFound:
				c.JSON(-1, gin.H{"error": lastError.Error()})
			case http.StatusMethodNotAllowed:
				c.JSON(-1, gin.H{"error": lastError.Error()})
			default:
				c.Status(-1)
			}
		}
	}
}
