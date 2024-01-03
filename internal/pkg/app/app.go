package app

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"lab1/internal/app/config"
	"lab1/internal/app/dsn"
	"lab1/internal/app/redis"
	"lab1/internal/app/repository"
	"lab1/internal/app/role"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Application struct {
	repo        *repository.Repository
	minioClient *minio.Client
	config      *config.Config
	redisClient *redis.Client
}

func (app *Application) Run() {
	log.Println("Server start up")

	r := gin.Default()

	r.Use(ErrorHandler())

	api := r.Group("/api")
	{

		modules := api.Group("/modules")
		{
			modules.GET("", app.WithAuthCheck(role.NotAuthorized, role.Customer, role.Moderator), app.GetAllModules)        // Список с поиском
			modules.GET("/:module_id", app.WithAuthCheck(role.NotAuthorized, role.Customer, role.Moderator), app.GetModule) // Одна услуга
			modules.DELETE("/:module_id", app.WithAuthCheck(role.Moderator), app.DeleteModule)                              // Удаление
			modules.PUT("/:module_id", app.WithAuthCheck(role.Moderator), app.ChangeModule)                                 // Изменение
			modules.POST("", app.WithAuthCheck(role.Moderator), app.AddModule)                                              // Добавление
			modules.POST("/:module_id/add_to_mission", app.WithAuthCheck(role.Customer, role.Moderator), app.AddToMission)  // Добавление в заявку
		}

		missions := api.Group("/missions")
		{
			missions.GET("", app.WithAuthCheck(role.Customer, role.Moderator), app.GetAllMissions)                                // Список (отфильтровать по дате формирования и статусу)
			missions.GET("/:mission_id", app.WithAuthCheck(role.Customer, role.Moderator), app.GetMission)                        // Одна заявка
			missions.PUT("", app.WithAuthCheck(role.Customer, role.Moderator), app.UpdateMission)                                 // Изменение (добавление)
			missions.DELETE("", app.WithAuthCheck(role.Moderator), app.DeleteMission)                                             //Удаление
			missions.DELETE("/delete_module/:module_id", app.WithAuthCheck(role.Customer, role.Moderator), app.DeleteFromMission) // Изменеие (удаление услуг)
			missions.PUT("/user_confirm", app.WithAuthCheck(role.Customer, role.Moderator), app.UserConfirm)                      // Сформировать создателем
			missions.PUT("/:mission_id/moderator_confirm", app.WithAuthCheck(role.Moderator), app.ModeratorConfirm)               // Завершить отклонить модератором
			missions.PUT("/:mission_id/funding", app.Funding)
		}

		user := api.Group("/user")
		{
			user.POST("/sign_up", app.Register)
			user.POST("/login", app.Login)
			user.POST("/logout", app.Logout)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Println("Host and Port: ", app.config.ServiceHost, app.config.ServicePort)
	r.Run(fmt.Sprintf("localhost:80")) //......................................changeable

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

	app.minioClient, err = minio.New(app.config.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4("", "", ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	app.redisClient, err = redis.New(app.config.Redis)
	if err != nil {
		return nil, err
	}

	return &app, nil
}
